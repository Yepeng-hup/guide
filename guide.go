package main

import (
	"guide/global"
	"guide/route"
	"log"
)


func main() {
	r := route.InitRoute()
	log.Println("INFO: Listening and serving HTTP on "+global.Host+":"+global.Port)
	if err := r.Run(global.Host+":"+global.Port); err != nil {
		log.Println("ERROR: error start fail", err)
	}
}



