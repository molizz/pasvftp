package ftp

import (
	"fmt"
	"io"
	"net"
	"time"
)

const (
	CodePasv = 227
	CodeQuit = 221
)

const (
	CommandSTOR = "STOR" // 从本地上传
	CommandNLST = "NLST" // 从远程下载
	CommandRETR = "RETR" // 从远程下载(文件)
	CommandLIST = "LIST" // 从远处下载
)

const (
	PasvStateUnknow = iota
	PasvStateUpload
	PasvStateDownload
)

type Proxy struct {
	client net.Conn
	origin net.Conn

	pasvState int
}

func NewProxy(client, origin net.Conn) *Proxy {
	p := new(Proxy)
	p.client = client
	p.origin = origin
	return p
}

func (p *Proxy) Work() error {
	NewClientReader(p.client).Reading(func(cmd *ClientCommand) {
		p.commandStateFromClient(cmd)
		err := p.sendMsgToOrigin(cmd.raw)
		if err != nil {
			panic(fmt.Errorf("send msg to origin was error: ", err))
		}
	})

	NewOriginReader(p.origin).Reading(func(result *OriginResult) {
		p.buildOriginResult(result)
		err := p.sendMsgToClient(result.raw)
		if err != nil {
			panic(err)
		}
	})

	return nil
}

func (p *Proxy) commandStateFromClient(cmd *ClientCommand) {
	switch cmd.Command {
	case CommandSTOR:
		p.pasvState = PasvStateUpload
	case CommandLIST, CommandNLST, CommandRETR:
		p.pasvState = PasvStateDownload
	default:
		p.pasvState = PasvStateUnknow
	}
}

func (p *Proxy) buildOriginResult(result *OriginResult) {
	switch result.Code {
	case CodePasv:
		pasvAddress := parsePasvOriginAddress(result.Msg)
		if pasvAddress != "" {
			pasvServer := NewPasvServer(pasvAddress)
			pasvServer.CopyFunc = func(client io.ReadWriter, origin io.ReadWriter) (int64, error) {
				switch p.pasvState {
				case PasvStateUpload:
					return io.Copy(origin, client)
				case PasvStateDownload:
					return io.Copy(client, origin)
				}
				return 0, nil
			}
			result.raw = pasvServer.buildPasvMsg()
			go func() {
				err := pasvServer.Work()
				if err != nil {
					fmt.Println("pasv server work was err: ", err)
				}
			}()
			time.Sleep(300 * time.Millisecond)
		}
	case CodeQuit:
		_ = p.origin.Close()
	}
}

func (p *Proxy) sendMsgToClient(raw []byte) error {
	_, err := p.client.Write(raw)
	return err
}
func (p *Proxy) sendMsgToOrigin(raw []byte) error {
	_, err := p.origin.Write(raw)
	return err
}
