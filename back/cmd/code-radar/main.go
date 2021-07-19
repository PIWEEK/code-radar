package main

import (
	"log"
	"os"
	"strings"

	"math/rand"
	"net/http"
	"encoding/json"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/google/uuid"
	"github.com/hashicorp/go-memdb"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
)

var info *RepoInfo
var db *memdb.MemDB

type RepoInfo struct {
	Name string
	Url string
}

type ResponseFile struct {
	Name string `json:"name"`
	Directory string `json:"directory,omitempty"`
	Extension string `json:"extension,omitempty"`
	Lines int `json:"lines"`
	Rating float32 `json:"rating"`
	IsDirectory bool `json:"isDirectory"`
}

type Response struct {
	Name string `json:"name"`
	Url string `json:"url"`
	Files []ResponseFile `json:"files"`
}

type FileDB struct {
	Id string
	Dir bool
	Path string
	Parent string
	Name string
}

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func InitializeDB() {
	var err error
	schema := CreateSchema()
	db, err = memdb.NewMemDB(schema)
	CheckError(err)
}

func CreateSchema() *memdb.DBSchema {
	schema := &memdb.DBSchema {
		Tables: map[string]*memdb.TableSchema {
			"files": &memdb.TableSchema {
				Name: "files",
				Indexes: map[string]*memdb.IndexSchema {
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.UUIDFieldIndex{Field: "Id"},
					},
					"path": &memdb.IndexSchema{
						Name:    "path",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Path"},
					},
					"parent": &memdb.IndexSchema{
						Name:    "parent",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Parent"},
					},
					"name": &memdb.IndexSchema{
						Name:    "name",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Name"},
					},
				},
			},
		},
	}
	return schema
}

func ProcessFile(t *memdb.Txn, path string, dir bool) error {
	raw, err := t.First("files", "path", path)

	parent := filepath.Dir(path)

	if err == nil && raw == nil {
		filedb := &FileDB{
			Id: uuid.NewString(),
			Dir: dir,
			Path: path,
			Parent: parent,
			Name: filepath.Base(path),
		}
		err = t.Insert("files", filedb)
	}

	return err
}

func ProcessCommit(t *memdb.Txn, commit *object.Commit) error {
	files, err := commit.Files()

	err = files.ForEach(func(file *object.File) error {
		path := file.Name

		err = ProcessFile(t, path, false)

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
				err = ProcessFile(t, currentDir, true)
			}
		}
		
		return err
	})

	return err
}

func ProcessRepository(url string){
	repo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions {
		URL: url,
		Progress: os.Stdout,
	})

	CheckError(err)

	ref, err := repo.Head()

	options := git.LogOptions {
		From: ref.Hash(),
	}
	
	commits, err := repo.Log(&options)
	CheckError(err)

	log.Println("[START] Processing repository")
	t := db.Txn(true)

	err = commits.ForEach(func (commit *object.Commit) error {
		return ProcessCommit(t, commit)
	})

	if err == nil {
		t.Commit()
	} else {
		t.Abort()
	}
	
	CheckError(err)
	log.Println("[END] Processing repository")
}

func ListFiles() [](*FileDB) {
	t := db.Txn(false)

	it, err := t.Get("files", "parent")
	CheckError(err)

	var result [](*FileDB)
	for obj := it.Next(); obj != nil; obj = it.Next() {
		p := obj.(*FileDB)
		result = append(result, p)
	}

	return result
}

func RetrieveFiles(w http.ResponseWriter, r *http.Request) {
	var responseFiles []ResponseFile

	for _, file := range ListFiles() {
		isDir := file.Dir

		var directory string

		if file.Parent != "." {
			directory = file.Parent
		}
		
		var extension string
		if !isDir {
			extension = filepath.Ext(file.Name)
		}
		
		responseFiles = append(responseFiles, ResponseFile {
			Name: file.Name,
			Directory: directory,
			Extension: extension,
			Lines: rand.Int() % 10000,
			Rating: rand.Float32(),
			IsDirectory: isDir,
		})
	}

	result := &Response {
		Name: info.Name,
		Url: info.Url,
		Files: responseFiles,
	}
	
	encoder := json.NewEncoder(w)
	encoder.Encode(result)
}

func Serve() error {
	router := mux.NewRouter().StrictSlash(false)

	router.HandleFunc("/files", RetrieveFiles)

	handler := cors.Default().Handler(router)

	log.Println("Listening http://localhost:8000/files")
	return http.ListenAndServe(":8000", handler)
}

func main() {
	if (len (os.Args) != 3) {
		log.Fatal("coderadar [name] [url]")
	}

	name := os.Args[1]
	url := os.Args[2]
	log.Println("> Repository:", url)

	info = &RepoInfo {
		Name: name,
		Url: url,
	}

	InitializeDB()
	ProcessRepository(url)

	log.Fatal(Serve())
}
