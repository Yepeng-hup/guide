package global

import (
	"os"
	"strings"
)

const  (
	FilePath = "text/url.txt"
)

var (
	SaveDataDir = os.Getenv("GUIDE_FILEDATA_DIR")
	Host = os.Getenv("GUIDE_HOST")
	Port = os.Getenv("GUIDE_PORT")
	whiteFile = os.Getenv("GUIDE_WHITE_LIST")
	IsStartWhitelist = os.Getenv("GUIDE_START_WHITE_LIST")
)


func IpList()[]string{
	whiteIpList := strings.Split(whiteFile, ",")
	return whiteIpList
}