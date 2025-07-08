package route

import (
	"fmt"
	"guide/core"
	"guide/core/mon"
	"guide/global"
	"guide/service"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	mlog = core.Wlogs{}
)

func InitRoute() *gin.Engine {
	gin.SetMode(core.Cfg.DebugMode)
	r := gin.Default()
	r.SetFuncMap(template.FuncMap{
		"checkFileTailStr": core.CheckFileTailStr,
	})
	r.Static("/sta", "static")
	r.LoadHTMLGlob("templates/*.tmpl")
	r.Use(mon.PromMonMiddleware())

	r.NoRoute(core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), func(c *gin.Context) {
		fullPath := filepath.Join(core.Cfg.FileDataDir, c.Request.URL.Path)
		if c.Request.URL.Path == "/root" {
			var body map[string]string
			if err := c.BindJSON(&body); err != nil {
				mlog.Error(err.Error())
				c.JSON(http.StatusOK, gin.H{
					"code": http.StatusInternalServerError,
				})
				return
			}
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

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/login", service.Login)
	r.POST("/loginck", service.LoginCk)
	r.GET("/logout", func(c *gin.Context) {
		cookie := http.Cookie{Name: "user", MaxAge: -1}
		http.SetCookie(c.Writer, &cookie)
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
		})
	})
	
	r.GET("/home", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/home"), service.HomeIndex)
	r.POST("/home/update/pwd", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), service.UpdateHomePwd)
	r.GET("/reboot", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/reboot"), service.RebootHost)
	r.GET("/home/data",core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/home/data"), service.LoginDataSource)

	user := r.Group("/user")
	user.GET("/index", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/user/index"), service.UserAdmin)
	user.POST("/create", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/user/create"), service.CreateUser)
	user.POST("/update/pwd", service.UpdatePwd)
	user.DELETE("/delete", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/user/delete"), service.DeleteUser)
	// user.POST("/update/info", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/user/update/info"), service.UpdateUserInfo)
	user.GET("/role/index", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/user/role/index"), service.RoleAdmin)
	user.POST("/role/create", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/user/role/create"), service.CreateRole)
	user.GET("/role/select", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/user/role/select"), service.SelectRole)
	user.POST("/role/per/select", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/user/role/per/select"), service.SelectRolePermission)
	user.DELETE("/role/delete", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/user/role/delete"), service.DelRole)
	user.DELETE("/role/per/delete", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/user/role/per/delete"), service.DelRolePermissionRoute)
	user.GET("/role/per/admin/select", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/user/role/per/admin/select"), service.SelectAdminRolePermission)

	url := r.Group("/url")
	url.GET("/index", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/url/index"), service.ShowDbUrl)
	url.POST("/show", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/url/show"), service.ShowTypeUrl)
	url.POST("/upload", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/url/upload"), service.RewriteUrl2)
	url.POST("/del", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/url/del"), service.DelUrl2)
	url.POST("/update", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/url/update"), service.UpdateUrlInfo2)
	url.POST("/type/create", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/url/type/create"), service.CreateType)
	url.GET("/type/list", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/url/type/list"), service.ShowUrlType)
	url.POST("/type/del", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/url/type/del"), service.DelUrlType)

	file := r.Group("/file")
	file.POST("/upload", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/file/upload"), service.UploadData)
	file.POST("/create", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/file/create"), service.CreateDir)
	file.POST("/file/create", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/file/file/create"), service.CreateFile)
	file.POST("/delete", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/file/delete"), service.DeleteDirAndFile)
	file.POST("/ys", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/file/ys"), service.CompressZipTar)
	file.POST("/jy", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/file/jy"), service.DecompressionZipTar)
	file.GET("/cat", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/file/cat"), service.CatFile)
	file.POST("/edit", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/file/edit"), service.UpdateFile)
	//file.GET("/download", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), service.DownloadFile)
	file.GET("/hs", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/file/hs"), service.ShowRecycle)
	file.POST("/hs/delete", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/file/hs/delete"), service.DeleteRecycleFile)
	file.POST("/ss", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/file/ss"), service.FileSearch)
	file.POST("/dir/check", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/file/dir/check"), service.RootDirCheck)

	cron := r.Group("/cron")
	cron.GET("/index", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/cron/index"), func(c *gin.Context) {
		c.HTML(http.StatusOK, "cron.tmpl", gin.H{})
	})
	cron.GET("/list", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/cron/list"), service.ShowCron)
	cron.POST("/cfg", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/cron/cfg"), service.CoutomCron)
	cron.POST("/delete", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/cron/delete"), service.DelCron)

	svc := r.Group("/svc")
	svc.GET("/index", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/svc/index"), func(c *gin.Context) {
		c.HTML(http.StatusOK, "protools.tmpl", gin.H{})
	})
	svc.POST("/cfg", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/svc/cfg"), service.SvcCfg)
	svc.POST("/delete", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/svc/delete"), service.DeleteSvc)
	svc.GET("/list", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/svc/list"), service.ShowSvcCfg)

	pwd := r.Group("/pwd")
	pwd.GET("/index", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.PasswdAdminWhitelist(), core.CookieCheck(), core.PermissionCheck("/pwd/index"), service.PwdIndex)
	pwd.GET("/list", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.PasswdAdminWhitelist(), core.CookieCheck(), core.PermissionCheck("/pwd/list"), service.ShowPwdList)
	pwd.GET("/bak", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.PasswdAdminWhitelist(), core.CookieCheck(), core.PermissionCheck("/pwd/bak"), service.UserPwdBackup)
	pwd.POST("/cfg", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.PasswdAdminWhitelist(), core.CookieCheck(), core.PermissionCheck("/pwd/cfg"), service.SavePwdToDb)
	pwd.POST("/cat", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.PasswdAdminWhitelist(), core.CookieCheck(), core.PermissionCheck("/pwd/cat"), service.CatPwd)
	pwd.POST("/delete", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.PasswdAdminWhitelist(), core.CookieCheck(), core.PermissionCheck("/pwd/delete"), service.DelUP)

	log := r.Group("/log")
	log.GET("/index", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/log/index"), func(c *gin.Context) {
		c.HTML(http.StatusOK, "syslog.tmpl", gin.H{})
	})
	// log.POST("/w", service.InsertLogToDb)
	log.POST("/r", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/log/r"), service.ShowLog)
	log.POST("d", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/log/d"), service.DeleteLimitLog)
	log.POST("/bit", service.FluentBit)

	security := r.Group("/security")
	security.GET("/index", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/security/index"), service.SecurityIndex)
	security.POST("/black/add", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/security/black/add"), service.AddBlacklistIp)
	security.DELETE("/black/mv", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/security/black/mv"), service.MoveBlacklistIp)
	security.GET("/black/show", core.SysIpWhitelist(core.Cfg.StartWhiteList), core.CookieCheck(), core.PermissionCheck("/security/black/show"), service.ShowDbBlacklistIp)
	return r
}
