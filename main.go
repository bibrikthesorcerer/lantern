package main

import (
	"flag"
	"mime"

	clog "github.com/charmbracelet/log"
)

var MusicDir string

func init() {
	flag.StringVar(&MusicDir, "dir", "", "Specify dir for audio")
	flag.Parse()
	mime.AddExtensionType(".css", "text/css; charset=utf-8")
	mime.AddExtensionType(".js", "application/javascript; charset=utf-8")
}

func main() {
	// http setup
	s, err := NewServer(MusicDir)
	if err != nil {
		clog.Fatalf("NewServer: %s", err)
	}
	SetUpMapping(s)

	printAddrQr(getURL())

	clog.Fatal(s.ListenAndServe())
}
