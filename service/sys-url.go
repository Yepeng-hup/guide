package service

import (
	"bufio"
	"fmt"
	"github.com/gin-gonic/gin"
	"guide/global"
	"log"
	"net/http"
	"os"
	"strings"
	"encoding/json"
)

var (
	lines []string
)

func writeFile(fileName, url, notes string) error {
	file, err := os.OpenFile(global.FilePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	defer file.Close()
	text := fileName + " " + url + " " + notes + "\n"
	if _, err := file.WriteString(text); err != nil {
		return fmt.Errorf(err.Error())
	}
	return nil
}

func rewriteFile() error {
	file, err := os.OpenFile(global.FilePath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0766)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	defer file.Close()
	for i := range lines {
		text := lines[i] + "\n"
		if _, err := file.WriteString(text); err != nil {
			return fmt.Errorf(err.Error())
		}
	}
	lines = nil
	return nil
}

func delUrlService(serviceName string) error {
	file, err := os.Open(global.FilePath)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if rel := strings.Contains(line, serviceName);rel != true {
			lines = append(lines, line)
		}
	}
	if err := file.Close();	err != nil {
		return fmt.Errorf(err.Error())
	}

	log.Printf("INFO: service and url has been delete ---> [%v]", serviceName)
	if err := rewriteFile();err != nil {
		return fmt.Errorf(err.Error())
	}
	return nil
}


func RewriteUrl(c *gin.Context){
	f := From{
		Url:   c.PostForm("url"),
		UName: c.PostForm("u-name"),
		Notes: c.PostForm("txt"),
	}
	if f.Url == "" && f.UName == ""{
		log.Println("WARN: add element is nil.")
	}else{
		if f.Notes == "" {
			f.Notes = "无备注"
		}
		err := writeFile(f.UName, f.Url, f.Notes)
		if err != nil {
			log.Println(err)
		}
	}
	c.Redirect(http.StatusFound, "/url/index")
}


func DelUrl(c *gin.Context){
	n := Name{
		Name: c.PostForm("u-name"),
	}
	if n.Name == ""{
		log.Println("WARN: del element is nil.")
	}else {
		err := delUrlService(n.Name)
		if err != nil {
			log.Println(err.Error())
		}
	}
	c.Redirect(http.StatusFound, "/url/index")
}

type UrlInfo struct {
	SrcUrlName string `json:"srcUrlName"`
	SrcUrl     string `json:"srcUrl"`
	SrcNotes   string `json:"srcNotes"`
	UrlName    string `json:"urlName"`
	Url        string `json:"url"`
	Notes      string `json:"notes"`
}

func UpdateUrlInfo(c *gin.Context) {
	data, err := c.GetRawData()
	if err != nil {
		log.Printf("Error getting raw data: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	var body UrlInfo
	err = json.Unmarshal(data, &body)
	if err != nil {
		log.Printf("Error unmarshaling JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	file, err := os.Open(global.FilePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}
	scanner := bufio.NewScanner(file)
	srcText := body.SrcUrlName + " " + body.SrcUrl + " " + body.SrcNotes
	text := body.UrlName + " " + body.Url + " " + body.Notes
	for scanner.Scan() {
		line := scanner.Text()
		if line == srcText{
			line = text
		}
		lines = append(lines, line)
	}

	if err := file.Close();	err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	if err := rewriteFile();err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
	})
}