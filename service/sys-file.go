package service

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
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
		//i, s := core.BitConvert(file.Size())
		//fmt.Println(i, s)
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
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"msg": err.Error(),
		})
		return
	}

	filename := filepath.Base(file.Filename)
	if strings.Contains(c.Query("path"), "..") {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid path.",
		})
		return
	}

	savePath := filepath.Join("file", c.Query("path"), filename)
	err = c.SaveUploadedFile(file, savePath)
	if err != nil {
		//c.String(http.StatusInternalServerError, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
	}
	c.Redirect(http.StatusFound, "/")
}


func DownloadData(c *gin.Context, absolutePath *string) {
	c.File(*absolutePath)
}


