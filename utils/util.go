package utils

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"time"
)

func JoinHostPort(host string, port uint) string {
	return net.JoinHostPort(host, fmt.Sprintf("%d", port))
}

var publicIp string

func PublicIp() (ip string, err error) {
	if os.Getenv("ENV") == "dev" {
		return "127.0.0.1", nil
	}

	if len(publicIp) > 0 {
		return publicIp, nil
	}
	conn, err := net.DialTimeout("tcp", "ns1.dnspod.net:6666", 10*time.Second)
	if err != nil {
		return
	}
	defer func() {
		conn.Close()
	}()
	err = conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		return
	}
	buff, err := ioutil.ReadAll(conn)
	if err != nil {
		return
	}
	ip = string(buff)
	publicIp = ip
	return
}
