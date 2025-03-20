package main

import (
	"guide/core"
	"guide/core/cmd"
	"guide/route"
	"log"
)

func main() {
	core.ReadJson("conf.d/debug.json")

	cmd.CliInit()

	r := route.InitRoute()

	if err := core.CreateGuideAllTable(); err != nil {
		log.Fatalln("ERROR: guide database init fail.", err.Error())
	}

	go core.InitUser()

	routes := r.Routes()
	core.InitPermissionRoute(routes)

	// if global.Mon == "true" {go service.CpuValueWDB()}
	log.Println("INFO: Server version -> 4.0, listening and serving HTTP on " + core.Cfg.ListenHost + ":" + core.Cfg.ListenPort)
	if err := r.Run(core.Cfg.ListenHost + ":" + core.Cfg.ListenPort); err != nil {
		log.Println("ERROR: error start fail", err)
	}

}
