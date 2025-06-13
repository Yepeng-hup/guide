package core

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"guide/global"
	"log"
	"os"
	"time"
)

var (
	tableList = []string{"cron", "service_tools", "user_passwd", "user", "error_log", "cpu", "mem", "url", "url_type", "roles", "user_roles", "roles_permission", "permission", "blacklist", "login_count"}
)

func ConnDb() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "guide.db")
	if err != nil {
		return nil, fmt.Errorf("conn db sqlite fail, %v", err.Error())
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(10 * time.Minute)
	return db, nil
}

func checkTableIfCreate() []string {
	t := make([]string, 0)
	database, err := ConnDb()
	if err != nil {
		log.Fatal(err.Error())
	}

	rows, err := database.Query("SELECT name FROM sqlite_master WHERE type='table';")
	if err != nil {
		log.Fatal(err.Error())
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			log.Fatal(err.Error())
			return nil
		}
		t = append(t, tableName)
	}
	return t
}

func CreateGuideAllTable() error {
	db, err := ConnDb()
	if err != nil {
		log.Fatal(err.Error())
	}

	list := checkTableIfCreate()
	for _, v := range tableList {
		if IfElement(list, v) {
		} else {
			switch v {
			case "cron":
				sql := `CREATE TABLE IF NOT EXISTS cron (id INTEGER PRIMARY KEY, cronNewDate TEXT DEFAULT (strftime('%Y-%m-%d %H:%M', 'now', 'localtime')), cronName TEXT, cronTime TEXT, cronCode TEXT, cronNotes TEXT);`
				_, err = db.Exec(sql)
				if err != nil {
					return fmt.Errorf("create table cron fail,%s", err.Error())
				}
			case "service_tools":
				sql := `CREATE TABLE IF NOT EXISTS service_tools (id INTEGER PRIMARY KEY, serviceName TEXT, serviceDate TEXT DEFAULT (strftime('%Y-%m-%d %H:%M', 'now', 'localtime')), startCmd TEXT, serviceNotes TEXT);`
				_, err = db.Exec(sql)
				if err != nil {
					return fmt.Errorf("create table service_tools fail,%s", err.Error())
				}
			case "user_passwd":
				sql := `CREATE TABLE IF NOT EXISTS user_passwd (id INTEGER PRIMARY KEY, serviceName TEXT, user TEXT, password TEXT, Notes TEXT);`
				_, err = db.Exec(sql)
				if err != nil {
					return fmt.Errorf("create table user_passwd fail,%s", err.Error())
				}
			case "user":
				sql := `CREATE TABLE IF NOT EXISTS user (id INTEGER PRIMARY KEY, userName TEXT, newUserDate TEXT DEFAULT (strftime('%Y-%m-%d %H:%M', 'now', 'localtime')), password TEXT, UNIQUE (userName));`
				_, err = db.Exec(sql)
				if err != nil {
					return fmt.Errorf("create table user fail,%s", err.Error())
				}
			case "error_log":
				sql := `CREATE TABLE IF NOT EXISTS error_log (id INTEGER PRIMARY KEY, newLogDate TEXT DEFAULT (strftime('%Y-%m-%d %H:%M', 'now', 'localtime')), types TEXT, logtext TEXT);`
				_, err = db.Exec(sql)
				if err != nil {
					return fmt.Errorf("create table error_log fail,%s", err.Error())
				}
			case "cpu":
				sql := `CREATE TABLE IF NOT EXISTS cpu (id INTEGER PRIMARY KEY, cpunum INT);`
				_, err = db.Exec(sql)
				if err != nil {
					return fmt.Errorf("create table cpu fail,%s", err.Error())
				}
			case "mem":
				sql := `CREATE TABLE IF NOT EXISTS mem (id INTEGER PRIMARY KEY, memnum FLOAT);`
				_, err = db.Exec(sql)
				if err != nil {
					return fmt.Errorf("create table mem fail,%s", err.Error())
				}
			case "url":
				sql := `CREATE TABLE IF NOT EXISTS url (id INTEGER PRIMARY KEY, urlName TEXT, urlAddress TEXT, urlType TEXT, urlNotes TEXT);`
				_, err = db.Exec(sql)
				if err != nil {
					return fmt.Errorf("create table url fail,%s", err.Error())
				}
			case "url_type":
				sql := `CREATE TABLE IF NOT EXISTS url_type (id INTEGER PRIMARY KEY, urlType TEXT);`
				_, err = db.Exec(sql)
				if err != nil {
					return fmt.Errorf("create table url_type fail,%s", err.Error())
				}
			case "roles":
				sql := `CREATE TABLE IF NOT EXISTS roles (id INTEGER PRIMARY KEY, roleName TEXT, newRoleDate TEXT DEFAULT (strftime('%Y-%m-%d %H:%M', 'now', 'localtime')), UNIQUE (roleName));`
				_, err = db.Exec(sql)
				if err != nil {
					return fmt.Errorf("create table roles fail,%s", err.Error())
				}
			case "user_roles":
				sql := `CREATE TABLE IF NOT EXISTS user_roles (id INTEGER PRIMARY KEY, userName TEXT, roleName TEXT, UNIQUE (userName));`
				_, err = db.Exec(sql)
				if err != nil {
					return fmt.Errorf("create table user_roles fail,%s", err.Error())
				}
			case "roles_permission":
				sql := `CREATE TABLE IF NOT EXISTS roles_permission (id INTEGER PRIMARY KEY, roleName TEXT, permission TEXT, label TEXT);`
				_, err = db.Exec(sql)
				if err != nil {
					return fmt.Errorf("create table roles_permission fail,%s", err.Error())
				}
			case "permission":
				sql := `CREATE TABLE IF NOT EXISTS permission (id INTEGER PRIMARY KEY, permission_url TEXT);`
				_, err = db.Exec(sql)
				if err != nil {
					return fmt.Errorf("create table permission fail,%s", err.Error())
				}
			case "blacklist":
				sql := `CREATE TABLE IF NOT EXISTS blacklist (id INTEGER PRIMARY KEY, newAddDate TEXT DEFAULT (strftime('%Y-%m-%d %H:%M', 'now', 'localtime')), ip TEXT);`
				_, err = db.Exec(sql)
				if err != nil {
					return fmt.Errorf("create table blacklist fail,%s", err.Error())
				}
			case "login_count":
				sql := `CREATE TABLE IF NOT EXISTS login_count (id INTEGER PRIMARY KEY, loginDate TEXT DEFAULT (strftime('%Y-%m-%d %H:%M', 'now', 'localtime')), userName TEXT);`
				_, err = db.Exec(sql)
				if err != nil {
					return fmt.Errorf("create table login_count fail,%s", err.Error())
				}

			default:
				return nil
			}
		}
	}
	return nil
}

func InitUser() {
	db, err := ConnDb()
	if err != nil {
		log.Println(err)
		return
	}

	sql := "SELECT userName FROM user WHERE userName = \"admin\""
	rows, err := db.Query(sql)

	if err != nil {
		mlog.Error("init fail.")
		mlog.Error(fmt.Sprintf("query user table fail,%s", err.Error()))
		os.Exit(1)
	}
	defer rows.Close()

	var user struct {
		User string
	}

	userList := make([]User, 0)

	for rows.Next() {
		err := rows.Scan(&user.User)
		if err != nil {
			mlog.Error("init fail.")
			mlog.Error(err.Error())
			os.Exit(1)
		}
		userList = append(userList, user)
	}

	if len(userList) < 1 {
		insertSQL := `INSERT INTO user (userName, password) VALUES (?,?);`

		encryptionPwd, err := PasswordEncryption("guide654321", global.NowKey)
		if err != nil {
			mlog.Error(fmt.Sprintf("init user encryption passwd fail,%s", err.Error()))
			os.Exit(1)
		}

		_, err = db.Exec(insertSQL, "admin", encryptionPwd)

		if err != nil {
			mlog.Error(fmt.Sprintf("init db user fail,%s", err.Error()))
			os.Exit(1)
		}
	} else {
		return
	}
}

var p1 *Permission

func readJson1() {
	jsonFilePath := "tools/permission.json"
	file, err := os.Open(jsonFilePath)
	if err != nil {
		log.Fatal("open json file: ", err.Error())
	}
	defer file.Close()
	f := bufio.NewReader(file)
	configObj := json.NewDecoder(f)
	if err = configObj.Decode(&p1); err != nil {
		log.Fatal(err.Error())
		return
	}
	return
}

func InitAdminRoleAndAdminRolePermission() {
	db, err := ConnDb()
	if err != nil {
		log.Println(err)
		return
	}
	sql := "SELECT userName FROM user_roles WHERE userName = \"admin\""
	rows, err := db.Query(sql)

	if err != nil {
		mlog.Error("adminRole init fail.")
		mlog.Error(fmt.Sprintf("query user_roles table fail,%s", err.Error()))
		os.Exit(1)
	}
	defer rows.Close()

	var user struct {
		User string
	}

	userList := make([]User, 0)

	for rows.Next() {
		err := rows.Scan(&user.User)
		if err != nil {
			mlog.Error("adminRole init fail.")
			mlog.Error(err.Error())
			os.Exit(1)
		}
		userList = append(userList, user)
	}

	if len(userList) < 1 {
		insertSQL := `INSERT INTO user_roles (userName, roleName) VALUES (?,?);`

		_, err = db.Exec(insertSQL, "admin", "role_admin")

		if err != nil {
			mlog.Error(fmt.Sprintf("init db role [role_admin] fail,%s", err.Error()))
			os.Exit(1)
		} else {
			// 写入role_admin权限
			//readJson1()
			//for _, v := range p1.User {
			//
			//}
			return

		}
	} else {
		return
	}

}

func InsertAct(params ...string) {
	db, err := ConnDb()
	if err != nil {
		log.Println(err)
		return
	}
	insertSQL := `INSERT INTO cron (cronName, cronTime, cronCode, cronNotes) VALUES (?, ?, ?, ?);`
	_, err = db.Exec(insertSQL, params[0], params[1], params[2], params[3])
	if err != nil {
		mlog.Error(fmt.Sprintf("insert data to cron fail,%s", err.Error()))
		return
	}
}

func InsertActSTools(p ...string) error {
	db, err := ConnDb()
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	insertSQL := `INSERT INTO service_tools (serviceName, startCmd, serviceNotes) VALUES (?, ?, ?);`
	_, err = db.Exec(insertSQL, p[0], p[1], p[2])
	if err != nil {
		return fmt.Errorf("insert data to service_tools fail,%s", err.Error())
	}
	return nil
}

func InsertRole(role string) error {
	db, err := ConnDb()
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	insertSQL := `INSERT INTO roles (roleName) VALUES (?);`
	_, err = db.Exec(insertSQL, role)
	if err != nil {
		return fmt.Errorf("insert data to roles fail,%s", err.Error())
	}
	return nil
}

func InsertRolePermission(permission ...string) error {
	db, err := ConnDb()
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	insertSQL := `INSERT INTO roles_permission (roleName, permission, label) VALUES (?,?,?);`
	_, err = db.Exec(insertSQL, permission[0], permission[1], permission[2])
	if err != nil {
		return fmt.Errorf("insert data to role_permission fail,%s", err.Error())
	}
	return nil
}

func InsertUserPwd(p ...string) error {
	db, err := ConnDb()
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	insertSQL := `INSERT INTO user_passwd (serviceName,user ,password ,Notes) VALUES (?, ?, ?, ?);`
	_, err = db.Exec(insertSQL, p[0], p[1], p[2], p[3])
	if err != nil {
		return fmt.Errorf("insert data to user_passwd fail,%s", err.Error())
	}
	return nil
}

func InsertUser(p ...string) error {
	db, err := ConnDb()
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	insertSQL := `INSERT INTO user (userName,password) VALUES (?, ?);`
	_, err = db.Exec(insertSQL, p[0], p[1])
	if err != nil {
		return fmt.Errorf("insert data to user fail,%s", err.Error())
	}
	return nil
}

func InsertActErrorLog(p ...string) error {
	db, err := ConnDb()
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	insertSQL := `INSERT INTO error_log (types, logtext) VALUES (?, ?);`
	_, err = db.Exec(insertSQL, p[0], p[1])
	if err != nil {
		return fmt.Errorf("insert data error_log fail,%s", err.Error())
	}
	return nil
}

func InsertActCPU(p ...int) error {
	db, err := ConnDb()
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	insertSQL := `INSERT INTO cpu (cpunum) VALUES (?);`
	_, err = db.Exec(insertSQL, p[0])
	if err != nil {
		return fmt.Errorf("insert data cpu fail,%s", err.Error())
	}
	return nil
}

func InsertActMEM(p ...float64) error {
	db, err := ConnDb()
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	insertSQL := `INSERT INTO mem (memnum) VALUES (?);`
	_, err = db.Exec(insertSQL, p[0])
	if err != nil {
		return fmt.Errorf("insert data mem fail,%s", err.Error())
	}
	return nil
}

func InsertActUrl(p ...string) error {
	db, err := ConnDb()
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	insertSQL := `INSERT INTO url (urlName, urlAddress, urlType, urlNotes) VALUES (?, ?, ?, ?);`
	_, err = db.Exec(insertSQL, p[0], p[1], p[2], p[3])
	if err != nil {
		return fmt.Errorf("insert data url fail,%s", err.Error())
	}
	return nil
}

func InsertActUrlType(p ...string) error {
	db, err := ConnDb()
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	insertSQL := `INSERT INTO url_type (urlType) VALUES (?);`
	_, err = db.Exec(insertSQL, p[0])
	if err != nil {
		return fmt.Errorf("insert data url_type fail,%s", err.Error())
	}
	return nil
}

func InsertUserAndRole(p ...string) error {
	db, err := ConnDb()
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	insertSQL := `INSERT INTO user_roles (userName,roleName) VALUES (?,?);`
	_, err = db.Exec(insertSQL, p[0], p[1])
	if err != nil {
		return fmt.Errorf("insert data user_roles fail,%s", err.Error())
	}
	return nil
}

func InsertActBlacklist(p ...string) error {
	db, err := ConnDb()
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	insertSQL := `INSERT INTO blacklist (ip) VALUES (?);`
	_, err = db.Exec(insertSQL, p[0])
	if err != nil {
		return fmt.Errorf("insert data to blacklist fail,%s", err.Error())
	}
	return nil
}

func InsertActLoginUser(p ...string) error {
	db, err := ConnDb()
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	insertSQL := `INSERT INTO login_count (userName) VALUES (?);`
	_, err = db.Exec(insertSQL, p[0])
	if err != nil {
		return fmt.Errorf("insert data to login_count fail,%s", err.Error())
	}
	return nil
}

func SelectAct(k, v string, b bool) ([]Cron, error) {
	db, err := ConnDb()
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	if b {
		rows, err := db.Query("SELECT * FROM cron")
		if err != nil {
			return nil, fmt.Errorf("query cron table fail,%s", err.Error())
		}
		defer rows.Close()
		var cron struct {
			CronId      string
			CronNewDate string
			CronName    string
			CronTime    string
			CronCode    string
			CronNotes   string
		}
		cronList := make([]Cron, 0)
		for rows.Next() {
			err := rows.Scan(&cron.CronId, &cron.CronNewDate, &cron.CronName, &cron.CronTime, &cron.CronCode, &cron.CronNotes)
			if err != nil {
				return nil, fmt.Errorf(err.Error())
			}
			cronList = append(cronList, cron)
		}
		return cronList, nil

	} else {
		ssql := "SELECT * FROM cron  WHERE ? = ?"
		rows, err := db.Query(ssql, k, v)
		if err != nil {
			return nil, fmt.Errorf("query cron table fail,%s", err.Error())
		}
		defer rows.Close()
		var cron struct {
			CronId      string
			CronNewDate string
			CronName    string
			CronTime    string
			CronCode    string
			CronNotes   string
		}
		cronList := make([]Cron, 0)
		for rows.Next() {
			err := rows.Scan(&cron.CronId, &cron.CronNewDate, &cron.CronName, &cron.CronTime, &cron.CronCode, &cron.CronNotes)
			if err != nil {
				return nil, fmt.Errorf(err.Error())
			}
			cronList = append(cronList, cron)
		}
		return cronList, nil
	}
}

func SelectActSTools(selectSql string) ([]ServiceTools, error) {
	db, err := ConnDb()
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	rows, err := db.Query(selectSql)
	if err != nil {
		return nil, fmt.Errorf("query service_tools table fail,%s", err.Error())
	}
	defer rows.Close()
	var serviceTools struct {
		Id           string
		ServiceName  string
		ServiceDate  string
		StartCmd     string
		ServiceNotes string
	}
	serviceToolsList := make([]ServiceTools, 0)
	for rows.Next() {
		err := rows.Scan(&serviceTools.Id, &serviceTools.ServiceName, &serviceTools.ServiceDate, &serviceTools.StartCmd, &serviceTools.ServiceNotes)
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}
		serviceToolsList = append(serviceToolsList, serviceTools)
	}
	return serviceToolsList, nil
}

func SelectUserPwd() ([]UserPwd, error) {
	db, err := ConnDb()
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	sql := "SELECT * FROM user_passwd"
	rows, err := db.Query(sql)
	if err != nil {
		return nil, fmt.Errorf("query user_passwd table fail,%s", err.Error())
	}
	defer rows.Close()
	var userPwd struct {
		Id          string
		ServiceName string
		User        string
		Passwd      string
		Notes       string
	}
	userPwdList := make([]UserPwd, 0)
	for rows.Next() {
		err := rows.Scan(&userPwd.Id, &userPwd.ServiceName, &userPwd.User, &userPwd.Passwd, &userPwd.Notes)
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}
		userPwdList = append(userPwdList, userPwd)
	}
	return userPwdList, nil
}

func SelectUser(selectSql string) ([]UserAll, error) {
	db, err := ConnDb()
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	rows, err := db.Query(selectSql)
	if err != nil {
		return nil, fmt.Errorf("query user table fail,%s", err.Error())
	}
	defer rows.Close()
	var user struct {
		Id          string
		UserName    string
		NewUserDate string
		Password    string
	}
	userList := make([]UserAll, 0)
	for rows.Next() {
		err := rows.Scan(&user.Id, &user.UserName, &user.NewUserDate, &user.Password)
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}
		userList = append(userList, user)
	}
	return userList, nil
}

func SelectUserAndRole(selectSql string) ([]UserRole, error) {
	db, err := ConnDb()
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	rows, err := db.Query(selectSql)
	if err != nil {
		return nil, fmt.Errorf("query user_roles table fail,%s", err.Error())
	}
	defer rows.Close()
	var userRole struct {
		Id       string
		UserName string
		RoleName string
	}
	userAndRoleList := make([]UserRole, 0)
	for rows.Next() {
		err := rows.Scan(&userRole.Id, &userRole.UserName, &userRole.RoleName)
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}
		userAndRoleList = append(userAndRoleList, userRole)
	}
	return userAndRoleList, nil
}

func SelectRole(selectSql string) ([]Roles, error) {
	db, err := ConnDb()
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	rows, err := db.Query(selectSql)
	if err != nil {
		return nil, fmt.Errorf("query roles table fail,%s", err.Error())
	}
	defer rows.Close()
	var roles struct {
		Id          string
		RoleName    string
		NewRoleDate string
	}
	roleList := make([]Roles, 0)
	for rows.Next() {
		err := rows.Scan(&roles.Id, &roles.RoleName, &roles.NewRoleDate)
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}
		roleList = append(roleList, roles)
	}
	return roleList, nil
}

func SelectRolePermission(selectSql string) ([]RolePermission, error) {
	db, err := ConnDb()
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	rows, err := db.Query(selectSql)
	if err != nil {
		return nil, fmt.Errorf("query roles_permission table fail,%s", err.Error())
	}
	defer rows.Close()
	var rolePermission struct {
		Id         string
		RoleName   string
		Permission string
		Label      string
	}
	rolePermissionList := make([]RolePermission, 0)
	for rows.Next() {
		err := rows.Scan(&rolePermission.Id, &rolePermission.RoleName, &rolePermission.Permission, &rolePermission.Label)
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}
		rolePermissionList = append(rolePermissionList, rolePermission)
	}
	return rolePermissionList, nil
}

func SelectLog(selectSql string) ([]ErrorLog, error) {
	db, err := ConnDb()
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	rows, err := db.Query(selectSql)
	if err != nil {
		return nil, fmt.Errorf("query error_log table fail,%s", err.Error())
	}
	defer rows.Close()
	var logs struct {
		Id         string
		Date       string
		LogType    string
		LogContent string
	}
	logList := make([]ErrorLog, 0)
	for rows.Next() {
		err := rows.Scan(&logs.Id, &logs.Date, &logs.LogType, &logs.LogContent)
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}
		logList = append(logList, logs)
	}
	return logList, nil
}

func SelectCPU(selectSql string) ([]Cpu, error) {
	db, err := ConnDb()
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	rows, err := db.Query(selectSql)
	if err != nil {
		return nil, fmt.Errorf("query cpu table fail,%s", err.Error())
	}
	defer rows.Close()
	var cpu struct {
		CpuNum int
	}
	cpuList := make([]Cpu, 0)
	for rows.Next() {
		err := rows.Scan(&cpu.CpuNum)
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}
		cpuList = append(cpuList, cpu)
	}
	return cpuList, nil
}

func SelectMEM(selectSql string) ([]Mem, error) {
	db, err := ConnDb()
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	rows, err := db.Query(selectSql)
	if err != nil {
		return nil, fmt.Errorf("query mem table fail,%s", err.Error())
	}
	defer rows.Close()
	var mem struct {
		MemNum float64
	}
	memList := make([]Mem, 0)
	for rows.Next() {
		err := rows.Scan(&mem.MemNum)
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}
		memList = append(memList, mem)
	}
	return memList, nil
}

func SelectUrl(selectSql string) ([]Url, error) {
	db, err := ConnDb()
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	rows, err := db.Query(selectSql)
	if err != nil {
		return nil, fmt.Errorf("query url table fail,%s", err.Error())
	}
	defer rows.Close()
	var url struct {
		UrlName    string
		UrlAddress string
		UrlType    string
		UrlNotes   string
	}
	urlList := make([]Url, 0)
	for rows.Next() {
		err := rows.Scan(&url.UrlName, &url.UrlAddress, &url.UrlType, &url.UrlNotes)
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}
		urlList = append(urlList, url)
	}
	return urlList, nil
}

func SelectUrlType(selectSql string) ([]UrlType, error) {
	db, err := ConnDb()
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	rows, err := db.Query(selectSql)
	if err != nil {
		return nil, fmt.Errorf("query typeName table fail,%s", err.Error())
	}
	defer rows.Close()
	var urlType struct {
		// Id         string
		TypeName string
	}
	urlTypeList := make([]UrlType, 0)
	for rows.Next() {
		err := rows.Scan(&urlType.TypeName)
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}
		urlTypeList = append(urlTypeList, urlType)
	}
	return urlTypeList, nil
}

func SelectUserPermission(user ...string) ([]string, error) {
	labelList := make([]string, 0)
	userAndRoleList, err := SelectUserAndRole(fmt.Sprintf("SELECT * FROM user_roles WHERE userName = \"%s\"", user[0]))
	if err != nil {
		mlog.Error("func [SelectUserAndRole] use error," + err.Error())
		return nil, err
	}

	if len(userAndRoleList) <= 0 {
		return nil, fmt.Errorf("user not binding role -> [%s]", user[0])
	}

	p, err := SelectRolePermission(fmt.Sprintf("select DISTINCT * from roles_permission WHERE roleName = \"%v\"", userAndRoleList[0].RoleName))
	if err != nil {
		mlog.Error("func [SelectRolePermission] use error," + err.Error())
		return nil, err
	}

	for _, v := range p {
		labelList = append(labelList, v.Label)
	}
	uniqueLabel := DeduplicateGeneric(labelList)
	return uniqueLabel, nil
}

func SelectBlacklistIp(selectSql string) ([]BIp, error) {
	db, err := ConnDb()
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	rows, err := db.Query(selectSql)
	if err != nil {
		return nil, fmt.Errorf("query blacklist table fail,%s", err.Error())
	}
	defer rows.Close()
	var blacklists struct {
		Id          string
		NewAddDate string
		Ip string
	}
	b := make([]BIp, 0)
	for rows.Next() {
		err := rows.Scan(&blacklists.Id,  &blacklists.NewAddDate, &blacklists.Ip)
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}
		b = append(b, blacklists)
	}
	return b, nil
}

func SelectLoginUser(selectSql string) ([]LoginCount, error) {
	db, err := ConnDb()
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	rows, err := db.Query(selectSql)
	if err != nil {
		return nil, fmt.Errorf("query login_count table fail,%s", err.Error())
	}
	defer rows.Close()
	var loginCount struct {
		// Id         string
		LoginDate string
		UserName string
	}
	loginCountList := make([]LoginCount, 0)
	for rows.Next() {
		err := rows.Scan(&loginCount.LoginDate,&loginCount.UserName)
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}
		loginCountList = append(loginCountList, loginCount)
	}
	return loginCountList, nil
}

func DeleteAct(p ...string) error {
	db, err := ConnDb()
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	deleteSQL := "DELETE FROM cron WHERE cronName = ?"
	stmt, err := db.Prepare(deleteSQL)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	defer stmt.Close()
	_, err = stmt.Exec(p[0])
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	mlog.Info(fmt.Sprintf("delete cron database record ok. name -> [%s].", p[0]))
	return nil
}

func DeleteActSTools(p ...string) error {
	db, err := ConnDb()
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	deleteSQL := "DELETE FROM service_tools WHERE serviceName = ?"
	stmt, err := db.Prepare(deleteSQL)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	defer stmt.Close()
	_, err = stmt.Exec(p[0])
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	mlog.Info(fmt.Sprintf("delete service_tools ok. name -> [%s].", p[0]))
	return nil
}

func DeleteErrLog() error {
	db, err := ConnDb()
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	deleteSQL := "DELETE FROM error_log WHERE id IN (SELECT id FROM error_log LIMIT 100);"
	stmt, err := db.Prepare(deleteSQL)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	mlog.Info("delete table eRROR_log data ok.")
	return nil
}

func DeleteUserPwd(p ...string) error {
	db, err := ConnDb()
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	deleteSQL := "DELETE FROM user_passwd WHERE serviceName = ?"
	stmt, err := db.Prepare(deleteSQL)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	defer stmt.Close()
	_, err = stmt.Exec(p[0])
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	mlog.Info(fmt.Sprintf("delete user_passwd ok. name -> [%s].", p[0]))
	return nil
}

func DeleteUserAndRole(p ...string) error {
	db, err := ConnDb()
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	deleteSQL := "DELETE FROM user_roles WHERE userName = ?"
	stmt, err := db.Prepare(deleteSQL)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	defer stmt.Close()
	_, err = stmt.Exec(p[0])
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	mlog.Info(fmt.Sprintf("delete user_roles ok. name -> [%s].", p[0]))
	return nil
}

func DeleteUser(p ...string) error {
	db, err := ConnDb()
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	if err := DeleteUserAndRole(p[0]); err != nil {
		mlog.Error(fmt.Sprintf("delete user_roles table data fail -> [%s]", p[0]))
		return err
	}
	deleteSQL := "DELETE FROM user WHERE userName = ?"
	stmt, err := db.Prepare(deleteSQL)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	defer stmt.Close()
	_, err = stmt.Exec(p[0])
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	mlog.Info(fmt.Sprintf("delete user ok. name -> [%s].", p[0]))
	return nil
}

func DeleteUrl(p ...string) error {
	db, err := ConnDb()
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	deleteSQL := "DELETE FROM url WHERE urlName = ?"
	stmt, err := db.Prepare(deleteSQL)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	defer stmt.Close()
	_, err = stmt.Exec(p[0])
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	mlog.Info(fmt.Sprintf("delete url ok. name -> [%s].", p[0]))
	return nil
}

func DeleteUrlType(p ...string) error {
	db, err := ConnDb()
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	deleteSQL := "DELETE FROM url_type WHERE urlType = ?"
	stmt, err := db.Prepare(deleteSQL)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	defer stmt.Close()
	_, err = stmt.Exec(p[0])
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	mlog.Info(fmt.Sprintf("delete url_type ok. name -> [%s].", p[0]))
	return nil
}

func DeleteRolePermission(roleName string) error {
	db, err := ConnDb()
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	deleteSQL := "DELETE FROM roles_permission  WHERE roleName = ?"
	stmt, err := db.Prepare(deleteSQL)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	defer stmt.Close()
	_, err = stmt.Exec(roleName)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	mlog.Info(fmt.Sprintf("delete roles permission ok. role -> [%s].", roleName))
	return nil
}

func DeleteRolePermissionRoute(roleName, permissionRoute string) error {
	db, err := ConnDb()
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	deleteSQL := "DELETE FROM roles_permission  WHERE roleName = ? AND permission = ?"
	stmt, err := db.Prepare(deleteSQL)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	defer stmt.Close()
	_, err = stmt.Exec(roleName, permissionRoute)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	mlog.Info(fmt.Sprintf("delete roles permission ok. role -> [%s], route -> [%s]", roleName, permissionRoute))
	return nil
}

func DeleteRole(p ...string) error {
	db, err := ConnDb()
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	userAndRoleList, err := SelectUserAndRole(fmt.Sprintf("SELECT * FROM user_roles WHERE roleName = \"%s\"", p[0]))
	if err != nil {
		mlog.Error(err.Error())
		return err
	}

	if len(userAndRoleList) >= 1 {
		for _, v := range userAndRoleList {
			mlog.Error(fmt.Sprintf("role binding user, delete fail [%s]->[%s].", v.UserName, v.RoleName))
		}
		return fmt.Errorf("role binding user, delete fail.")
	}

	deleteSQL := "DELETE FROM roles WHERE roleName = ?"
	stmt, err := db.Prepare(deleteSQL)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	defer stmt.Close()
	_, err = stmt.Exec(p[0])
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	if err := DeleteRolePermission(p[0]); err != nil {
		return err
	}
	mlog.Info(fmt.Sprintf("delete roles ok. name -> [%s].", p[0]))
	return nil
}


func DeleteBlacklistIp(blacklistIp string) error {
	db, err := ConnDb()
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	deleteSQL := "DELETE FROM blacklist  WHERE ip = ?"
	stmt, err := db.Prepare(deleteSQL)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	defer stmt.Close()
	_, err = stmt.Exec(blacklistIp)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	mlog.Info(fmt.Sprintf("delete blacklist ip ok. ip -> [%s]", blacklistIp))
	return nil
}


func UpdateUser(p ...string) error {
	db, err := ConnDb()
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	_, err = db.Exec("UPDATE user SET userName = ? WHERE newUserDate = ? AND id = ?", p[0], p[1], p[2])
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	return nil
}

func UpdateUserPwd(p ...string) error {
	db, err := ConnDb()
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	_, err = db.Exec("UPDATE user SET password = ? WHERE userName = ?", p[0], p[1])
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	return nil
}
