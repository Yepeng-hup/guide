package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"guide/core"
	"guide/core/cmd"
	"guide/global"
	"log"
	"net/http"
)

type (
	LoginFrom struct {
		User     string
		Password string
	}
)

func Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.tmpl", gin.H{})
}

func LoginCk(c *gin.Context) {
	f := LoginFrom{
		User:     c.PostForm("Name"),
		Password: c.PostForm("Password"),
	}
	user, err := core.SelectUser(fmt.Sprintf("SELECT id,userName,newUserDate,password FROM user WHERE userName = \"%s\"", f.User))
	if err != nil {
		log.Println(err.Error())
		return
	}
	if len(user) >= 1 {
		if f.User == user[0].UserName && f.Password == user[0].Password {
			cookie := http.Cookie{Name: "user", Value: f.User, MaxAge: 108000}
			http.SetCookie(c.Writer, &cookie)
			c.Redirect(http.StatusMovedPermanently, "/user/index")
			return
		}
	}
	c.HTML(http.StatusOK, "login.tmpl", gin.H{
		"error": "user or password input error",
	})
}

func UserAdmin(c *gin.Context) {
	userList, err := core.SelectUser("SELECT id,userName,newUserDate,password FROM user")
	if err != nil {
		log.Println(err.Error())
		return
	}
	c.HTML(http.StatusOK, "user.tmpl", gin.H{
		"userList": userList,
	})
}

func CreateUser(c *gin.Context) {
	data, _ := c.GetRawData()
	var body map[string]string
	_ = json.Unmarshal(data, &body)
	userName := body["userName"]
	userPassword := body["password"]
	err := core.InsertUser(userName, userPassword)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadGateway,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
	})
}

// ***** 待测试函数
func UpdatePwd(c *gin.Context) {
	// data {"key": xxxx, "userName": "xx", "password": "xxxx"}

	data, _ := c.GetRawData()
	var body map[string]string
	_ = json.Unmarshal(data, &body)
	key := body["key"]
	userName := body["userName"]
	userPassword := body["password"]
	fmt.Println(userName, userPassword, key)

	// 根据全局key做约束
	if key != global.NowKey {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadGateway,
			"msg":  "error,illegal request.",
		})
	} else {
		if err := core.UpdateUserPwd(userPassword, userName); err != nil {
			log.Println(err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
		})
	}
}

func DeleteUser(c *gin.Context) {
	data, _ := c.GetRawData()
	var body map[string]string
	_ = json.Unmarshal(data, &body)
	userName := body["userName"]
	user, err := c.Cookie("user")
	if err != nil {
		log.Println(err.Error())
		return
	}

	if user != "admin" {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadGateway,
		})
		log.Println("ERROR: This login user does not have permission --> ", user)
		return
	}
	if err := core.DeleteUser(userName); err != nil {
		log.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
	})
}

func UpdateUserInfo(c *gin.Context) {
	data, _ := c.GetRawData()
	var body map[string]string
	_ = json.Unmarshal(data, &body)
	userId := body["userId"]
	userName := body["userName"]
	newUserDate := body["newUserDate"]
	user, err := c.Cookie("user")
	if err != nil {
		log.Println(err.Error())
		return
	}

	if user != "admin" {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadGateway,
		})
		log.Println("ERROR: This login user does not have permission --> ", user)
		return
	}
	if err := core.UpdateUser(userName, newUserDate, userId); err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadGateway,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
	})
}

func RebootHost(c *gin.Context) {
	osType := cmd.ShowSys()
	switch osType {
	case "linux":
		cmdCode := "reboot"
		if err := cmd.UseCmd(cmdCode); err != nil {
			log.Println("ERROR: reboot host fail,", err)
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusBadGateway,
			})
			return
		}
		//log.Println("INFO: reboot host ok.")
		//c.JSON(http.StatusOK, gin.H{
		//	"code": http.StatusOK,
		//})
		return
	case "windows":
		cmdCode := "shutdown.exe -r -f -t 0"
		if err := cmd.UseCmd(cmdCode); err != nil {
			log.Println("ERROR: reboot host fail,", err)
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusBadGateway,
			})
			return
		}
		//log.Println("INFO: reboot host ok.")
		//c.JSON(http.StatusOK, gin.H{
		//	"code": http.StatusOK,
		//})
		return
	default:
		log.Printf("%s", "WARN: unsupported operating system.")
		return
	}

}
