package route

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"guide/core"
	"guide/global"
	"guide/service"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)


func InitRoute() *gin.Engine {
	gin.SetMode("release")
	r := gin.Default()
	r.SetFuncMap(template.FuncMap{
		"checkFileTailStr": core.CheckFileTailStr,
	})
	r.Static("/sta","static")
	r.LoadHTMLGlob("templates/*.tmpl")

	r.NoRoute(core.SysIpWhitelist(global.IsStartWhitelist), func(c *gin.Context) {
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

	r.GET("/readme",core.SysIpWhitelist(global.IsStartWhitelist), func(c *gin.Context){
		c.HTML(http.StatusOK, "readme.tmpl",gin.H{})
	})

	url := r.Group("/url")
		url.GET("/index", core.SysIpWhitelist(global.IsStartWhitelist),func(c *gin.Context) {
			relStr := core.ShowUrl()
			c.HTML(http.StatusOK, "url.tmpl", gin.H{
				"UrlPic": relStr,
			})
		})
		url.POST("/upload", core.SysIpWhitelist(global.IsStartWhitelist),service.RewriteUrl)
		url.POST("/del", core.SysIpWhitelist(global.IsStartWhitelist),service.DelUrl)

	file := r.Group("/file")
		file.POST("/upload", core.SysIpWhitelist(global.IsStartWhitelist),service.UploadData)
		file.POST("/create", core.SysIpWhitelist(global.IsStartWhitelist),service.CreateDir)
		file.POST("/file/create", core.SysIpWhitelist(global.IsStartWhitelist), service.CreateFile)
		file.POST("/delete", core.SysIpWhitelist(global.IsStartWhitelist),service.DeleteDirAndFile)
		file.POST("/ys", core.SysIpWhitelist(global.IsStartWhitelist), service.CompressZipTar)
		file.POST("/jy", core.SysIpWhitelist(global.IsStartWhitelist), service.DecompressionZipTar)
		file.GET("/cat", core.SysIpWhitelist(global.IsStartWhitelist), service.CatFile)
		file.POST("/edit", core.SysIpWhitelist(global.IsStartWhitelist), service.UpdateFile)

	cron := r.Group("/cron")
		cron.GET("/index", core.SysIpWhitelist(global.IsStartWhitelist),func(c *gin.Context) {
			c.HTML(http.StatusOK, "cron.tmpl", gin.H{
			})
		})
		cron.GET("/list",core.SysIpWhitelist(global.IsStartWhitelist), service.ShowCron)
		cron.POST("/cfg", core.SysIpWhitelist(global.IsStartWhitelist), service.CoutomCron)
		cron.POST("/delete", core.SysIpWhitelist(global.IsStartWhitelist), service.DelCron)

	svc := r.Group("/svc")
		svc.GET("/index", core.SysIpWhitelist(global.IsStartWhitelist),func(c *gin.Context) {
			c.HTML(http.StatusOK, "protools.tmpl", gin.H{
			})
		})
		svc.POST("/cfg", core.SysIpWhitelist(global.IsStartWhitelist), service.SvcCfg)
		svc.POST("/delete", core.SysIpWhitelist(global.IsStartWhitelist), service.DeleteSvc)
		svc.GET("/list", core.SysIpWhitelist(global.IsStartWhitelist), service.ShowSvcCfg)

	pwd := r.Group("/pwd")
		pwd.GET("/index", core.SysIpWhitelist(global.IsStartWhitelist),core.PasswdAdminWhitelist(), service.PwdIndex)
		pwd.GET("/list", core.SysIpWhitelist(global.IsStartWhitelist),core.PasswdAdminWhitelist(), service.ShowPwdList)
		pwd.GET("/bak", core.SysIpWhitelist(global.IsStartWhitelist),core.PasswdAdminWhitelist(), service.UserPwdBackup)
		pwd.POST("/cfg",core.SysIpWhitelist(global.IsStartWhitelist),core.PasswdAdminWhitelist(), service.SavePwdToDb)
		pwd.POST("/cat", core.SysIpWhitelist(global.IsStartWhitelist),core.PasswdAdminWhitelist(), service.CatPwd)
		pwd.POST("/delete",core.SysIpWhitelist(global.IsStartWhitelist),core.PasswdAdminWhitelist(),service.DelUP)

	return r
}
