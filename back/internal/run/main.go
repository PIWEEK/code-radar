package run

import (
	"log"
	"flag"
	"net"

	"github.com/go-git/go-git/v5"
	"github.com/piweek/code-radar/internal/api"
	"github.com/piweek/code-radar/internal/db"
	"github.com/piweek/code-radar/internal/parser"
)

type RunArgs struct {
	Url string
	Port int
	Local bool
}

func GetFreePort() (int, error) {
	listener, err := net.Listen("tcp", ":0")

	if err != nil {
    return -1, err
	}

	port := listener.Addr().(*net.TCPAddr).Port
	return port, nil
}

func ProcessArgs() RunArgs {
	isLocal := flag.Bool("local", false, "the url refers to a local repository path")
	port := flag.Int("port", -1, "listening port")
	
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

	var finalPort int = *port

	if finalPort == -1 {
		finalPort, _ = GetFreePort()
	}

	return RunArgs {
		Url: url,
		Port: finalPort,
		Local: *isLocal,
	}
}

func Main() {
	var err error

	args := ProcessArgs()

	if err = db.InitializeDB(); err != nil {
		panic(err)
	}

	var repo *git.Repository

	if args.Local {
		repo = parser.InitLocalRepository(args.Url)
	} else {
		repo = parser.InitRemoteRepository(args.Url)
	}
	
	if err = parser.ProcessRepository(repo); err != nil {
		panic(err)
	}

	log.Fatal(api.Start(args.Port))
}
