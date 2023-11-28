package main

import (
	"fmt"
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
	if err != nil {
		log.Fatalln(err,"guide database init fail.")
	}
	fmt.Printf("\033[32m%s\033[0m\n", global.Logo)
	log.Println("INFO: Server version -> 2.5, listening and serving HTTP on "+global.Host+":"+global.Port)
	if err := r.Run(global.Host+":"+global.Port); err != nil {
		log.Println("ERROR: error start fail", err)
	}
}



