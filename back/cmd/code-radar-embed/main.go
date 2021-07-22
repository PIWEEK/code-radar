package main

import (
	"embed"
	"net/http"
	"github.com/piweek/code-radar/internal/run"
	"github.com/piweek/code-radar/internal/api"
)

//go:embed build favicon.svg global.css index.html
var content embed.FS

func main() {
	api.FileSystem = http.FS(content)
	run.Main()
}
