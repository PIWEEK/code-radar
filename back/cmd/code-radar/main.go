package main

import (
	"net/http"
	"github.com/piweek/code-radar/internal/run"
	"github.com/piweek/code-radar/internal/api"
)

func main() {
	api.FileSystem = http.Dir("./static/")
	run.Main()
}
