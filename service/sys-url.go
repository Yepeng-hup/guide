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
		rel := strings.Contains(line, serviceName)
		if rel != true {
			lines = append(lines, line)
		}
	}
	err = file.Close()
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	log.Printf("service and url has been delete ---> [%v]", serviceName)
	err = rewriteFile()
	if err != nil {
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
		log.Println("add element is nil.")
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
		log.Println("del element is nil.")
	}else {
		err := delUrlService(n.Name)
		if err != nil {
			log.Println(err.Error())
		}
	}
	c.Redirect(http.StatusFound, "/url/index")
}