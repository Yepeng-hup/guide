package global

import (
	"os"
	"strings"
)

const  (
	FilePath = "text/url.txt"
	Logo = `
							  _ _ _ _ _    _ __ _ _ _ _ _  _ _ _ __
							/  *_*_||*|    |*|| * ||* _\*\ | *_*_*_|
		Welcome to 			| *|  _ |*|    |*|| * ||*|  |*||*|_ __
							| *|_|*||*|_ _ |*|| * ||*|_ |*||*|_ _ _
							\ _*_*_|\*_*_*_*/ |_*_||*_*_*/ |_*_*_*_|
	`
)

var (
	SaveDataDir = os.Getenv("GUIDE_FILEDATA_DIR")
	Host = os.Getenv("GUIDE_HOST")
	Port = os.Getenv("GUIDE_PORT")
	whiteFile = os.Getenv("GUIDE_WHITE_LIST")
	IsStartWhitelist = os.Getenv("GUIDE_START_WHITE_LIST")
	InterfaceName = os.Getenv("GUIDE_INTERFACE_NAME")
)


func IpList()[]string{
	whiteIpList := strings.Split(whiteFile, ",")
	return whiteIpList
}