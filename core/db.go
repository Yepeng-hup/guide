package core

import (
	"database/sql"
	"fmt"
	"guide/global"
	"log"
	"os"
	_ "github.com/mattn/go-sqlite3"
)

var tableList = []string{"cron", "service_tools", "user_passwd", "user", "error_log", "cpu", "mem", "url", "url_type"}

func ConnDb() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "guide.db")
	if err != nil {
		return nil, fmt.Errorf("EEROR: conn db sqlite fail, %v", err.Error())
	}
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
				createTableCron := `CREATE TABLE IF NOT EXISTS cron (id INTEGER PRIMARY KEY, cronNewDate TEXT DEFAULT (strftime('%Y-%m-%d %H:%M', 'now', 'localtime')), cronName TEXT, cronTime TEXT, cronCode TEXT, cronNotes TEXT);`
				_, err = db.Exec(createTableCron)
				if err != nil {
					return fmt.Errorf("ERROR: create table cron fail,%s", err.Error())
				}
			case "service_tools":
				createTableServiceTools := `CREATE TABLE IF NOT EXISTS service_tools (id INTEGER PRIMARY KEY, serviceName TEXT, serviceDate TEXT DEFAULT (strftime('%Y-%m-%d %H:%M', 'now', 'localtime')), startCmd TEXT, serviceNotes TEXT);`
				_, err = db.Exec(createTableServiceTools)
				if err != nil {
					return fmt.Errorf("ERROR: create table service_tools fail,%s", err.Error())
				}
			case "user_passwd":
				createTableUserPasswd := `CREATE TABLE IF NOT EXISTS user_passwd (id INTEGER PRIMARY KEY, serviceName TEXT, user TEXT, password TEXT, Notes TEXT);`
				_, err = db.Exec(createTableUserPasswd)
				if err != nil {
					return fmt.Errorf("ERROR: create table user_passwd fail,%s", err.Error())
				}
			case "user":
				createTableUser := `CREATE TABLE IF NOT EXISTS user (id INTEGER PRIMARY KEY, userName TEXT, newUserDate TEXT DEFAULT (strftime('%Y-%m-%d %H:%M', 'now', 'localtime')), password TEXT);`
				_, err = db.Exec(createTableUser)
				if err != nil {
					return fmt.Errorf("ERROR: create table user fail,%s", err.Error())
				}
			case "error_log":
				createTableErrorLog := `CREATE TABLE IF NOT EXISTS error_log (id INTEGER PRIMARY KEY, newLogDate TEXT DEFAULT (strftime('%Y-%m-%d %H:%M', 'now', 'localtime')), types TEXT, logtext TEXT);`
				_, err = db.Exec(createTableErrorLog)
				if err != nil {
					return fmt.Errorf("ERROR: create table error_log fail,%s", err.Error())
				}
			case "cpu":
				createTableCPU := `CREATE TABLE IF NOT EXISTS cpu (id INTEGER PRIMARY KEY, cpunum INT);`
				_, err = db.Exec(createTableCPU)
				if err != nil {
					return fmt.Errorf("ERROR: create table cpu fail,%s", err.Error())
				}
			case "mem":
				createTableMEM := `CREATE TABLE IF NOT EXISTS mem (id INTEGER PRIMARY KEY, memnum FLOAT);`
				_, err = db.Exec(createTableMEM)
				if err != nil {
					return fmt.Errorf("ERROR: create table mem fail,%s", err.Error())
				}
			case "url":
				createTableUrl := `CREATE TABLE IF NOT EXISTS url (id INTEGER PRIMARY KEY, urlName TEXT, urlAddress TEXT, urlType TEXT, urlNotes TEXT);`
				_, err = db.Exec(createTableUrl)
				if err != nil {
					return fmt.Errorf("ERROR: create table url fail,%s", err.Error())
				}
			case "url_type":
				createTableUrlType := `CREATE TABLE IF NOT EXISTS url_type (id INTEGER PRIMARY KEY, urlType TEXT);`
				_, err = db.Exec(createTableUrlType)
				if err != nil {
					return fmt.Errorf("ERROR: create table url_type fail,%s", err.Error())
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
		log.Println("ERROR: init fail.")
		log.Println("ERROR: query user table fail,", err.Error())
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
			log.Println("ERROR: init fail.")
			log.Println(err.Error())
			os.Exit(1)
		}
		userList = append(userList, user)
	}

	if len(userList) < 1 {
		insertSQL := `INSERT INTO user (userName, password) VALUES (?,?);`
		encryptionPwd, err := PasswordEncryption("guide654321", global.NowKey)
		if err != nil {
			log.Println("ERROR: init user encryption passwd fail, ", err.Error())
			return
		}
		_, err = db.Exec(insertSQL, "admin", encryptionPwd)
		if err != nil {
			log.Println("ERROR: init db user fail, ", err.Error())
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
		log.Println("ERROR: insert data to cron fail, ", err.Error())
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
		return fmt.Errorf("ERROR: insert data to service_tools fail,%s", err.Error())
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
		return fmt.Errorf("ERROR: insert data to user_passwd fail,%s", err.Error())
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
		return fmt.Errorf("ERROR: insert data to user fail,%s", err.Error())
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
		return fmt.Errorf("ERROR: insert data error_log fail,%s", err.Error())
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
		return fmt.Errorf("ERROR: insert data cpu fail,%s", err.Error())
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
		return fmt.Errorf("ERROR: insert data mem fail,%s", err.Error())
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
		return fmt.Errorf("ERROR: insert data url fail,%s", err.Error())
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
		return fmt.Errorf("ERROR: insert data url_type fail,%s", err.Error())
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
			return nil, fmt.Errorf("ERROR: query cron table fail,%s", err.Error())
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
			return nil, fmt.Errorf("ERROR: query cron table fail,%s", err.Error())
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
		return nil, fmt.Errorf("ERROR: query service_tools table fail,%s", err.Error())
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
		return nil, fmt.Errorf("ERROR: query user_passwd table fail,%s", err.Error())
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
		return nil, fmt.Errorf("ERROR: query user table fail,%s", err.Error())
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

func SelectLog(selectSql string) ([]ErrorLog, error) {
	db, err := ConnDb()
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	rows, err := db.Query(selectSql)
	if err != nil {
		return nil, fmt.Errorf("ERROR: query error_log table fail,%s", err.Error())
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
		return nil, fmt.Errorf("ERROR: query cpu table fail,%s", err.Error())
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
		return nil, fmt.Errorf("ERROR: query mem table fail,%s", err.Error())
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
		return nil, fmt.Errorf("ERROR: query url table fail,%s", err.Error())
	}
	defer rows.Close()
	var url struct {
		// Id         string
		UrlName       string
		UrlAddress    string
		UrlType 	  string
		UrlNotes	  string
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
		return nil, fmt.Errorf("ERROR: query typeName table fail,%s", err.Error())
	}
	defer rows.Close()
	var urlType struct {
		// Id         string
		TypeName       string
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
	log.Printf("INFO: delete cron database record ok. name -> [%s].", p[0])
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
	log.Printf("INFO: delete service_tools ok. name -> [%s].", p[0])
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
	log.Printf("INFO: delete table error_log data ok.")
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
	log.Printf("INFO: delete user_passwd ok. name -> [%s].", p[0])
	return nil
}

func DeleteUser(p ...string) error {
	db, err := ConnDb()
	if err != nil {
		return fmt.Errorf(err.Error())
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
	log.Printf("INFO: delete user ok. name -> [%s].", p[0])
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
	log.Printf("INFO: delete url ok. name -> [%s].", p[0])
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
	log.Printf("INFO: delete url_type ok. name -> [%s].", p[0])
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

