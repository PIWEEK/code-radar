package global

var info *Info

type Info struct {
	Name string
	Url string
}

func InitInfo(name string, url string) {
	info = &Info {
		Name: name,
		Url: url,
	}
}

func GetInfo() *Info {
	return info
}
