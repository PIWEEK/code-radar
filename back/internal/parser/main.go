package parser

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/piweek/code-radar/internal/db"
)

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func ProcessCommit(commit *object.Commit) error {
	files, err := commit.Files()

	err = files.ForEach(func(file *object.File) error {
		path := file.Name

		err = db.InsertFile(path, false)

		if (err != nil) {
			return err
		}

		dirs := strings.Split(filepath.Dir(path), "/")
		var currentDir string

		for _, dir := range dirs {
			if currentDir == "" {
				currentDir = dir
			} else {
				currentDir = currentDir + "/" + dir
			}

			if currentDir != "." {
				err = db.InsertFile(currentDir, true)
			}
		}
		
		return err
	})

	return err
}

func ProcessRepository(url string) error {
	repo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions {
		URL: url,
			Progress: os.Stdout,
		})

	if err != nil {
		return err
	}

	ref, err := repo.Head()

	if err != nil {
		return err
	}

	options := git.LogOptions {
		From: ref.Hash(),
	}
	
	commits, err := repo.Log(&options)

	if err != nil {
		return err
	}

	log.Println("[START] Processing repository")

	db.StartTransaction()
	err = commits.ForEach(func (commit *object.Commit) error {
		return ProcessCommit(commit)
	})

	if err == nil {
		db.Commit()
	} else {
		db.Rollback()
	}

	if err != nil {
		return err
	}

	log.Println("[END] Processing repository")

	return nil
}
