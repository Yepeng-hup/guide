package service

import (
	"github.com/gin-gonic/gin"
	"guide/core"
	"log"
	"net/http"
)

func SvcCfg(c *gin.Context)  {
	s := SvcFrom{
		SvcName: c.PostForm("svcName"),
		SvcCmd: c.PostForm("svcCmd"),
		SvcNotes: c.PostForm("svcNotes"),
	}

	if s.SvcName == ""|| s.SvcCmd == ""||s.SvcNotes == ""{
		log.Println("WARN: add element is nil.")
		c.Redirect(http.StatusFound, "/svc/index")
		return
	}

	err := core.InsertActSTools(s.SvcName, s.SvcCmd, s.SvcNotes)
	if err != nil {
		log.Println(err)
		c.Redirect(http.StatusFound, "/svc/index")
		return
	}
	log.Printf("INFO: service-tools use success, use name -> [%s].", s.SvcName)
	c.Redirect(http.StatusFound, "/svc/index")
	return
}


func ShowSvcCfg(c *gin.Context){
	svcList, err := core.SelectActSTools("select * from service_tools")
	if err != nil{
		log.Println(err)
		return
	}

	c.HTML(http.StatusOK, "protoolscat.tmpl", gin.H{
		"svcList": svcList,
	})
}
