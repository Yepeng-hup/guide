package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"guide/core"
	"net/http"
	"strings"
)

func SvcCfg(c *gin.Context) {
	s := SvcFrom{
		SvcName:  c.PostForm("svcName"),
		SvcCmd:   c.PostForm("svcCmd"),
		SvcNotes: c.PostForm("svcNotes"),
	}

	if s.SvcName == "" || s.SvcCmd == "" || s.SvcNotes == "" {
		mlog.Warn("add element is nil.")
		c.Redirect(http.StatusFound, "/svc/index")
		return
	}

	err := core.InsertActSTools(s.SvcName, s.SvcCmd, s.SvcNotes)
	if err != nil {
		mlog.Error(err.Error())
		c.Redirect(http.StatusFound, "/svc/index")
		return
	}
	mlog.Info(fmt.Sprintf("service-tools use success, use name -> [%s].", s.SvcName))
	c.Redirect(http.StatusFound, "/svc/index")
}

func ShowSvcCfg(c *gin.Context) {
	svcList, err := core.SelectActSTools("select * from service_tools")
	if err != nil {
		mlog.Error(err.Error())
		return
	}

	c.HTML(http.StatusOK, "protoolscat.tmpl", gin.H{
		"svcList": svcList,
	})
}

func DeleteSvc(c *gin.Context) {
	svcList := strings.Fields(c.PostForm("svc"))
	err := core.DeleteActSTools(svcList[0])
	if err != nil {
		mlog.Error(err.Error())
		return
	}
	c.Redirect(http.StatusFound, "/svc/list")
}
