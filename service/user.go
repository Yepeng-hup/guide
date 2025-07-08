package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"guide/core"
	"guide/core/cmd"
	"guide/global"
	"net/http"
)

type (
	LoginFrom struct {
		User     string `json:"user"`
		Password string `json:"password"`
	}

	User struct {
		UserName string `json:"userName"`
		RoleName string `json:"roleName"`
		Password string `json:"password"`
	}
)

var (
	mlog = core.Wlogs{}
)

func Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.tmpl", gin.H{})
}

func LoginCk(c *gin.Context) {
	var f LoginFrom

	if err := c.ShouldBindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "json analysis parameter error.",
		})
		return
	}

	user, err := core.SelectUser(fmt.Sprintf("SELECT id,userName,newUserDate,password FROM user WHERE userName = \"%s\"", f.User))
	if err != nil {
		mlog.Error(err.Error())
		return
	}
	if len(user) >= 1 {
		pwd, err := core.PasswordDecrypt(user[0].Password, global.NowKey)
		if err != nil {
			mlog.Error(err.Error())
			c.HTML(http.StatusOK, "login.tmpl", gin.H{
				"error": err,
			})
			return
		}
		if f.User == user[0].UserName && f.Password == pwd {
			permissionList, err := core.SelectUserPermission(f.User)
			if err != nil {
				mlog.Error(err.Error())
				c.JSON(http.StatusOK, gin.H{
					"code": http.StatusInternalServerError,
					"msg":  err.Error(),
				})
				return
			}

			cookie := http.Cookie{Name: "user", Value: f.User, MaxAge: 408000}
			http.SetCookie(c.Writer, &cookie)
			mlog.Info(fmt.Sprintf("user -> [%s] login success.", f.User))
			if err := core.InsertActLoginUser(f.User); err != nil {
				mlog.Error(err.Error())
			}

			c.JSON(http.StatusOK, gin.H{
				"code":       http.StatusOK,
				"permission": permissionList,
				"loginUser":  f.User,
			})
			return
		}
	}

	mlog.Error(fmt.Sprintf("user -> [%s] login fail.", f.User))
	//c.HTML(http.StatusOK, "login.tmpl", gin.H{
	//	"error": "user or password input error",
	//})
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusInternalServerError,
		"msg":  "user or password input error",
	})
}

func UserAdmin(c *gin.Context) {
	userList, err := core.SelectUser("SELECT id,userName,newUserDate,password FROM user")
	if err != nil {
		mlog.Error(err.Error())
		return
	}

	userAndRoleList, err := core.SelectUserAndRole("SELECT * FROM user_roles")
	if err != nil {
		mlog.Error(err.Error())
		return
	}

	c.HTML(http.StatusOK, "user.tmpl", gin.H{
		"userList":        userList,
		"userAndRoleList": userAndRoleList,
	})
}

func CreateUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "json analysis parameter error.",
		})
		return
	}
	pwd, err := core.PasswordEncryption(user.Password, global.NowKey)
	if err != nil {
		mlog.Error(err.Error())
	}
	if err := core.InsertUser(user.UserName, pwd); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadGateway,
			"msg":  err.Error(),
		})
		return
	}

	if err := core.InsertUserAndRole(user.UserName, user.RoleName); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadGateway,
			"msg":  err.Error(),
		})
		return
	}
	mlog.Info(fmt.Sprintf("add user and user add role [%s] success ---> [%s]", user.RoleName, user.UserName))
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
			mlog.Error(err.Error())
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusBadGateway,
			})
			return
		}

		if err := core.UpdateUserPwd(pwd, userName); err != nil {
			mlog.Error(err.Error())
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
		mlog.Error(err.Error())
		return
	}

	if user != "admin" {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadGateway,
		})
		mlog.Error(fmt.Sprintf("This login user does not have permission --> %s", user))
		return
	}
	if err := core.DeleteUser(userName); err != nil {
		mlog.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
	})
}

// func UpdateUserInfo(c *gin.Context) {
// 	data, _ := c.GetRawData()
// 	var body map[string]string
// 	_ = json.Unmarshal(data, &body)
// 	userId := body["userId"]
// 	userName := body["userName"]
// 	newUserDate := body["newUserDate"]
// 	user, err := c.Cookie("user")
// 	if err != nil {
// 		mlog.Error(err.Error())
// 		return
// 	}

// 	if user != "admin" {
// 		c.JSON(http.StatusOK, gin.H{
// 			"code": http.StatusBadGateway,
// 		})
// 		mlog.Error(fmt.Sprintf("This login user does not have permission --> %s", user))
// 		return
// 	}
// 	if err := core.UpdateUser(userName, newUserDate, userId); err != nil {
// 		mlog.Error(err.Error())
// 		c.JSON(http.StatusOK, gin.H{
// 			"code": http.StatusBadGateway,
// 		})
// 		return
// 	}
// 	mlog.Info(fmt.Sprintf("update user [%s] info ok.", userName))
// 	c.JSON(http.StatusOK, gin.H{
// 		"code": http.StatusOK,
// 	})
// }

func RebootHost(c *gin.Context) {
	osType := core.ShowSys()
	switch osType {
	case "linux":
		cmdCode := "reboot"
		if err := cmd.UseCmd(cmdCode); err != nil {
			mlog.Error(fmt.Sprintf("reboot host fail,%s", err))
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusBadGateway,
			})
			return
		}
		return
	case "windows":
		cmdCode := "shutdown.exe -r -f -t 0"
		if err := cmd.UseCmd(cmdCode); err != nil {
			mlog.Error(fmt.Sprintf("reboot host fail,%s", err))
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusBadGateway,
			})
			return
		}
		return
	default:
		mlog.Warn("unsupported operating system.")
		return
	}

}
