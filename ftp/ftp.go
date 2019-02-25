package ftp

import (
	"net"
	"time"

	"github.com/molizz/pasvftp/utils"
)

type FtpProxy struct {
	localPort  uint
	remoteHost string
	remotePort uint

	done bool

	log *utils.Logger
	ln  net.Listener
}

func (this *FtpProxy) dialOrigin() (originConn net.Conn, err error) {
	originConn, err = net.DialTimeout("tcp", utils.JoinHostPort(this.remoteHost, this.remotePort), 10*time.Second)
	if err != nil {
		return nil, err
	}
	return originConn, nil
}

func (this *FtpProxy) handle(client net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			this.log.Printf("ftp proxy was panic: ", err)
		}
		_ = client.Close()
	}()

	originConn, err := this.dialOrigin()
	if err != nil {
		this.log.Info("could't dial to origin server: ", err)
		return
	}
	defer originConn.Close()

	proxy := NewProxy(client, originConn)
	err = proxy.Work()
	if err != nil {
		this.log.Info("proxy work was error: ", err)
		return
	}
}

func (this *FtpProxy) Start() (err error) {
	this.log.Info("Listen port", this.localPort)

	this.ln, err = net.Listen("tcp", utils.JoinHostPort("", this.localPort))
	if err != nil {
		this.log.Info("Can't Listen port: ", this.localPort, " ", err)
		return err
	}
	for !this.done {
		conn, err := this.ln.Accept()
		if err != nil {
			// this.log.Info("Can't Accept: ", this.localPort, " ", err)
			if netErr, ok := err.(net.Error); ok && netErr.Temporary() {
				time.Sleep(100 * time.Millisecond)
				continue
			} else {
				return err
			}
		}
		this.log.Info("Accept ", conn.RemoteAddr().String())
		go this.handle(conn)
	}
	return nil
}

func (this *FtpProxy) Stop() error {
	this.done = true
	if this.ln != nil {
		_ = this.ln.Close()
	}
	return nil
}

func (this *FtpProxy) LocalPort() uint {
	return this.localPort
}

func (this *FtpProxy) Traffic() (uint64, error) {
	return 0, nil
}

func NewFtpProxy(localPort uint, remoteHost string, remotePort uint) *FtpProxy {
	return &FtpProxy{
		localPort:  localPort,
		remoteHost: remoteHost,
		remotePort: remotePort,
		log:        utils.NewLogger("FtpProxy"),
	}
}
