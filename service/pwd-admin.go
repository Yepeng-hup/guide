package service

import (
	"github.com/gin-gonic/gin"
	"guide/core"
	"guide/global"
	"net/http"
	"os"
	"strings"
)

type (
	PwdFrom struct {
		Svc    string
		User   string
		Passwd string
		Notes  string
	}
)

func PwdIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "pwd.tmpl", gin.H{})
}

func SavePwdToDb(c *gin.Context) {
	f := &PwdFrom{
		Svc:    c.PostForm("svcName"),
		User:   c.PostForm("loginName"),
		Passwd: c.PostForm("loginPwd"),
		Notes:  c.PostForm("pwdNotes"),
	}

	if f.Svc == "" || f.Passwd == "" || f.User == "" || f.Notes == "" {
		mlog.Warn("add element is nil.")
		c.Redirect(http.StatusFound, "/pwd/index")
		return
	}

	encryptionPassword, err := core.PasswordEncryption(f.Passwd, global.NowKey)
	if err != nil {
		mlog.Error(err.Error())
		return
	}

	if err := core.InsertUserPwd(f.Svc, f.User, encryptionPassword, f.Notes); err != nil {
		mlog.Error(err.Error())
		c.Redirect(http.StatusFound, "/pwd/index")
		return
	}
	c.Redirect(http.StatusFound, "/pwd/index")
}

func ShowPwdList(c *gin.Context) {
	// show all user and passwd
	list, err := core.SelectUserPwd()
	if err != nil {
		mlog.Error(err.Error())
		c.Redirect(http.StatusFound, "/pwd/index")
		return
	}
	c.HTML(http.StatusOK, "pwdcat.tmpl", gin.H{
		"userPwdList": list,
	})
}

func CatPwd(c *gin.Context) {
	encryptionPassword := c.PostForm("pwd")
	ePasswordList := strings.Fields(encryptionPassword)
	dpwd, err := core.PasswordDecrypt(ePasswordList[2], global.NowKey)
	if err != nil {
		mlog.Error(err.Error())
		c.Redirect(http.StatusFound, "/pwd/list")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":      http.StatusOK,
		"dPassword": dpwd,
	})
}

func DelUP(c *gin.Context) {
	str := c.PostForm("pwd")
	list := strings.Fields(str)
	err := core.DeleteUserPwd(list[0])
	if err != nil {
		mlog.Error(err.Error())
		c.Redirect(http.StatusFound, "/pwd/list")
		return
	}
	c.Redirect(http.StatusFound, "/pwd/list")
}

func UserPwdBackup(c *gin.Context) {
	list, err := core.SelectUserPwd()
	if err != nil {
		mlog.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "error: db select fail.",
		})
		return
	}
	p := core.BackupPrefix()
	file, err := os.OpenFile("backup_"+p+".txt", os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0766)
	if err != nil {
		mlog.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "error: create file fail.",
		})
		return
	}
	defer file.Close()
	for _, v := range list {
		text := v.ServiceName + "  " + v.User + "  " + v.Passwd + "\n"
		if _, err := file.WriteString(text); err != nil {
			mlog.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "error: bakup data write file fail.",
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "data backup success.",
	})
}
