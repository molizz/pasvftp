package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
)

func JoinHostPort(host string, port uint) string {
	return net.JoinHostPort(host, fmt.Sprintf("%d", port))
}

var publicIp string

func init() {
	var err error
	publicIp, err = PublicIp()
	if err != nil {
		panic(err)
	}
}

func PublicIp() (ip string, err error) {
	if os.Getenv("ENV") == "dev" {
		return "127.0.0.1", nil
	}

	if len(publicIp) > 0 {
		return publicIp, nil
	}

	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("resp status is err:" + resp.Status)
	}

	buff, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	IP := net.ParseIP(string(buff))
	if IP == nil {
		return "", errors.New("ip is invalid:" + string(buff))
	}
	ip = IP.String()
	return
}
