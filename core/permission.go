package core

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

var p *Permission = nil

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

func InitPermissionRoute(allRoute gin.RoutesInfo) {
	readJson()
	go func() {
		for _, v := range p.File {
			go func(v string) {
				p, err := SelectRolePermission(fmt.Sprintf("select * from roles_permission WHERE roleName = \"%v\" and permission = \"%v\"", "role_admin", v))
				if len(p) <= 0 {
					if err := InsertRolePermission("role_admin", v, "file"); err != nil {
						mlog.Error(fmt.Sprintf("add route fail -> %s,%s", v, err.Error()))
					}
				}
				if err != nil {
					log.Println(err.Error())
				}
			}(v)
		}
		for _, v := range p.User {
			go func(v string) {
				p, err := SelectRolePermission(fmt.Sprintf("select * from roles_permission WHERE roleName = \"%v\" and permission = \"%v\"", "role_admin", v))
				if len(p) <= 0 {
					if err := InsertRolePermission("role_admin", v, "user"); err != nil {
						mlog.Error(fmt.Sprintf("add route fail -> %s,%s", v, err.Error()))
					}
				}
				if err != nil {
					log.Println(err.Error())
				}
			}(v)
		}
		for _, v := range p.Url {
			go func(v string) {
				p, err := SelectRolePermission(fmt.Sprintf("select * from roles_permission WHERE roleName = \"%v\" and permission = \"%v\"", "role_admin", v))
				if len(p) <= 0 {
					if err := InsertRolePermission("role_admin", v, "url"); err != nil {
						mlog.Error(fmt.Sprintf("add route fail -> %s,%s", v, err.Error()))
					}
				}
				if err != nil {
					log.Println(err.Error())
				}
			}(v)
		}
		for _, v := range p.Log {
			go func(v string) {
				p, err := SelectRolePermission(fmt.Sprintf("select * from roles_permission WHERE roleName = \"%v\" and permission = \"%v\"", "role_admin", v))
				if len(p) <= 0 {
					if err := InsertRolePermission("role_admin", v, "log"); err != nil {
						mlog.Error(fmt.Sprintf("add route fail -> %s,%s", v, err.Error()))
					}
				}
				if err != nil {
					log.Println(err.Error())
				}
			}(v)
		}
		for _, v := range p.Cron {
			go func(v string) {
				p, err := SelectRolePermission(fmt.Sprintf("select * from roles_permission WHERE roleName = \"%v\" and permission = \"%v\"", "role_admin", v))
				if len(p) <= 0 {
					if err := InsertRolePermission("role_admin", v, "cron"); err != nil {
						mlog.Error(fmt.Sprintf("add route fail -> %s,%s", v, err.Error()))
					}
				}
				if err != nil {
					log.Println(err.Error())
				}
			}(v)
		}
		for _, v := range p.Service {
			go func(v string) {
				p, err := SelectRolePermission(fmt.Sprintf("select * from roles_permission WHERE roleName = \"%v\" and permission = \"%v\"", "role_admin", v))
				if len(p) <= 0 {
					if err := InsertRolePermission("role_admin", v, "service"); err != nil {
						mlog.Error(fmt.Sprintf("add route fail -> %s,%s", v, err.Error()))
					}
				}
				if err != nil {
					log.Println(err.Error())
				}
			}(v)
		}
		for _, v := range p.Passwd {
			go func(v string) {
				p, err := SelectRolePermission(fmt.Sprintf("select * from roles_permission WHERE roleName = \"%v\" and permission = \"%v\"", "role_admin", v))
				if len(p) <= 0 {
					if err := InsertRolePermission("role_admin", v, "passwd"); err != nil {
						mlog.Error(fmt.Sprintf("add route fail -> %s,%s", v, err.Error()))
					}
				}
				if err != nil {
					log.Println(err.Error())
				}
			}(v)
		}
		for _, v := range p.Security {
			go func(v string) {
				p, err := SelectRolePermission(fmt.Sprintf("select * from roles_permission WHERE roleName = \"%v\" and permission = \"%v\"", "role_admin", v))
				if len(p) <= 0 {
					if err := InsertRolePermission("role_admin", v, "security"); err != nil {
						mlog.Error(fmt.Sprintf("add route fail -> %s,%s", v, err.Error()))
					}
				}
				if err != nil {
					log.Println(err.Error())
				}
			}(v)
		}
		for _, v := range p.Other {
			go func(v string) {
				p, err := SelectRolePermission(fmt.Sprintf("select * from roles_permission WHERE roleName = \"%v\" and permission = \"%v\"", "role_admin", v))
				if len(p) <= 0 {
					if err := InsertRolePermission("role_admin", v, "other"); err != nil {
						mlog.Error(fmt.Sprintf("add route fail -> %s,%s", v, err.Error()))
					}
				}
				if err != nil {
					log.Println(err.Error())
				}
			}(v)
		}
	}()
}
