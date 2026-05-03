package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"mime"
	"os"

	clog "github.com/charmbracelet/log"
)

var musicDir string
var port int
var showConfig bool

func initParseCLIFlags() {
	flag.StringVar(&musicDir, "dir", "", "Path to music directory")
	flag.IntVar(&port, "port", 0, "Port on which server runs")
	flag.BoolVar(&showConfig, "show-config", false, "Print current config and exit")
	flag.Parse()
	mime.AddExtensionType(".css", "text/css; charset=utf-8")
	mime.AddExtensionType(".js", "application/javascript; charset=utf-8")
}

func flagsOverrideConfig(c *Config) {
	if musicDir != "" {
		c.Dir = musicDir
	}

	if port != 0 {
		c.Port = port
	}
}

func main() {
	initParseCLIFlags()

	conf, err := ensureConfig()
	if err != nil {
		clog.Fatalf("can't ensure config: %v", err)
	}

	if showConfig {
		data, _ := json.MarshalIndent(conf, "", " ")
		fmt.Println(string(data))
		os.Exit(0)
	}

	// override config if needed
	flagsOverrideConfig(conf)

	// http setup
	s, err := NewServer(conf)
	if err != nil {
		clog.Fatalf("NewServer setup fail: %s", err)
	}

	printAddrQr(getURL(*s.conf))

	clog.Fatal(s.RunServer())
}
