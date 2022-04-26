package assets

import (
	"embed"
	"io/fs"
)

//go:embed *
var FS embed.FS

func GetFile(name string) (fs.File, error) {
	return FS.Open(name)
}

func ReadFile(name string) ([]byte, error) {
	return FS.ReadFile(name)
}
