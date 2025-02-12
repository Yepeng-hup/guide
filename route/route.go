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
	gin.SetMode(global.Mode)
	r := gin.Default()
	r.SetFuncMap(template.FuncMap{
		"checkFileTailStr": core.CheckFileTailStr,
	})
	r.Static("/sta", "static")
	r.LoadHTMLGlob("templates/*.tmpl")

	r.NoRoute(core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), func(c *gin.Context) {
		fullPath := filepath.Join(global.SaveDataDir, c.Request.URL.Path)
		if c.Request.URL.Path == "/root" {
			var body map[string]string
			c.BindJSON(&body)
			password := body["password"]
			jmDirPwd, _ := core.PasswordEncryption(global.DirAccessPwd, global.NowKey)
			jmPwd, _ := core.PasswordEncryption(password, global.NowKey)
			dirPwd, _ := core.PasswordDecrypt(jmDirPwd, global.NowKey)
			pwd, _ := core.PasswordDecrypt(jmPwd, global.NowKey)

			if pwd != dirPwd {
				c.JSON(http.StatusOK, gin.H{
					"code": http.StatusInternalServerError,
				})
				return
			}
			//c.Redirect(http.StatusMovedPermanently, "/root/?password="+password)
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusOK,
				"pwd":  jmPwd,
			})
			return
		}

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
	url.GET("/index", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), service.ShowDbUrl)
	url.POST("/show", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), service.ShowTypeUrl)
	url.POST("/upload", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), service.RewriteUrl2)
	url.POST("/del", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), service.DelUrl2)
	url.POST("/update", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), service.UpdateUrlInfo2)
	url.POST("/type/create", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), service.CreateType)
	url.GET("/type/list", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), service.ShowUrlType)
	url.POST("/type/del", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), service.DelUrlType)

	file := r.Group("/file")
	file.POST("/upload", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), service.UploadData)
	file.POST("/create", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), service.CreateDir)
	file.POST("/file/create", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), service.CreateFile)
	file.POST("/delete", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), service.DeleteDirAndFile)
	file.POST("/ys", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), service.CompressZipTar)
	file.POST("/jy", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), service.DecompressionZipTar)
	file.GET("/cat", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), service.CatFile)
	file.POST("/edit", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), service.UpdateFile)
	file.GET("/hs", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), service.ShowRecycle)
	file.POST("/hs/delete", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), service.DeleteRecycleFile)
	file.POST("/ss", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), service.FileSearch)
	file.POST("/dir/check", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), service.RootDirCheck)

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
	// log.POST("/w", service.InsertLogToDb)
	log.POST("/r", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), service.ShowLog)
	log.POST("d", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), service.DeleteLimitLog)
	log.POST("/bit", service.FluentBit)

	security := r.Group("/security")
	security.GET("/index", core.SysIpWhitelist(global.IsStartWhitelist), core.CookieCheck(), service.SecurityIndex)

	return r
}
