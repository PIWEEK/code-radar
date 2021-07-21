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
	var prevTree *object.Tree

	// If the previous commit is empty we compare with the empty tree
	if (previous != nil) {
		prevTree, _ = previous.Tree()
	} else {
		prevTree = &object.Tree{}
	}

	options := &object.DiffTreeOptions {
		DetectRenames: true,
	}

	commitTree, _ := commit.Tree()

	ctx := context.Background()
	changes, err := object.DiffTreeWithOptions(ctx, prevTree, commitTree, options)

	if err != nil {
		return err
	}

	statsMap := make(map[string]object.FileStat)

	patch, _ := changes.Patch()
	stats := patch.Stats()

	for _, stat := range stats {
		statsMap[stat.Name] = stat
	}

	user := commit.Author.Email
	date := commit.Author.When

	for _, change := range changes {
		action, _ := change.Action()
		name := change.To.Name
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

		  	if (fileStat.Addition != 0 || fileStat.Deletion != 0) {
		  		err = db.UpdateFile(name, false, fileStat.Addition, fileStat.Deletion, user, date)
		  	}

		  case merkletrie.Delete:
		  	name = change.From.Name
		  	err = db.DeleteFile(name, user, date)
		  }

		  if err != nil {
		  	return err
		  }
	}

	return nil
}

func InitLocalRepository(url string) *git.Repository {
	log.Println("Init local repository", url)
	repo, err := git.PlainOpenWithOptions(url, &git.PlainOpenOptions{DetectDotGit: true})

	if err != nil {
		panic(err)
	}

	return repo
}

func InitRemoteRepository(url string) *git.Repository {
	log.Println("Init remote repository", url)
	repo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions {
		URL: url,
		Progress: os.Stdout,
	})

	if err != nil {
		panic(err)
	}

	return repo
}

func ProcessRepository(repo *git.Repository) error {
	var err error

	ref, err := repo.Head()

	if err != nil {
		return err
	}

	commitIt, err := repo.Log(&git.LogOptions{
		From: ref.Hash(),
		Order: git.LogOrderCommitterTime,
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
