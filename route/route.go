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
	gin.SetMode("release")
	r := gin.Default()
	r.Static("/sta","static")
	r.LoadHTMLGlob("templates/*")

	r.NoRoute(core.IpWhitelistMiddleware(global.IsStartWhitelist), func(c *gin.Context) {
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

	r.GET("/readme",core.IpWhitelistMiddleware(global.IsStartWhitelist), func(c *gin.Context){
		c.HTML(http.StatusOK, "readme.tmpl",gin.H{})
	})

	url := r.Group("/url")
		url.GET("/index", core.IpWhitelistMiddleware(global.IsStartWhitelist),func(c *gin.Context) {
			relStr := core.ShowUrl()
			c.HTML(http.StatusOK, "url.tmpl", gin.H{
				"UrlPic": relStr,
			})
		})
		url.POST("/upload", core.IpWhitelistMiddleware(global.IsStartWhitelist),service.RewriteUrl)
		url.POST("/del", core.IpWhitelistMiddleware(global.IsStartWhitelist),service.DelUrl)

	file := r.Group("/file")
		file.POST("/upload", core.IpWhitelistMiddleware(global.IsStartWhitelist),service.UploadData)
		file.POST("/create", core.IpWhitelistMiddleware(global.IsStartWhitelist),service.CreateDir)
		file.POST("/delete", core.IpWhitelistMiddleware(global.IsStartWhitelist),service.DeleteDirAndFile)

	return r
}
