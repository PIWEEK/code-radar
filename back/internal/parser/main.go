package parser

import (
	"log"
	"fmt"
	"os"
	"context"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/go-git/go-git/v5/utils/merkletrie"
	"github.com/piweek/code-radar/internal/db"
)

func ProcessCommit(commit *object.Commit, previous *object.Commit) error {
	if (previous != nil) {
		options := &object.DiffTreeOptions {
			DetectRenames: true,
		}

		commitTree, _ := commit.Tree()
		prevTree, _ := previous.Tree()

		ctx := context.Background()
		changes, err := object.DiffTreeWithOptions(ctx, prevTree, commitTree, options)

		statsMap := make(map[string]object.FileStat)

		patch, _ := changes.Patch()

		stats := patch.Stats()

		for _, stat := range stats {
			statsMap[stat.Name] = stat
		}

		if err != nil {
			return err
		}

		user := commit.Author.Email
		date := commit.Author.When

		for _, change := range changes {
			name := change.To.Name
			action, _ := change.Action()
			fileStat := statsMap[name]

			switch action {
			  case merkletrie.Insert:
				  err = db.UpdateFile(name, false, fileStat.Addition, 0, user, date)
			  	
			  case merkletrie.Modify:
			  	nameFrom := change.From.Name

				  if nameFrom != name {
				  	err = db.MoveFile(nameFrom, name, user, date)

						if err != nil {
							return err
						}
				  }
				  
				  err = db.UpdateFile(name, false, fileStat.Addition, fileStat.Deletion, user, date)
			  	
			  case merkletrie.Delete:
				  err = db.DeleteFile(name, user, date)
			}

			if err != nil {
				return err
			}
		}
	}

	return nil
}

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

	commitIt, err := repo.Log(&git.LogOptions{
		From: ref.Hash(),
	})

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

	fmt.Printf("Processing %d/%d\r", 0, len(commits))

	for i := len(commits) - 1 ; i >= 0; i-- {
		fmt.Printf("Processing %d/%d\r", len(commits) - i, len(commits))

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
