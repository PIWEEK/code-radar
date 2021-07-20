package db

type FileDB struct {
	Id string // Unique identifier UUID
	Dir bool
	Path string
	Parent string
	Name string
	Lines int
	Rating float32
}

