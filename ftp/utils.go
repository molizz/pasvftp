package ftp

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
)

func parsePasvOriginAddress(raw string) string {
	r := regexp.MustCompile(`\d+`)
	matches := r.FindAllString(raw, -1)
	matchedCount := len(matches)

	uintBuff := make([]int, matchedCount)
	for i, e := range matches {
		uintBuff[i], _ = strconv.Atoi(e)
	}
	// host
	if len(uintBuff) == 6 {
		ip := net.IPv4(byte(uintBuff[0]), byte(uintBuff[1]), byte(uintBuff[2]), byte(uintBuff[3]))
		port := uint(uintBuff[4]*256 + uintBuff[5])
		return fmt.Sprintf("%s:%d", ip.String(), port)
	}
	return ""
}

func parseAddressToPasvMode(address string) string {
	h, p, err := net.SplitHostPort(address)
	if err != nil {
		panic(err)
	}
	iPort, _ := strconv.Atoi(p)

	h = strings.Replace(h, ".", ",", -1)
	return fmt.Sprintf("%s,%d,%d", h, iPort>>8, iPort&0xFF)
}
