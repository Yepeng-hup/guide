package main

import (
	"guide/route"
	"log"
	"os"
)


func main() {
	r := route.InitRoute()
	host := os.Getenv("GUIDE_HOST")
	port := os.Getenv("GUIDE_PORT")
	log.Println("start ok ---> Listening and serving HTTP on "+host+":"+port)
	if err := r.Run(host+":"+port); err != nil {
		log.Println("error start fail", err)
	}
}



