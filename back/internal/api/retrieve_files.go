package api

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"time"

	"github.com/piweek/code-radar/internal/db"
	"github.com/piweek/code-radar/internal/global"
)

type RetrieveFilesResponse struct {
	Name string `json:"name"`
	Url string `json:"url"`
	FirstCommit string `json:"firstCommit"`
	LastCommit string `json:"lastCommit"`
	Users []string `json:"users"`
	Files []RetrieveFilesResponseFile `json:"files"`
}

type RetrieveFilesResponseFile struct {
	Name string `json:"name"`
	Directory string `json:"directory,omitempty"`
	Extension string `json:"extension,omitempty"`
	Lines int `json:"lines"`
	Rating float32 `json:"rating"`
	IsDirectory bool `json:"isDirectory"`
	History []RetrieveFilesResponseHistoryEntry `json:"history"`
	Owners map[string]interface{} `json:"owners"`
}

type RetrieveFilesResponseHistoryEntry struct {
	User string `json:"user"`
	Added int `json:"added"`
	Deleted int `json:"deleted"`
	Date string `json:"date"`
}

func RetrieveFiles(w http.ResponseWriter, r *http.Request) {
	var err error
	var responseFiles []RetrieveFilesResponseFile

	files, err := db.All()

	if err != nil {
		panic(err)
	}

	for _, file := range files {
		isDir := file.Dir
		directory := file.Parent

		var extension string
		if !isDir {
			extension = filepath.Ext(file.Name)
		}

		var history []RetrieveFilesResponseHistoryEntry

		for _, h := range file.History {
			history = append(history, RetrieveFilesResponseHistoryEntry {
				User: h.User,
				Added: h.Added,
				Deleted: h.Deleted,
				Date: h.Date.Format(time.RFC3339),
			})
		}

		owners := make(map[string]interface{})
		for k, v := range file.Owners {
			owners[k] = v
		}

		responseFiles = append(responseFiles, RetrieveFilesResponseFile {
			Name: file.Name,
			Directory: directory,
			Extension: extension,
			Lines: file.Lines,
			Rating: file.Rating,
			IsDirectory: isDir,
			History: history,
			Owners: owners,
		})
	}

	result := &RetrieveFilesResponse {
		Name: global.GetInfo().Name,
		Url: global.GetInfo().Url,
		FirstCommit: global.GetInfo().FirstCommit.Format(time.RFC3339),
		LastCommit: global.GetInfo().LastCommit.Format(time.RFC3339),
		Users: global.GetInfo().Users,
		Files: responseFiles,
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(result)
}

