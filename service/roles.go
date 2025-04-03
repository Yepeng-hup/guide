package service

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"guide/core"
	"log"
	"net/http"
	"os"
)

type (
	permission struct {
		User     []string `json:"user"`
		Url      []string `json:"url"`
		File     []string `json:"file"`
		Service  []string `json:"service"`
		Passwd   []string `json:"passwd"`
		Log      []string `json:"log"`
		Security []string `json:"security"`
		Cron     []string `json:"cron"`
		Other    []string `json:"other"`
	}
)

func RoleAdmin(c *gin.Context) {
	roleList, err := core.SelectRole("SELECT * FROM roles")
	if err != nil {
		mlog.Error(err.Error())
		return
	}

	c.HTML(http.StatusOK, "roles.tmpl", gin.H{
		"roleList": roleList,
	})

}

func SelectRole(c *gin.Context) {
	roleList, err := core.SelectRole("SELECT * FROM roles")
	if err != nil {
		mlog.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":     http.StatusOK,
		"roleList": roleList,
	})

}

func SelectRolePermission(c *gin.Context) {
	role := c.PostForm("roleName")
	p, err := core.SelectRolePermission(fmt.Sprintf("select * from roles_permission WHERE roleName = \"%v\"", role))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":       http.StatusOK,
		"permission": p,
	})
}

var p *permission = nil

func readJson() {
	jsonFilePath := "tools/permission.json"
	file, err := os.Open(jsonFilePath)
	if err != nil {
		log.Fatal("open json file: ", err.Error())
	}
	defer file.Close()
	f := bufio.NewReader(file)
	configObj := json.NewDecoder(f)
	if err = configObj.Decode(&p); err != nil {
		log.Fatal(err.Error())
		return
	}
	return
}

type RequestData struct {
	RoleName   string   `json:"roleName"`
	Permission []string `json:"permission"`
}

func useCreateRole(role, permissionName string) error {
	s := core.Security{}
	if s.CheckForm(role) {
		mlog.Error("error Illegal characters.")
		return fmt.Errorf("%s", "error Illegal characters.")
	}
	readJson()
	roleList, _ := core.SelectRole(fmt.Sprintf("select * from roles WHERE roleName = \"%v\"", role))
	if len(roleList) == 0 {
		mlog.Info(fmt.Sprintf("create role -> %v", role))
		if err := core.InsertRole(role); err != nil {
			return err
		}
	}

	switch permissionName {
	case "user":
		for _, v := range p.User {
			p, err := core.SelectRolePermission(fmt.Sprintf("select * from roles_permission WHERE roleName = \"%v\" and permission = \"%v\"", role, v))
			if len(p) <= 0 {
				if err := core.InsertRolePermission(role, v, "user"); err != nil {
					mlog.Error(fmt.Sprintf("add route fail -> %s,%s", v, err.Error()))
				}
			} else {
				mlog.Warn(fmt.Sprintf("permission route [%v] already exists.", v))
			}
			if err != nil {
				mlog.Error(err.Error())
			}
		}
		mlog.Info(fmt.Sprintf("%s add [user] permission.", role))
	case "url":
		for _, v := range p.Url {
			p, err := core.SelectRolePermission(fmt.Sprintf("select * from roles_permission WHERE roleName = \"%v\" and permission = \"%v\"", role, v))
			if len(p) <= 0 {
				if err := core.InsertRolePermission(role, v, "url"); err != nil {
					mlog.Error(fmt.Sprintf("add route fail -> %s,%s", v, err.Error()))
				}
			} else {
				mlog.Warn(fmt.Sprintf("permission route [%v] already exists.", v))
			}
			if err != nil {
				mlog.Error(err.Error())
			}
		}
		mlog.Info(fmt.Sprintf("%s add [url] permission.", role))
	case "file":
		for _, v := range p.File {
			p, err := core.SelectRolePermission(fmt.Sprintf("select * from roles_permission WHERE roleName = \"%v\" and permission = \"%v\"", role, v))
			if len(p) <= 0 {
				if err := core.InsertRolePermission(role, v, "file"); err != nil {
					mlog.Error(fmt.Sprintf("add route fail -> %s,%s", v, err.Error()))
				}
			} else {
				mlog.Warn(fmt.Sprintf("permission route [%v] already exists.", v))
			}
			if err != nil {
				mlog.Error(err.Error())
			}
		}
		mlog.Info(fmt.Sprintf("%s add [file] permission.", role))
	case "service":
		for _, v := range p.Service {
			p, err := core.SelectRolePermission(fmt.Sprintf("select * from roles_permission WHERE roleName = \"%v\" and permission = \"%v\"", role, v))
			if len(p) <= 0 {
				if err := core.InsertRolePermission(role, v, "service"); err != nil {
					mlog.Error(fmt.Sprintf("add route fail -> %s,%s", v, err.Error()))
				}
			} else {
				mlog.Warn(fmt.Sprintf("permission route [%v] already exists.", v))
			}
			if err != nil {
				mlog.Error(err.Error())
			}
		}
		mlog.Info(fmt.Sprintf("%s add [service] permission.", role))
	case "passwd":
		for _, v := range p.Passwd {
			p, err := core.SelectRolePermission(fmt.Sprintf("select * from roles_permission WHERE roleName = \"%v\" and permission = \"%v\"", role, v))
			if len(p) <= 0 {
				if err := core.InsertRolePermission(role, v, "passwd"); err != nil {
					mlog.Error(fmt.Sprintf("add route fail -> %s,%s", v, err.Error()))
				}
			} else {
				mlog.Warn(fmt.Sprintf("permission route [%v] already exists.", v))
			}
			if err != nil {
				mlog.Error(err.Error())
			}
		}
		mlog.Info(fmt.Sprintf("%s add [passwd] permission.", role))
	case "log":
		for _, v := range p.Log {
			p, err := core.SelectRolePermission(fmt.Sprintf("select * from roles_permission WHERE roleName = \"%v\" and permission = \"%v\"", role, v))
			if len(p) <= 0 {
				if err := core.InsertRolePermission(role, v, "log"); err != nil {
					mlog.Error(fmt.Sprintf("add route fail -> %s,%s", v, err.Error()))
				}
			} else {
				mlog.Warn(fmt.Sprintf("permission route [%v] already exists.", v))
			}
			if err != nil {
				mlog.Error(err.Error())
			}
		}
		mlog.Info(fmt.Sprintf("%s add [log] permission.", role))
	case "security":
		for _, v := range p.Security {
			p, err := core.SelectRolePermission(fmt.Sprintf("select * from roles_permission WHERE roleName = \"%v\" and permission = \"%v\"", role, v))
			if len(p) <= 0 {
				if err := core.InsertRolePermission(role, v, "security"); err != nil {
					mlog.Error(fmt.Sprintf("add route fail -> %s,%s", v, err.Error()))
				}
			} else {
				mlog.Warn(fmt.Sprintf("permission route [%v] already exists.", v))
			}
			if err != nil {
				mlog.Error(err.Error())
			}
		}
		mlog.Info(fmt.Sprintf("%s add [security] permission.", role))
	case "cron":
		for _, v := range p.Cron {
			p, err := core.SelectRolePermission(fmt.Sprintf("select * from roles_permission WHERE roleName = \"%v\" and permission = \"%v\"", role, v))
			if len(p) <= 0 {
				if err := core.InsertRolePermission(role, v, "cron"); err != nil {
					mlog.Error(fmt.Sprintf("add route fail -> %s,%s", v, err.Error()))
				}
			} else {
				mlog.Warn(fmt.Sprintf("permission route [%v] already exists.", v))
			}

			if err != nil {
				mlog.Error(err.Error())
			}
		}
		mlog.Info(fmt.Sprintf("%s add [cron] permission.", role))
	case "other":
		for _, v := range p.Other {
			p, err := core.SelectRolePermission(fmt.Sprintf("select * from roles_permission WHERE roleName = \"%v\" and permission = \"%v\"", role, v))
			if len(p) <= 0 {
				if err := core.InsertRolePermission(role, v, "other"); err != nil {
					mlog.Error(fmt.Sprintf("add route fail -> %s,%s", v, err.Error()))
				}
			} else {
				mlog.Warn(fmt.Sprintf("permission route [%v] already exists.", v))
			}
			if err != nil {
				mlog.Error(err.Error())
			}
		}
		mlog.Info(fmt.Sprintf("%s add [other] permission.", role))
	default:
		mlog.Error(fmt.Sprintf("%s add permission [%s] fail, error not is permission.", role, permissionName))
		return nil
	}
	return nil
}

func CreateRole(c *gin.Context) {
	var req RequestData

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "json analysis parameter error.",
		})
		return
	}

	for _, v := range req.Permission {
		if err := useCreateRole(req.RoleName, v); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusInternalServerError,
				"msg":  err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
	})
}

type Role struct {
	RoleName string `json:"roleName"`
}

type deleteRoute struct {
	RoleName        string `json:"roleName"`
	PermissionRoute string `json:"permissionRoute"`
}

func DelRole(c *gin.Context) {
	var r Role
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "json analysis role error.",
		})
		return
	}

	if err := core.DeleteRole(r.RoleName); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
	})
}

func DelRolePermissionRoute(c *gin.Context) {
	var d deleteRoute
	if err := c.ShouldBindJSON(&d); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "json analysis role error.",
		})
		return
	}

	if err := core.DeleteRolePermissionRoute(d.RoleName, d.PermissionRoute); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
	})
}
