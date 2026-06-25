package web

import (
	"embed"
	"io/fs"
	"os"
)

//go:embed static
var staticFiles embed.FS

func GetFS() fs.FS {
	if os.Getenv("DEV") == "1" {
		return os.DirFS("internal/web/static")
	}

	fsys, err := fs.Sub(staticFiles, "static")
	if err != nil {
		panic(err)
	}
	return fsys
}
