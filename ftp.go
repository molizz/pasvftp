package main

import (
	"fmt"
	"github.com/molizz/pasvftp/ftp"
)

func main() {
	fmt.Println("pasvftp by moli")

	p := ftp.NewFtpProxy(2121, "ftp.helloworld.com", 21)
	_ = p.Start()
}
