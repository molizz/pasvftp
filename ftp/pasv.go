package ftp

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/molizz/bitproxy/utils"
)

const (
	DialOriginPasvTimeout       = 10 * time.Second
	AcceptClinetPasvConnTimeout = 10 * time.Second
)

type PasvServer struct {
	originPasvAddr string
	originPasvHost string
	originPasvPort string
	originPasvConn net.Conn

	proxyPasvHost string
	proxyPasvPort string

	CopyFunc func(client io.ReadWriter, origin io.ReadWriter) (int64, error)
}

func NewPasvServer(originPasvAddr string) *PasvServer {
	originPasvHost, originPasvPort, err := net.SplitHostPort(originPasvAddr)
	if err != nil {
		panic(err)
	}

	ip, _ := utils.PublicIp()
	p := &PasvServer{
		originPasvHost: originPasvHost,
		originPasvPort: originPasvPort,
		originPasvAddr: originPasvAddr,
		proxyPasvHost:  ip,
		proxyPasvPort:  originPasvPort,
	}
	return p
}

func (p *PasvServer) validAndPrepare() (err error) {
	if p.CopyFunc == nil {
		panic("copy func is nil")
	}
	return nil
}

// 连接origin服务器端口
// 发送本地端口信息到客户端
// 监听本地端口
//
func (p *PasvServer) Work() (err error) {
	if err = p.validAndPrepare(); err != nil {
		return err
	}

	clientPasvConn, err := p.acceptClientPasvConn()
	if err != nil {
		return err
	}

	defer func() {
		if clientPasvConn != nil {
			_ = clientPasvConn.Close()
		}
		if p.originPasvConn != nil {
			_ = p.originPasvConn.Close()
		}
	}()

	err = p.dialOriginPasv()
	if err != nil {
		return err
	}

	_, _ = p.CopyFunc(clientPasvConn, p.originPasvConn)
	_ = clientPasvConn.Close()
	_ = p.originPasvConn.Close()
	fmt.Println("Pasv模式完成")
	return nil
}

func (p *PasvServer) acceptClientPasvConn() (net.Conn, error) {
	addr := net.JoinHostPort("", p.proxyPasvPort)
	var err error
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	defer l.Close()

	t := time.AfterFunc(AcceptClinetPasvConnTimeout, func() {
		if p.originPasvConn != nil {
			_ = p.originPasvConn.Close()
		}
		_ = l.Close()
	})
	defer t.Stop()

	conn, err := l.Accept()
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (p *PasvServer) dialOriginPasv() error {
	var err error
	p.originPasvConn, err = net.DialTimeout("tcp", p.originPasvAddr, DialOriginPasvTimeout)
	if err != nil {
		return err
	}
	return nil
}

func (p *PasvServer) proxyPasvAddr() string {
	return net.JoinHostPort(p.proxyPasvHost, p.proxyPasvPort)
}

func (p *PasvServer) buildPasvMsg() []byte {
	msg := []byte(fmt.Sprintf("227 Entering Passive Mode (%s).\r\n", parseAddressToPasvMode(p.proxyPasvAddr())))
	return msg
}
