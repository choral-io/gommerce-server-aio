package static

import (
	"embed"
	"io/fs"
)

//go:embed *
//go:embed assets
var efs embed.FS

func FS() fs.FS {
	return efs
}
