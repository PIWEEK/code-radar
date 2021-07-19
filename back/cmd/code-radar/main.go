package main

import (
	"log"
	"os"
	"strings"

	"net/http"
	"encoding/json"
	"path/filepath"

	"github.com/gorilla/mux"
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
	Name string `json:name`
	Extension string `json:extension`
	Path []string `json:path`
}

type Response struct {
	Name string `json:name`
	Url string `json:url`
	Files []ResponseFile `json:files`
}

type FileDB struct {
	Id string
	Path string
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
	log.Println("Processing commit:", commit.Hash)

	files, err := commit.Files()

	t := db.Txn(true)

	err = files.ForEach(func(file *object.File) error {
		raw, err := t.First("files", "path", file.Name)

		if err == nil && raw == nil {
			filedb := &FileDB{
				Id: uuid.NewString(),
				Path: file.Name,
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

	err = commits.ForEach(ProcessCommit)
	CheckError(err)
}

func ListFiles() []string {
	t := db.Txn(false)
	it, err := t.Get("files", "path")
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
	
	for _, file := range ListFiles() {
		responseFiles = append(responseFiles, ResponseFile {
			Name: filepath.Base(file),
			Extension: filepath.Ext(file),
			Path: strings.Split(filepath.Dir(file), "/"),
		})
	}

	result := &Response {
		Url: info.Url,
		Files: responseFiles,
	}
	
	encoder := json.NewEncoder(w)
	encoder.Encode(result)
}

func Serve() error {
	router := mux.NewRouter().StrictSlash(false)

	router.HandleFunc("/files", RetrieveFiles)

	log.Println("Listening http://localhost:8000/files")
	return http.ListenAndServe(":8000", router)
}

func main() {
	if (len (os.Args) != 2) {
		log.Fatal("coderadar [url]")
	}
	
	url := os.Args[1]
	log.Println("> Repository:", url)

	info = &RepoInfo {
		Url: url,
	}

	InitializeDB()
	ProcessRepository(url)

	log.Fatal(Serve())
}
