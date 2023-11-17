package cmd

import (
	"guide/global"
	"log"
)

func CliInit (){
	saveDataDir := global.SaveDataDir
	w := global.WhiteFile
	lsStartWhitelist := global.IsStartWhitelist
	host := global.Host
	port := global.Port
	//interfaceName := global.InterfaceName
	if saveDataDir == ""{
		log.Fatalln("ERROR: GUIDE_FILEDATA_DIR is undefined")
	}else if w == ""{
		log.Fatalln("ERROR: GUIDE_WHITE_LIST is undefined")
	}else if host == ""{
		log.Fatalln("ERROR: GUIDE_HOST is undefined")
	}else if port == ""{
		log.Fatalln("ERROR: GUIDE_PORT is undefined")
	}else if lsStartWhitelist == ""{
		log.Fatalln("ERROR: GUIDE_START_WHITE_LIST is undefined")
	}
	return
}