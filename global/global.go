package global

import "os"

const  (
	FilePath = "text/url.txt"
)

var (
	GinDebug  = os.Getenv("GUIDE_GIN_DEBUG")
	SaveDataDir = os.Getenv("GUIDE_FILEDATA_DIR")
	Host = os.Getenv("GUIDE_HOST")
	Port = os.Getenv("GUIDE_PORT")
)