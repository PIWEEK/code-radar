package api

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"math/rand"
	"fmt"
	"time"

	"github.com/piweek/code-radar/internal/db"
	"github.com/piweek/code-radar/internal/global"
)

type RetrieveFilesResponse struct {
	Name string `json:"name"`
	Url string `json:"url"`
	Files []RetrieveFilesResponseFile `json:"files"`
}

type RetrieveFilesResponseFile struct {
	Name string `json:"name"`
	Directory string `json:"directory,omitempty"`
	Extension string `json:"extension,omitempty"`
	Lines int `json:"lines"`
	Rating float32 `json:"rating"`
	IsDirectory bool `json:"isDirectory"`
	History []string `json:"history"`
}


func RetrieveFiles(w http.ResponseWriter, r *http.Request) {
	var err error
	var responseFiles []RetrieveFilesResponseFile

	files, err := db.ListFiles()

	if err != nil {
		panic(err)
	}

	for _, file := range files {
		isDir := file.Dir

		var directory string

		if file.Parent != "." {
			directory = file.Parent
		}

		var extension string
		if !isDir {
			extension = filepath.Ext(file.Name)
		}

		var history []string

		for _, h := range file.History {
			history = append(history, fmt.Sprintf("[%s]%s: +%d -%d",
				h.Date.Format(time.RFC3339),
				h.User,
				h.Added,
				h.Deleted,
			))
		}

		responseFiles = append(responseFiles, RetrieveFilesResponseFile {
			Name: file.Name,
			Directory: directory,
			Extension: extension,
			Lines: file.Lines,
			Rating: rand.Float32(),
			IsDirectory: isDir,
			History: history,
		})
	}

	result := &RetrieveFilesResponse {
		Name: global.GetInfo().Name,
		Url: global.GetInfo().Url,
		Files: responseFiles,
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(result)
}

