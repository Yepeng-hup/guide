package main

import (
	"bufio"
	"fmt"
	"github.com/gin-gonic/gin"
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
	}
	name struct {
		Name string
	}
)

func writeFile(fileName, url string) error {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	defer file.Close()
	text := fileName + " " + url + "\n"
	if _, err := file.WriteString(text); err != nil {
		return fmt.Errorf(err.Error())
	}
	return nil
}

func readFileLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func sliceUewMap(textSlice []string) map[string]string {
	imageMap := make(map[string]string)
	for _, s := range textSlice {
		parts := strings.Split(s, " ")
		imageMap[parts[0]] = parts[1]
	}
	return imageMap
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
	r.LoadHTMLGlob("templates/*")

	r.GET("/index", func(c *gin.Context) {
		rel, err := readFileLines(filePath)
		if err != nil {
			log.Println(err.Error())
		}
		relMap := sliceUewMap(rel)
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"UrlPic": relMap,
		})
	})

	r.POST("/upload", func(c *gin.Context) {
		f := from{
			url:   c.PostForm("url"),
			uName: c.PostForm("u-name"),
		}
		if f.url == "" && f.uName == ""{
			log.Println("add element is nil.")
		}else{
			err := writeFile(f.uName, f.url)
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
	log.Println("start ok ---> Listening and serving HTTP on 0.0.0.0:7878")
	if err := r.Run("0.0.0.0:7878"); err != nil {
		log.Println("error start fail", err)
	}
}
