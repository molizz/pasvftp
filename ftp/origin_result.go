package ftp

import (
	"fmt"
	"strconv"
	"strings"
)

type OriginResult struct {
	raw  []byte
	Code int
	Msg  string

	IsMul bool
}

func parseOriginResult(raw []byte) (*OriginResult, error) {
	// 222- 第三个字符为 - 的, 表示返回的内容是多行的
	if len(raw) >= 4 && raw[3] == '-' {
		n, err := strconv.Atoi(string(raw[:3]))
		if err != nil {
			panic(err)
		}
		return BuildIsMulOriginResult(n), nil
	}

	var builder strings.Builder
	builder.Write(raw)

	rawArray := strings.SplitN(builder.String(), " ", 2)
	if len(rawArray) < 2 {
		return nil, fmt.Errorf("origin result was error: %s", builder.String())
	}

	code, err := strconv.Atoi(rawArray[0])
	if err != nil {
		return nil, err
	}

	result := &OriginResult{
		raw:  raw,
		Code: code,
		Msg:  rawArray[1],
	}
	return result, nil
}

func BuildIsMulOriginResult(code int) *OriginResult {
	return &OriginResult{IsMul: true, Code: code}
}
