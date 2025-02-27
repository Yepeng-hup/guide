package global

import (
	"os"
)

const (
	FilePath = "text/url.txt"
)

var (
	NowKey       = os.Getenv("GUIDE_KEY")
	DirAccessPwd = os.Getenv("GUIDE_DIR_ACCESS_PWD")
)
