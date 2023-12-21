package static

import (
	"embed"
	"io/fs"
)

//go:embed *
//go:embed assets
var staticFS embed.FS

func GetStaticFS() fs.FS {
	return staticFS
}
