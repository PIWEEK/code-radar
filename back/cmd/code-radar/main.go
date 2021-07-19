package main

import (
	"log"
	"os"
	// "strings"

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
	Extension string `json:"extension"`
	Lines int `json:lines`
	Rating float32 `json:rating`
	// Path string `json:path`
}

type Response struct {
	Name string `json:"name"`
	Url string `json:"url"`
	Files []ResponseFile `json:"files"`
}

type FileDB struct {
	Id string
	Path string
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
				},
			},
		},
	}
	return schema
}

func ProcessCommit(commit *object.Commit) error {
	files, err := commit.Files()

	t := db.Txn(true)

	err = files.ForEach(func(file *object.File) error {
		path := file.Name
		raw, err := t.First("files", "path", path)

		if err == nil && raw == nil {
			filedb := &FileDB{
				Id: uuid.NewString(),
				Path: path,
			}
			err = t.Insert("files", filedb)
		}

		return err
	})

	if err == nil {
		t.Commit()
	} else {
		t.Abort()
		return err
	}

	return err
}

func ProcessRepository(url string) {
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
	err = commits.ForEach(ProcessCommit)
	CheckError(err)
	log.Println("[END] Processing repository")
}

func ListFiles(path string) []string {
	t := db.Txn(false)
	it, err := t.Get("files", "path_prefix", path)
	CheckError(err)

	var result []string
	for obj := it.Next(); obj != nil; obj = it.Next() {
		p := obj.(*FileDB)
		result = append(result, p.Path)
	}

	return result
}

func RetrieveFiles(w http.ResponseWriter, r *http.Request) {
	var responseFiles []ResponseFile

	path := r.URL.Query().Get("path")
	
	for _, file := range ListFiles(path) {
		// log.Println("??", filepath.Dir(file), path)
		if path == "" && filepath.Dir(file) == "." ||
			 path != "" && filepath.Dir(file) == path {
			responseFiles = append(responseFiles, ResponseFile {
				Name: filepath.Base(file),
				Extension: filepath.Ext(file),
				Lines: rand.Int() % 10000,
				Rating: rand.Float32(),
				//Path: file,
			})
		}
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
