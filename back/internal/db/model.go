package db

import (
	"time"
)

type FileHistory struct {
	User string
	Added int
	Deleted int
	Date time.Time
}

type FileDB struct {
	Id string // Unique identifier UUID
	Dir bool
	Path string
	Parent string
	Name string
	Lines int
	Rating float32
	History []FileHistory
	Owners map[string]float32
}

func (file *FileDB) Copy() *FileDB {
	var newFile FileDB
	newFile = (*file)
	newFile.History = make([]FileHistory, len(newFile.History))

	// Copy history elements
	for i, h := range file.History {
		newFile.History[i] = h
	}

	newFile.Owners = make(map[string]float32)

	for k, v := range file.Owners {
		newFile.Owners[k] = v
	}

	return &newFile
}
