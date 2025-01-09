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
		pwd, err := core.PasswordDecrypt(user[0].Password, global.NowKey)
		if err != nil {
			log.Println(err)
			c.HTML(http.StatusOK, "login.tmpl", gin.H{
				"error": err,
			})
			return
		}
		if f.User == user[0].UserName && f.Password == pwd {
			cookie := http.Cookie{Name: "user", Value: f.User, MaxAge: 408000}
			http.SetCookie(c.Writer, &cookie)
			c.Redirect(http.StatusMovedPermanently, "/url/index")
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
	pwd, err := core.PasswordEncryption(userPassword, global.NowKey)
	if err != nil {log.Println(err)}
	if err := core.InsertUser(userName, pwd); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadGateway,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
	})
}

// ***** not test function
func UpdatePwd(c *gin.Context) {
	// data {
	//		"key": "xxxxxxxxxxxxxxxxxxxxxxxx", 24字节
	//		"userName": "xiaomi",
	//		"password": "xiaomi.666"
	//	}

	data, _ := c.GetRawData()
	var body map[string]string
	_ = json.Unmarshal(data, &body)
	key := body["key"]
	userName := body["userName"]
	userPassword := body["password"]

	// Constrain based on global key
	if key != global.NowKey {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadGateway,
			"msg":  "error,illegal request.",
		})
		return
	} else {
		pwd, err := core.PasswordEncryption(userPassword, global.NowKey)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusBadGateway,
			})
			return
		}

		if err := core.UpdateUserPwd(pwd, userName); err != nil {
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
	osType := core.ShowSys()
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
		return
	default:
		log.Printf("%s", "WARN: unsupported operating system.")
		return
	}

}
