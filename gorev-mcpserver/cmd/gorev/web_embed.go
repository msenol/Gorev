package main

import (
	"embed"
)

//go:embed web/dist
var WebDistFS embed.FS
