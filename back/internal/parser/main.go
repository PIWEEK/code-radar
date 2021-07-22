package parser

import (
	"log"
	"fmt"
	"os"
	"context"
	"math"
	"sort"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/go-git/go-git/v5/utils/merkletrie"

	"github.com/piweek/code-radar/internal/db"
	"github.com/piweek/code-radar/internal/global"
)

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

func CalculateOwners(file *db.FileDB) map[string]float32 {
	var totalChanges float32 = 0
	owners := make(map[string]float32)

	for _, h := range file.History {
		v := float32(h.Added + h.Deleted)
		totalChanges += v
		owners[h.User] += v
	}

	if totalChanges != 0 {
		for k, v := range owners {
			owners[k] = v / totalChanges
		}
	}

	return owners
}

func ProcessFiles() error {
	files, _ := db.ListAllFiles()

	min := math.MaxInt32
	max := math.MinInt32

	for _, file := range files {
		var chgs = len(file.History)

		if chgs < min {
			min = chgs
		}

		if chgs > max {
			max = chgs
		}
	}

	for _, file := range files {
		var v = len(file.History)
		rating := float32(v - min) / float32(max - min)

		newFile := file.Copy()
		newFile.Rating = rating
		newFile.Owners = CalculateOwners(file)

		err := db.SaveFile(newFile)

		if err != nil {
			return err
		}
	}

	return nil
}

func ProcessFolders() error {
	dirs, _ := db.ListAllDirectories()

	// We sort by length of the path so we know every children
	// has been processed before the current
	sort.Slice(dirs, func(i, j int) bool {
    return len(strings.Split(dirs[i].Path, "/")) >
			len(strings.Split(dirs[j].Path, "/"))
	})

	for _, dir := range dirs {
		files, _ := db.ListFolderFiles(dir.Path)

		var rating float32 = 0.0

		// Keep the max rating of the children
		for _, file := range files {
			if file.Rating > rating {
				rating = file.Rating
			}
		}

		newDir := dir.Copy()
		newDir.Rating = rating
		newDir.Owners = CalculateOwners(dir)

		err := db.SaveFile(newDir)

		if err != nil {
			return err
		}
	}

	return nil
}

func ProcessRepository(repo *git.Repository) error {
	var err error

	remote, _ := repo.Remote("origin")
	remoteUrl := remote.Config().URLs[0]

	parts := strings.Split(remoteUrl, "/")
	remoteName := parts[len(parts) - 1]
	remoteName = strings.ReplaceAll(remoteName, ".git", "")

	global.InitInfo(remoteName, remoteUrl)

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
		global.AddUser(c.Author.Email)
		commits = append(commits, c)
		return nil
	})

	first := commits[len(commits) - 1].Author.When
	last := commits[0].Author.When

	global.SetCommitRange(first, last)

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

	err = ProcessFiles()

	if err == nil {
		err = ProcessFolders()
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
