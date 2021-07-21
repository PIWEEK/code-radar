package main

import (
	"log"
	"flag"

	"github.com/go-git/go-git/v5"
	"github.com/piweek/code-radar/internal/api"
	"github.com/piweek/code-radar/internal/db"
	"github.com/piweek/code-radar/internal/parser"
)

func main() {
	var err error

	isLocal := flag.Bool("local", false, "the url refers to a local repository path")
	flag.Parse()

	args := flag.Args()

	if len(args) != 0 && len(args) != 1 {
		log.Fatal("coderadar [--local] [<url|path>]")
	}

	var url string

	if len(args) == 1 {
		url  = args[0]
	} else {
		tmp := true
		isLocal = &tmp
		url = "."
	}

	err = db.InitializeDB()

	if err != nil {
		panic(err)
	}

	var repo *git.Repository

	if *isLocal {
		repo = parser.InitLocalRepository(url)
	} else {
		repo = parser.InitRemoteRepository(url)
	}
	
	err = parser.ProcessRepository(repo)

	if err != nil {
		panic(err)
	}

	log.Fatal(api.Start())
}
