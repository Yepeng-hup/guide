package route

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"guide/core"
	"guide/global"
	"guide/service"
	"net/http"
	"os"
	"path/filepath"
	_ "path/filepath"
)


func InitRoute() *gin.Engine {
	gin.SetMode(global.GinDebug)
	r := gin.Default()
	r.Static("/sta","static")
	r.LoadHTMLGlob("templates/*")

	r.NoRoute(func(c *gin.Context) {
		fullPath := filepath.Join(global.SaveDataDir, c.Request.URL.Path)
		fileInfo, err := os.Stat(fullPath)
		if err != nil {
			fmt.Println(err.Error())
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		if fileInfo.IsDir() {
			service.CutDirAndFile(c, &fullPath)
		} else {
			service.DownloadData(c, &fullPath)
		}

	})

	r.GET("/readme", func(c *gin.Context){
		c.HTML(http.StatusOK, "readme.tmpl",gin.H{})
	})

	url := r.Group("/url")
		url.GET("/index", func(c *gin.Context) {
			relStr := core.ShowUrl()
			c.HTML(http.StatusOK, "url.tmpl", gin.H{
				"UrlPic": relStr,
			})
		})
		url.POST("/upload", service.RewriteUrl)
		url.POST("/del", service.DelUrl)

	file := r.Group("/file")
		file.POST("/upload", service.UploadData)
		file.POST("/create", service.CreateDir)

	return r
}
