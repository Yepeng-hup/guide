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
	r.Static("/sta", "static")
	r.LoadHTMLGlob("templates/*.tmpl")

	r.NoRoute(core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), func(c *gin.Context) {
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

	r.GET("/login", service.Login)
	r.POST("/login", service.LoginCk)
	r.GET("/logout", func(c *gin.Context) {
		cookie := http.Cookie{Name: "user", MaxAge: -1}
		http.SetCookie(c.Writer, &cookie)
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
		})
	})

	r.GET("/readme", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "readme.tmpl", gin.H{})
	})
	r.GET("/reboot", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), service.RebootHost)

	user := r.Group("/user")
	user.GET("/index", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), service.UserAdmin)
	user.POST("/create", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), service.CreateUser)
	user.POST("/update/pwd", service.UpdatePwd)
	user.DELETE("/delete", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), service.DeleteUser)
	user.POST("/update/info", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), service.UpdateUserInfo)

	url := r.Group("/url")
	url.GET("/index", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), func(c *gin.Context) {
		relStr := core.ShowUrl()
		c.HTML(http.StatusOK, "url.tmpl", gin.H{
			"UrlPic": relStr,
		})
	})
	url.POST("/upload", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), service.RewriteUrl)
	url.POST("/del", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), service.DelUrl)
	url.POST("/update", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), service.UpdateUrlInfo)

	file := r.Group("/file")
	file.POST("/upload", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), service.UploadData)
	file.POST("/create", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), service.CreateDir)
	file.POST("/file/create", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), service.CreateFile)
	file.POST("/delete", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), service.DeleteDirAndFile)
	file.POST("/ys", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), service.CompressZipTar)
	file.POST("/jy", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), service.DecompressionZipTar)
	file.GET("/cat", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), service.CatFile)
	file.POST("/edit", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), service.UpdateFile)

	cron := r.Group("/cron")
	cron.GET("/index", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "cron.tmpl", gin.H{})
	})
	cron.GET("/list", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), service.ShowCron)
	cron.POST("/cfg", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), service.CoutomCron)
	cron.POST("/delete", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), service.DelCron)

	svc := r.Group("/svc")
	svc.GET("/index", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "protools.tmpl", gin.H{})
	})
	svc.POST("/cfg", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), service.SvcCfg)
	svc.POST("/delete", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), service.DeleteSvc)
	svc.GET("/list", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), service.ShowSvcCfg)

	pwd := r.Group("/pwd")
	pwd.GET("/index", core.SysIpWhitelist(global.IsStartWhitelist), core.PasswdAdminWhitelist(), core.CookieCheck(), service.PwdIndex)
	pwd.GET("/list", core.SysIpWhitelist(global.IsStartWhitelist), core.PasswdAdminWhitelist(), core.CookieCheck(), service.ShowPwdList)
	pwd.GET("/bak", core.SysIpWhitelist(global.IsStartWhitelist), core.PasswdAdminWhitelist(), core.CookieCheck(), service.UserPwdBackup)
	pwd.POST("/cfg", core.SysIpWhitelist(global.IsStartWhitelist), core.PasswdAdminWhitelist(), core.CookieCheck(), service.SavePwdToDb)
	pwd.POST("/cat", core.SysIpWhitelist(global.IsStartWhitelist), core.PasswdAdminWhitelist(), core.CookieCheck(), service.CatPwd)
	pwd.POST("/delete", core.SysIpWhitelist(global.IsStartWhitelist), core.PasswdAdminWhitelist(), core.CookieCheck(), service.DelUP)

	log := r.Group("/log")
	log.GET("/index", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "syslog.tmpl", gin.H{})
	})
	log.POST("/w", service.InsertLogToDb)
	log.POST("/r", service.ShowLog)
	log.POST("d", service.DeleteLimitLog)

	security := r.Group("/security")
	security.GET("/index", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), service.SecurityIndex)

	return r
}
