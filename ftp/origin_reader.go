package ftp

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
)

type OriginReadCallbackFunc func(*OriginResult)

type OriginReader struct {
	origin net.Conn
}

func NewOriginReader(origin net.Conn) *OriginReader {
	return &OriginReader{
		origin: origin,
	}
}

func (o *OriginReader) Reading(callbackFunc OriginReadCallbackFunc) {
	origin := bufio.NewReader(o.origin)

	for {
		cmdRaw, err := origin.ReadBytes('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		cmd, err := parseOriginResult(cmdRaw)
		if err != nil {
			fmt.Println(err)
			return
		}
		if cmd.IsMul {
			var rawBuff bytes.Buffer
			rawBuff.Write(cmdRaw)
			for {
				line, err := origin.ReadBytes('\n')
				if err != nil {
					break
				}
				rawBuff.Write(line)
				if bytes.HasSuffix(line, []byte("End\r\n")) ||
					bytes.HasSuffix(line, []byte("End of list\r\n")) {
					break
				}
			}
			cmd.raw = rawBuff.Bytes()
		}
		callbackFunc(cmd)
	}
}
