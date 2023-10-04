package service

import (
	"github.com/gin-gonic/gin"
	"guide/global"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)


func CutDirAndFile(c *gin.Context, fullPath *string) {
	files, _ := ioutil.ReadDir(*fullPath)
	dirList := make([]DirectoryAnchor, 0)
	fileList := make([]FileAnchor, 0)
	dirList = append(dirList, DirectoryAnchor{
		DirectoryName: "..",
		Href: strings.TrimRight(c.Request.URL.Path, "/") + "/..",
	})
	for _, file := range files {
		href := strings.ReplaceAll(c.Request.URL.Path +"/"+ file.Name(), "//", "/")
		if file.IsDir() {
			dirList = append(dirList, DirectoryAnchor{
				DirectoryName: file.Name(),
				Href:          href,
				Size: file.Size(),
				Time: file.ModTime(),
				Power: file.Mode(),
			})

		} else {
			fileList = append(fileList, FileAnchor{
				FileName: file.Name(),
				Href:     href,
				Size: file.Size()/1024/1024,
				Time: file.ModTime(),
				Power: file.Mode(),
			})
		}
	}
	c.HTML(http.StatusOK, "file.tmpl", gin.H{
		"dirList": dirList,
		"fileList": fileList,
		"currentDir": c.Request.URL.Path,
	})
}


func UploadData(c *gin.Context) {
	file, _ := c.FormFile("file")
	if file == nil {
		c.Redirect(http.StatusFound, "/")
		return
	}
	filename := filepath.Base(file.Filename)
	if strings.Contains(c.Query("path"), "..") {
		c.Redirect(http.StatusFound, "/")
		return
	}
	savePath := filepath.Join(global.SaveDataDir, c.Query("path"), filename)
	err := c.SaveUploadedFile(file, savePath)
	if err != nil {
		log.Println("ERROR: file save fail,", err.Error())
		c.Redirect(http.StatusFound, "/")
		return
	}
	c.Redirect(http.StatusFound, c.Query("path"))
}


func DownloadData(c *gin.Context, absolutePath *string) {
	c.File(*absolutePath)
}


func CreateDir(c *gin.Context){
	dirName := c.PostForm("name")
	err := os.Mkdir(dirName, 0755)
	if err != nil {
		log.Println("ERROR: create dir and file fail.", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"message": "目录创建失败.",
		})
	}
	log.Printf("INFO: create dir and file success ---> [%v].", dirName)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "目录创建成功.",
	})
}


func DeleteDirAndFile(c *gin.Context){

}


