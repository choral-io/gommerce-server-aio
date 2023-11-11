package static

import (
	"embed"
	"io/fs"
)

//go:embed favicon.ico
var favicon embed.FS

func GetFavicon() fs.FS {
	return favicon
}
