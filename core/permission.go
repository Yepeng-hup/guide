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

// debug code

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
	//var lock sync.Mutex
	//const ginStr = "/sta/*filepath"
	//routePathList := make([]string, 0)
	//for _, route := range allRoute {
	//	//fmt.Println(route.Path)
	//	if route.Path == ginStr {
	//		continue
	//	}
	//	routePathList = append(routePathList, route.Path)
	//}
	//fmt.Println(routePathList)
	//go func() {
	//	lock.Lock()
	//	for _, route := range allRoute {
	//		fmt.Println(route.Path)
	//		// 覆盖写入,表的列必须加上唯一索引或key
	//		//db.Exec("REPLACE INTO permission_all (permission_all_path)VALUES(?)", route.Path)
	//		// 其他方法
	//		// sql := INSERT INTO permission_all (permission_all_path) VALUES ('/login') ON DUPLICATE KEY UPDATE permission_all_path = VALUES(permission_all_path);
	//	}
	//	lock.Unlock()
	//	log.Println("init start permission table")
	//}()
}
