package parser

import (
	"log"
	"os"
	//	"path/filepath"
	//"strings"
	"context"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/go-git/go-git/v5/utils/merkletrie"
	"github.com/piweek/code-radar/internal/db"
)

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func ProcessCommit(commit *object.Commit, previous *object.Commit) error {
	
	stats, err := commit.Stats()

	if err != nil {
		return err
	}

	// log.Println("Processing", commit.Hash)
	// log.Println("====")
	// log.Println(commit.Message)
	// log.Println(commit.Author.When, commit.Author.Name, commit.Author.Email)
	// log.Println("====")

	statsMap := make(map[string]object.FileStat)

	for _, stat := range stats {
		statsMap[stat.Name] = stat
	}

	if (previous != nil) {
		options := &object.DiffTreeOptions {
			DetectRenames: true,
		}

		commitTree, _ := commit.Tree()
		prevTree, _ := previous.Tree()

		ctx := context.Background()
		changes, err := object.DiffTreeWithOptions(ctx, prevTree, commitTree, options)

		if err != nil {
			return err
		}

		for _, change := range changes {
			name := change.To.Name

			switch acc, _ := change.Action(); acc {
			  case merkletrie.Insert:
				  err = db.UpdateFile(name, false, statsMap[name].Addition, 0, 0 /*rating*/)
			  	
			  case merkletrie.Modify:
			  	nameFrom := change.From.Name

				  if nameFrom != name {
				  	db.MoveFile(nameFrom, name)
				  }
				  
				  err = db.UpdateFile(name, false, statsMap[name].Addition, statsMap[name].Deletion, 0)
				
			  	
			  case merkletrie.Delete:
				  err = db.DeleteFile(name)
			}

			if err != nil {
				return err
			}
		}
	}

	return nil
}

/*
func ProcessCommit_newold(commit *object.Commit) error {
	stats, err := commit.Stats()

	if err != nil {
		return err
	}

	log.Println("====")
	log.Println(commit.Message)
	log.Println(commit.Author.When, commit.Author.Name, commit.Author.Email)
	log.Println("====")
	for _, stat := range stats {
		log.Println(stat.Name, stat.Addition, stat.Deletion)
	}

	return nil
}

func ProcessCommit_old(commit *object.Commit) error {
	files, err := commit.Files()

	err = files.ForEach(func(file *object.File) error {
		var currentDir string

		path := file.Name

		// Insert entries for each of its parents
		dirs := strings.Split(filepath.Dir(path), "/")

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

		err = db.InsertFile(path, false)

		if (err != nil) {
			return err
		}
		
		return err
	})

	return err
}
*/
   
func ProcessRepository(url string) error {
	var repo *git.Repository
	var err error

	if (url == "") {
		log.Println("Using current directory")
		repo, err = git.PlainOpenWithOptions(".", &git.PlainOpenOptions{DetectDotGit: true})
	} else {
		log.Println("Cloning ", url)
		repo, err = git.Clone(memory.NewStorage(), nil, &git.CloneOptions {
			URL: url,
			Progress: os.Stdout,
		})
	}

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
	
	commitIt, err := repo.Log(&options)

	if err != nil {
		return err
	}

	log.Println("[START] Processing repository")

	var commits [](*object.Commit)

	err = commitIt.ForEach(func (c *object.Commit) error {
		commits = append(commits, c)
		return nil
	})

	db.StartTransaction()

	for i := len(commits) - 1 ; i >= 0; i-- {
		log.Printf("Processing %d/%d\n", len(commits) - i, len(commits))
		var prev *object.Commit

		if i < len(commits) - 1 {
			prev = commits[i+1]
		}
		
		err = ProcessCommit(commits[i], prev)

		if err != nil {
			break;
		}
	}
	
	if err == nil {
		db.Commit()
	} else {
		db.Rollback()
		return err
	}

	log.Println("[END] Processing repository")

	return nil
}
