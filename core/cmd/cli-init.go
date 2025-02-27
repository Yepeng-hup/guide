package cmd

import (
	"guide/global"
	"log"
)

func CliInit() {
	newKey := global.NowKey
	dirAccessPwd := global.DirAccessPwd

	if newKey == "" {
		log.Fatalln("ERROR: GUIDE_KEY is undefined")
	} else if dirAccessPwd == "" {
		log.Fatalln("ERROR: GUIDE_DIR_ACCESS_PWD is undefined")
	}
}
