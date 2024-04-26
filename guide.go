package main

import (
	"guide/core"
	"guide/core/cmd"
	"guide/global"
	"guide/route"
	"log"
)

func main() {
	cmd.CliInit()
	r := route.InitRoute()
	err := core.CreateGuideAllTable()
	go core.InitUser()
	if err != nil {
		log.Fatalln(err, "guide database init fail.")
	}
	log.Println("INFO: Server version -> 4.0, listening and serving HTTP on " + global.Host + ":" + global.Port)
	if err := r.Run(global.Host + ":" + global.Port); err != nil {
		log.Println("ERROR: error start fail", err)
	}
}
