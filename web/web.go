package web

import (
	"embed"
)

//go:embed all:build
var BuildFS embed.FS

//go:embed all:headshot
var HeadshotFS embed.FS
