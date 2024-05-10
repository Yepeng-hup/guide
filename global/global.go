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
	WhiteFile = os.Getenv("GUIDE_WHITE_LIST")
	IsStartWhitelist = os.Getenv("GUIDE_START_WHITE_LIST")
	PasswdAdminWhitelist = os.Getenv("GUIDE_PWD_ADMIN_WHITELIST")
	NowKey = os.Getenv("GUIDE_KEY")
	//InterfaceName = os.Getenv("GUIDE_INTERFACE_NAME")
	Mon = os.Getenv("GUIDE_MON")
)


func IpList()[]string{
	whiteIpList := strings.Split(WhiteFile, ",")
	return whiteIpList
}

