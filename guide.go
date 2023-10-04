package main

import (
	"guide/global"
	"guide/route"
	"log"
)


func main() {
	r := route.InitRoute()
	log.Println("start ok ---> Listening and serving HTTP on "+global.Host+":"+global.Port)
	if err := r.Run(global.Host+":"+global.Port); err != nil {
		log.Println("error start fail", err)
	}
}



