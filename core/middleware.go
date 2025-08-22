package core

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func ifSysIpInWhitelist(ip *string) bool {
	whiteIpList := strings.Split(Cfg.WhiteList, ",")
	for _, allowedIP := range whiteIpList {
		if allowedIP == *ip {
			return true
		}
	}
	return false
}

// admin ip num 1
func ifPwdIpWhitelist(ip *string) bool {
	whitelistIp := Cfg.PasswdAdminWhiteList
	if whitelistIp == *ip {
		return true
	}
	return false
}

func getClientIP(r *http.Request) (string, error) {
	// show client ipaddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", fmt.Errorf("ERROR: show client ip fail,%s", err.Error())
	}
	return ip, nil
}

func SysIpWhitelist(onOff string) gin.HandlerFunc {
	//if start whitelist
	if onOff != "true" {
		return func(c *gin.Context) {
			c.Next()
		}
	}
	return func(c *gin.Context) {
		clientIP, err := getClientIP(c.Request)
		if err != nil {
			log.Println(err.Error())
			return
		}
		// check clientIp in whitelist
		if !ifSysIpInWhitelist(&clientIP) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "The ip does not have permission, please contact the management personnel to activate the whitelist",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func PasswdAdminWhitelist() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP, err := getClientIP(c.Request)
		if err != nil {
			log.Println(err.Error())
			return
		}
		if !ifPwdIpWhitelist(&clientIP) {
			log.Printf("ERROR: %s host does not allow access to password management.", clientIP)
			c.Redirect(http.StatusFound, "/")
			return
		}
		c.Next()
	}
}

func CookieCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := c.Cookie("user"); err != nil {
			mlog.Error(fmt.Sprintf("not login,%s", err.Error()))
			c.Redirect(http.StatusMovedPermanently, "/login")
			return
		}
		return
	}
}

func discernPermission(roleName string) []string {
	var urlSlices = make([]string, 0)
	sql := fmt.Sprintf("SELECT * from roles_permission WHERE roleName = \"%s\"", roleName)
	routeList, _ := SelectRolePermission(sql)
	if len(routeList) <= 0 {
		mlog.Error("show role permission fail. LIST not index, index nil.")
		return nil
	}

	for _, v := range routeList {
		urlSlices = append(urlSlices, v.Permission)
	}
	return urlSlices
}

func IfPermission(permissionSlices []string, url string) bool {
	for _, v := range permissionSlices {
		if v == url {
			return true
		}
	}
	return false
}

func PermissionCheck(routePath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		user, err := c.Cookie("user")
		if err != nil {
			mlog.Error(fmt.Sprintf("not use cookie and cookie config fail, %s", err.Error()))
			if method == "GET" {
				c.Redirect(http.StatusMovedPermanently, "/")
				return
			} else {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"code": http.StatusForbidden,
					"msg":  "You do not have permission to access, please contact the administrator.",
				})
				return
			}
		}

		sql := fmt.Sprintf("select * from user_roles where userName = \"%s\"", user)
		r, _ := SelectUserAndRole(sql)
		if len(r) <= 0 {
			mlog.Error(fmt.Sprintf("not user [%s] is role.", user))
			return
		}

		urlSlices := discernPermission(r[0].RoleName)
		boolRel := IfPermission(urlSlices, routePath)
		if boolRel {
			c.Next()
		} else {
			//
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code": http.StatusForbidden,
				"msg":  "You do not have permission to access, please contact the administrator.",
			})
			return
		}
		c.Next()
	}
}
