package main

import (
	"fmt"
	"net"

	clog "github.com/charmbracelet/log"
	qrcode "github.com/skip2/go-qrcode"
)

func printAddrQr(url string) {
	qr, err := qrcode.New(url, qrcode.Medium)
	if err != nil {
		clog.Errorf("couldn't generate address QR: %s", err)
		return
	}
	fmt.Println("\nFollow QR:\n\n" + qr.ToSmallString(false))
	fmt.Printf("\nOr type: %s\n\n", url)
}

func getURL(c Config) string {
	conn, err := net.Dial("udp", "8.8.8.8:80") // udp doesn't need handshake
	if err != nil {
		clog.Warnf("couldn't establish connection to resolve local ip: %s", err)
		return fmt.Sprintf("http://localhost:%d", c.Port)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr).IP.String()

	return fmt.Sprintf("http://%s:%d", localAddr, c.Port) //TODO: use port from server config
}
