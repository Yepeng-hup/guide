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
	passwdAdminWhitelist := global.PasswdAdminWhitelist
	newKey := global.NowKey
	mon := global.Mon
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
	}else if passwdAdminWhitelist == ""{
		log.Fatalln("ERROR: GUIDE_PWD_ADMIN_WHITELIST is undefined")
	}else if newKey == "" {
		log.Fatalln("ERROR: GUIDE_KEY is undefined")
	}else if mon == "" {
		log.Fatalln("ERROR: GUIDE_MON is undefined")
	}
	return
}