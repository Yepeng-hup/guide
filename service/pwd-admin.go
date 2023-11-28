package service

import (
	"github.com/gin-gonic/gin"
	"guide/core"
	"guide/global"
	"log"
	"net/http"
	"strings"
)

type (
	PwdFrom struct {
		Svc string
		User string
		Passwd string
		Notes string
	}


)


func PwdIndex(c *gin.Context){
	c.HTML(http.StatusOK, "pwd.tmpl", gin.H{})
}


func SavePwdToDb(c *gin.Context){
	f := &PwdFrom{
		Svc: c.PostForm("svcName"),
		User: c.PostForm("loginName"),
		Passwd: c.PostForm("loginPwd"),
		Notes: c.PostForm("pwdNotes"),
	}

	if f.Svc == ""|| f.Passwd == ""|| f.User == ""|| f.Notes == "" {
		log.Println("WARN: add element is nil.")
		c.Redirect(http.StatusFound, "/pwd/index")
		return
	}

	encryptionPassword, err := core.PasswordEncryption(f.Passwd, global.NowKey)
	if err != nil {
		log.Printf("ERROR: %s", err)
		return
	}

	if err := core.InsertUserPwd(f.Svc, f.User, encryptionPassword, f.Notes); err != nil {
		log.Println(err)
		c.Redirect(http.StatusFound, "/pwd/index")
		return
	}
	c.Redirect(http.StatusFound, "/pwd/index")
	return
}


func ShowPwdList(c *gin.Context){
	// show all user and passwd
	list, err := core.SelectUserPwd()
	if err != nil {
		log.Println(err.Error())
		c.Redirect(http.StatusFound, "/pwd/index")
		return
	}
	c.HTML(http.StatusOK, "pwdcat.tmpl", gin.H{
			"userPwdList": list,
	})
}


func CatPwd(c *gin.Context){
	encryptionPassword := c.PostForm("pwd")
	ePasswordList := strings.Fields(encryptionPassword)
	dpwd, err := core.PasswordDecrypt(ePasswordList[2], global.NowKey)
	if err != nil {
		log.Printf("ERROR: %s", err)
		c.Redirect(http.StatusFound, "/pwd/list")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"dPassword": dpwd,
	})
	return
}


func DelUP(c *gin.Context){

}


func UserPwdBackup(c *gin.Context){

}