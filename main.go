package main

import (
	"bufio"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	filePath = "text/url.txt"
	lines    []string
)

type (
	from struct {
		url   string
		uName string
		notes string
	}
	name struct {
		Name string
	}

	News struct {
		UName string
		Url string
		Notes string
	}

)

func writeFile(fileName, url, notes string) error {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	defer file.Close()
	text := fileName + " " + url + " " + notes+"\n"
	if _, err := file.WriteString(text); err != nil {
		return fmt.Errorf(err.Error())
	}
	return nil
}


func str()([]News,){
	var structSlice []News
	file, err := os.Open(filePath)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		slice := strings.Split(line, " ")
		txtStruct := News{
			UName:   slice[0],
			Url: slice[1],
			Notes: slice[2],
		}
		//fmt.Println(myStruct.uName)
		structSlice = append(structSlice, txtStruct)
	}

	if err := scanner.Err(); err != nil {
		log.Println(err.Error())
		return nil
	}
	return structSlice
}


func rewriteFile() error {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0766)
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
	file, err := os.Open(filePath)
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

func main() {
	gin.SetMode("release")
	r := gin.Default()
	r.Static("/sta","static")
	r.LoadHTMLGlob("templates/*")

	r.GET("/index", func(c *gin.Context) {
		relStr := str()
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"UrlPic": relStr,
		})
	})

	r.POST("/upload", func(c *gin.Context) {
		f := from{
			url:   c.PostForm("url"),
			uName: c.PostForm("u-name"),
			notes: c.PostForm("txt"),
		}
		if f.url == "" && f.uName == ""{
			log.Println("add element is nil.")
		}else{
			if f.notes == "" {
				f.notes = "无备注"
			}
			err := writeFile(f.uName, f.url, f.notes)
			if err != nil {
				log.Println(err)
			}
		}
		c.Redirect(http.StatusFound, "/index")
	})

	r.POST("/del", func(c *gin.Context) {
		n := name{
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
		c.Redirect(http.StatusFound, "/index")
	})
	err := godotenv.Load("guide.env")
	if err != nil {
		log.Fatal(err)
	}
	host := os.Getenv("GUIDE_HOST")
	port := os.Getenv("GUIDE_PORT")
	log.Println("start ok ---> Listening and serving HTTP on "+host+":"+port+"/index")
	if err := r.Run(host+":"+port); err != nil {
		log.Println("error start fail", err)
	}
}
