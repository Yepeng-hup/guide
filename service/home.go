package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"guide/core"
	"guide/global"
	"math/rand"
	"net/http"
)

type (
	updatePwd struct {
		UserName string `json:"userName"`
		Password string `json:"password"`
	}
)

var (
	messageDemo = []string{
		"无边落木萧萧下，不尽长江滚滚来。",
		"安得广厦千万间，大庇天下寒士俱欢颜。",
		"感时花溅泪，恨别鸟惊心。",
		"会当凌绝顶，一览众山小。",
		"待到秋来九月八，我花开后百花杀。",
		"他年我若为青帝，报与桃花一处开。",
		"关山难越，谁悲失路之人。",
		"萍水相逢，尽是他乡之客。",
	}
)

func showMessage() string {
	randomIndex := rand.Intn(len(messageDemo))
	randomValue := messageDemo[randomIndex]
	return randomValue
}

func showUserData() int {
	sql := "select * from user"
	v, err := core.SelectUser(sql)
	if err != nil {
		mlog.Error(err.Error())
		return 0
	}
	return len(v)
}

func showRoleData() int {
	sql := "select * from roles"
	v, err := core.SelectRole(sql)
	if err != nil {
		mlog.Error(err.Error())
		return 0
	}
	return len(v)
}


func LoginDataSource(c *gin.Context){
	valueList := []string{}
	countsMap := make(map[string]int)
	sql := "select loginDate, userName from login_count where userName != 'admin'"
	r, err := core.SelectLoginUser(sql)
	if err != nil {
		mlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"cate": []string{},
			"values": []string{},
		})
		return
	}

	for _, v := range r {
		valueList = append(valueList, v.UserName)
	}

    for _, item := range valueList {
        countsMap[item]++
    }
	
	keys := make([]string, 0, len(countsMap))
    values := make([]int, 0, len(countsMap))
	for k, v := range countsMap {
        keys = append(keys, k)
        values = append(values, v)
    }

	c.JSON(http.StatusOK, gin.H{
		"cate": keys,
    	"values": values,
	})
}

func HomeIndex(c *gin.Context) {
	randomValue := showMessage()
	userNum := showUserData()
	roleNum := showRoleData()
	c.HTML(http.StatusOK, "home.tmpl", gin.H{
		"messageValue": randomValue,
		"userNum":      userNum,
		"roleNum":      roleNum,
	})
}

func UpdateHomePwd(c *gin.Context) {
	var u updatePwd
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "json analysis parameter error.",
		})
		return
	}

	pwd, err := core.PasswordEncryption(u.Password, global.NowKey)
	if err != nil {
		mlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadGateway,
			"msg":  err.Error(),
		})
		return
	}

	if err := core.UpdateUserPwd(pwd, u.UserName); err != nil {
		mlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadGateway,
			"msg":  err.Error(),
		})
		return
	}

	mlog.Info(fmt.Sprintf("[%s] update pwd OK.", u.UserName))
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "update user pwd OK.",
	})

}
