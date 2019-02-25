package ftp

import (
	"bufio"
	"fmt"
	"net"
)

type ClientReadCallbackFunc func(*ClientCommand)

type ClientReader struct {
	client net.Conn
}

func NewClientReader(client net.Conn) *ClientReader {
	return &ClientReader{
		client: client,
	}
}

func (c *ClientReader) Reading(callbackFunc ClientReadCallbackFunc) {
	client := bufio.NewReader(c.client)

	go func() {
		for {
			cmdRaw, err := client.ReadBytes('\n')
			if err != nil {
				fmt.Println(err)
				return
			}

			cmd, err := parseCommand(cmdRaw)
			if err != nil {
				fmt.Println(err)
				return
			}
			callbackFunc(cmd)
		}
	}()
}
