package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"mime"
	"os"

	"github.com/bibrikthesorcerer/lantern/internal/config"
	"github.com/bibrikthesorcerer/lantern/internal/web"
	clog "github.com/charmbracelet/log"
)

//go:embed static
var staticFiles embed.FS

var musicDir string
var port int
var showConfig bool

func initParseCLIFlags() {
	flag.StringVar(&musicDir, "dir", "", "Path to music directory")
	flag.IntVar(&port, "port", 0, "Port on which server runs")
	flag.BoolVar(&showConfig, "show-Config", false, "Print current Config and exit")
	flag.Parse()
	mime.AddExtensionType(".css", "text/css; charset=utf-8")
	mime.AddExtensionType(".js", "application/javascript; charset=utf-8")
}

func flagsOverrideConfig(c *config.Config) {
	if musicDir != "" {
		c.Dir = musicDir
	}

	if port != 0 {
		c.Port = port
	}
}
func getFS() fs.FS {
	if os.Getenv("DEV") == "1" {
		clog.Info("DEV==1, using disk FS")
		return os.DirFS("./static")
	}

	staticFS, _ := fs.Sub(staticFiles, "static")
	return staticFS
}

func main() {
	initParseCLIFlags()

	conf, err := config.EnsureConfig()
	if err != nil {
		clog.Fatalf("can't ensure Config: %v", err)
	}

	if showConfig {
		data, _ := json.MarshalIndent(conf, "", " ")
		fmt.Println(string(data))
		os.Exit(0)
	}

	// override Config if needed
	flagsOverrideConfig(conf)

	// http setup
	s, err := web.NewServer(conf, getFS())
	if err != nil {
		clog.Fatalf("NewServer setup fail: %s", err)
	}

	needsImport := s.Repo.NeedsImport()
	if needsImport {
		clog.Info("library not present, starting full scan")
		if err := s.Repo.ImportLibrary(s.Conf.Dir); err != nil {
			clog.Fatalf("library import failed: %v", err)
		}
		clog.Info("full scan completed")
	} else {
		go func() {
			clog.Info("starting library sync")

			if err := s.Repo.Sync(s.Conf.Dir); err != nil {
				clog.Errorf("library sync failed: %v", err)
				return
			}

			clog.Info("library sync completed")
		}()
	}

	web.PrintAddrQr(*s.Conf)

	clog.Fatal(s.RunServer()) // TODO: graceful shutdown
}
