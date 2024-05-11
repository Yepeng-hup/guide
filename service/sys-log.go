package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"guide/core"
	"log"
	"net/http"
)

type (
	Data struct {
		LogType string `json:"type"`
		Log     string `json:"log"`
	}
	LogType struct {
		Type string `json:"logType"`
	}
	//DelLogNum struct {
	//	LogNum int `json:"logNum"`
	//}
)

func InsertLogToDb(c *gin.Context) {
	var data Data
	if err := c.BindJSON(&data); err != nil {
		log.Println(err.Error())
		return
	}

	if err := core.InsertActErrorLog(data.LogType, data.Log); err != nil {
		log.Println(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
	})
	return
}

func ShowLog(c *gin.Context) {
	var logType LogType
	if err := c.BindJSON(&logType); err != nil {
		log.Println(err.Error())
		return
	}
	rel, err := core.SelectLog(fmt.Sprintf("SELECT id,newLogDate,types,logtext FROM error_log WHERE types = \"%s\"", logType.Type))
	if err != nil {
		log.Println(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":   http.StatusOK,
		"logRel": rel,
	})
}

func DeleteLimitLog(c *gin.Context) {
	//var delLogNum DelLogNum
	//if err := c.BindJSON(&delLogNum); err != nil {
	//	log.Println(err.Error())
	//	return
	//}
	if err := core.DeleteErrLog(); err != nil {
		log.Println(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
	})
}
