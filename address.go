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
	}
	fmt.Println(qr.ToSmallString(false))
	fmt.Printf("\nOr type: %s\n\n", url)
}

func getURL() string {
	conn, err := net.Dial("udp", "8.8.8.8:80") // udp doesn't need handshake
	if err != nil {
		clog.Errorf("couldn't establish connection to resolve local ip: %s", err)
		return "http://localhost:8080"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr).IP.String()

	return fmt.Sprintf("http://%s:8080", localAddr) //TODO: use port from server config
}
