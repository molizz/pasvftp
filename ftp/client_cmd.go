package ftp

import (
	"strings"
)

type ClientCommand struct {
	raw     []byte
	Command string
	Params  string
}

func parseCommand(raw []byte) (*ClientCommand, error) {
	var rawString strings.Builder
	rawString.Write(raw)

	cmdRaw := strings.Split(rawString.String(), " ")
	var p1, p2 string
	if len(cmdRaw) >= 1 {
		p1 = cmdRaw[0]
	}
	if len(cmdRaw) >= 2 {
		p2 = cmdRaw[1]
	}

	result := &ClientCommand{
		raw:     raw,
		Command: p1,
		Params:  p2,
	}
	return result, nil
}
