package core

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"guide/global"
	"log"
	"net"
	"net/http"
)


func ifSysIpInWhitelist(ip *string) bool {
	whitelist := global.IpList()
	for _, allowedIP := range whitelist {
		if allowedIP == *ip {
			return true
		}
	}
	return false
}


func ifPwdIpWhitelist(ip *string) bool {
	whitelistIp := global.PasswdAdminWhitelist
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


func PasswdAdminWhitelist()gin.HandlerFunc{
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
			log.Println("not login,", err.Error())
			c.Redirect(http.StatusMovedPermanently, "/login")
			return
		}
		//fmt.Println("================= cookie ok.")
		return
	}
}
