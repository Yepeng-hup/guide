package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type (
	LoginFrom struct {
		User string
		Password string
	}
)


func Login(c *gin.Context)  {
	c.HTML(http.StatusOK, "login.tmpl", gin.H{})
}


func LoginCk(c *gin.Context)  {
	f := LoginFrom{
		User: c.PostForm("Name"),
		Password: c.PostForm("Password"),
	}
	if f.User == "admin" && f.Password == "guide" {
		cookie := http.Cookie{Name: "user", Value: f.User, MaxAge: 10800}
		http.SetCookie(c.Writer, &cookie)
		c.Redirect(http.StatusMovedPermanently, "/user/index")
		return
	}else {
		c.HTML(http.StatusOK, "login.tmpl", gin.H{
			"error": "user or password input error",
		})
	}
}



func UserAdmin(c *gin.Context)  {
	c.HTML(http.StatusOK, "user.tmpl", gin.H{})
}
