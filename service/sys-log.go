package service

import (
	"fmt"
	"guide/core"
	"log"
	"net/http"
	"regexp"
	"github.com/gin-gonic/gin"
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

// func InsertLogToDb(c *gin.Context) {
// 	var data Data
// 	if err := c.BindJSON(&data); err != nil {
// 		log.Println(err.Error())
// 		return
// 	}

// 	if err := core.InsertActErrorLog(data.LogType, data.Log); err != nil {
// 		log.Println(err.Error())
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{
// 		"code": http.StatusOK,
// 	})
// }

func screenErrorAndWarn(logContent string) (bool, string) {
	// Define the regular expression pattern to match, ignoring capitalization.
	patternErr := "(?i)error"
	patternWarn := "(?i)warn"
	patternFail := "(?i)fail"

	// Compile regular expressions.
	regexErr := regexp.MustCompile(patternErr)
	if regexErr.MatchString(logContent) {
		return true, "error"
	} else {
		regexWarn := regexp.MustCompile(patternWarn)
		if regexWarn.MatchString(logContent) {
			return true, "warn"
		} else {
			regexFail := regexp.MustCompile(patternFail)
			if regexFail.MatchString(logContent) {
				return true, "other"
			}
		}
	}
	return false, ""
}

func FluentBit(c *gin.Context){
	var jsonData []map[string]interface{}
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, v := range jsonData {
		b, r := screenErrorAndWarn(v["log"].(string))
		if b {
			// 
			if err := core.InsertActErrorLog(r, v["log"].(string)); err != nil {
				log.Println(err.Error())
				return
			}
		}else {
			c.JSON(http.StatusOK, gin.H{
				"info": "not log type",
			})
		}
	}
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
