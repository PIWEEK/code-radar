package global

import (
	"time"
)

var info *Info

type Info struct {
	Name string
	Url string
	FirstCommit time.Time 
	LastCommit time.Time
	Users []string
}

func InitInfo(name string, url string) {
	info = &Info {
		Name: name,
		Url: url,
	}
}

func AddUser(user string) {
	for _, u := range info.Users {
		if u == user {
			// User already in. Return
			return
		}
	}
	info.Users = append(info.Users, user)
}

func SetCommitRange(first time.Time, last time.Time) {
	info.FirstCommit = first
	info.LastCommit = last
}

func GetInfo() *Info {
	return info
}
