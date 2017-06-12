package db

import (
	"database/sql"

	"fmt"

	// use mysql
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
)

var (
	db       *sql.DB
	dbUser   string
	dbPasswd string
	dbDirver string
	dbDBName string
	dbHost   string
	dbPort   string
)

func tryInitDBInstance() *sql.DB {
	if nil == db {
		err := initDB("poker", "poker1111", "mysql", "poker", "127.0.0.1", "3306")
		if nil != err {
			glog.Fatal("DB init failed:", err)
		}
	}

	return db
}

// initDB new dababase
func initDB(user string, passwd string, dirver string, database string, host string, port string) error {
	//need check params?
	dbUser = user
	dbPasswd = passwd
	dbDirver = dirver
	dbDBName = database
	dbHost = host
	dbPort = port

	dbInstance, err := open()
	if nil == err {
		db = dbInstance
	}

	return err
}

func open() (*sql.DB, error) {
	// "user:password@tcp(127.0.0.1:3306)/database"
	// "user:password@tcp(host:port)/database"
	s := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPasswd, dbHost, dbPort, dbDBName)
	d, err := sql.Open(dbDirver, s)
	return d, err
}

// Query query sql
func Query(statement string, args ...string) (*sql.Rows, error) {
	tryInitDBInstance()
	stmtOut, err := db.Prepare(statement)
	if nil != err {
		glog.Errorf("SQL prepare (%v) failed: %v", statement, err)
		return nil, err
	}
	defer stmtOut.Close()

	var rows *sql.Rows
	if 0 == len(args) {
		rows, err = stmtOut.Query()
	} else {
		rows, err = stmtOut.Query(args)
	}
	if nil != err {
		glog.Errorf("SQL (%v) query (%v) failed: %v", statement, args, err)
	}

	return rows, err
}

// Count count
func Count(statement string, args ...string) (uint64, error) {
	tryInitDBInstance()
	stmtOut, err := db.Prepare(statement)
	if nil != err {
		glog.Errorf("SQL prepare (%v) failed: %v", statement, err)
		return 0, err
	}
	defer stmtOut.Close()

	var count uint64
	if 0 == len(args) {
		stmtOut.QueryRow().Scan(&count)
	} else {
		stmtOut.QueryRow(args).Scan(&count)
	}

	return count, err
}

// Exec execute sql
func Exec(statement string, args ...string) (bool, error) {
	tryInitDBInstance()
	stmtOut, err := db.Prepare(statement)
	if nil != err {
		glog.Errorf("SQL prepare (%v) failed: %v", statement, err)
		return false, err
	}
	defer stmtOut.Close()

	if 0 == len(args) {
		_, err = stmtOut.Exec()
	} else {
		_, err = stmtOut.Exec(args)
	}
	if nil != err {
		glog.Errorf("SQL (%v) query (%v) failed: %v", statement, args, err)
		return false, err
	}

	return false, nil
}
