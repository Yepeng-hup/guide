package main

import (
	"guide/core"
	"guide/core/cmd"
	"guide/global"
	"guide/route"
	_"guide/service"
	"log"
)

func main() {
	cmd.CliInit()
	r := route.InitRoute()
	if err := core.CreateGuideAllTable(); err != nil {
		log.Fatalln("ERROR: guide database init fail.", err.Error())
	}
	go core.InitUser()
	// if global.Mon == "true" {go service.CpuValueWDB()}
	log.Println("INFO: Server version -> 4.0, listening and serving HTTP on " + global.Host + ":" + global.Port)
	if err := r.Run(global.Host + ":" + global.Port); err != nil {
		log.Println("ERROR: error start fail", err)
	}
}
