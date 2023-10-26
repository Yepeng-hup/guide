package core

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func ConnDb()(*sql.DB,error){
	db, err := sql.Open("sqlite3", "guide.db")
	if err != nil {
		return nil, fmt.Errorf("EEROR: conn db sqlite fail, %v", err.Error())
	}
	return db, nil
}


func CreateCronTable()error{
	db, err := ConnDb()
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	createTableSQL := `CREATE TABLE IF NOT EXISTS cron (id INTEGER PRIMARY KEY, cronNewDate TEXT DEFAULT (strftime('%Y-%m-%d %H:%M', 'now', 'localtime')), cronName TEXT, cronTime TEXT, cronCode TEXT, cronNotes TEXT);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("ERROR: create table cron fail,%s", err.Error())
	}
	return nil
}

func InsertAct(params ...string){
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


func SelectAct(k,v string, b bool)([]Cron, error){
	db, err := ConnDb()
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	if b {
		rows, err := db.Query("SELECT * FROM cron")
		if err != nil {
			return nil, fmt.Errorf("ERROR: query cron table fail,%s",err.Error())
		}
		defer rows.Close()
		var cron struct{
			CronId string
			CronNewDate string
			CronName string
			CronTime string
			CronCode string
			CronNotes string
		}
		cronList := make([]Cron, 0)
		for rows.Next() {
			err := rows.Scan(&cron.CronId,&cron.CronNewDate,&cron.CronName, &cron.CronTime, &cron.CronCode, &cron.CronNotes)
			if err != nil {
				return nil, fmt.Errorf(err.Error())
			}
			cronList = append(cronList, cron)
		}
		return cronList, nil

	}else {
		ssql := "SELECT * FROM cron  WHERE ? = ?"
		rows, err := db.Query(ssql, k, v)
		if err != nil {
			return nil, fmt.Errorf("ERROR: query cron table fail,%s",err.Error())
		}
		defer rows.Close()
		var cron struct{
			CronId string
			CronNewDate string
			CronName string
			CronTime string
			CronCode string
			CronNotes string
		}
		cronList := make([]Cron, 0)
		for rows.Next() {
			err := rows.Scan(&cron.CronId,&cron.CronNewDate,&cron.CronName, &cron.CronTime, &cron.CronCode, &cron.CronNotes)
			if err != nil {
				return nil, fmt.Errorf(err.Error())
			}
			cronList = append(cronList, cron)
		}
		return cronList, nil
	}
}


func DeleteAct(p ...string)error{
	db, err := ConnDb()
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	deleteSQL := "DELETE FROM cron WHERE cronName = ?"

	// use SQL
	stmt, err := db.Prepare(deleteSQL)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	defer stmt.Close()
	// use sql and afferent delete parameter
	_, err = stmt.Exec(p[0])
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	log.Printf("INFO: delete cron ok. name -> [%s].", p[0])
	return nil
}