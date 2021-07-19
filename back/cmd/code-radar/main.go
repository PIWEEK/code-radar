package main

import (
	"log"
	"os"

	"github.com/piweek/code-radar/internal/api"
	"github.com/piweek/code-radar/internal/db"
	"github.com/piweek/code-radar/internal/global"
	"github.com/piweek/code-radar/internal/parser"
)

func main() {
	var err error

	if (len (os.Args) != 3) {
		log.Fatal("coderadar [name] [url]")
	}

	name := os.Args[1]
	url  := os.Args[2]

	log.Println("Repository:", url)
	global.InitInfo(name, url)

	err = db.InitializeDB()

	if err != nil {
		panic(err)
	}
	
	err = parser.ProcessRepository(url)

	if err != nil {
		panic(err)
	}

	log.Fatal(api.Start())
}
