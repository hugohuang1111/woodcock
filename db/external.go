package db

import (
	"errors"

	"github.com/golang/glog"
)

// UserRegister user register
func UserRegister(name, passwd string) (bool, string) {
	n := "'" + name + "'"
	pw := "'" + passwd + "'"
	sql := newstatement().insert("user").columns("name", "passwd").values(n, pw).toString()
	suc, err := Exec(sql)
	if nil == err {
		return suc, ""
	}

	return suc, err.Error()
}

//UserLogin user login
func UserLogin(name, passwd string) (bool, uint64) {
	sql := newstatement().selects("id").from("user").where("name", name, false).where("passwd", passwd, false).toString()
	rows, err := Query(sql)
	if nil != err {
		glog.Errorf("UserLogin sql fail:%v", err)
		return false, 0
	}
	defer rows.Close()

	var uid uint64
	if rows.Next() {
		err = rows.Scan(&uid)
		if nil != err {
			glog.Errorf("UserLogin failed:%v", err)
			return false, 0
		}
	} else {
		return false, 0
	}

	return true, uid
}

// UserCount get user count
func UserCount() (uint64, error) {
	sql := newstatement().count("*").from("user").toString()
	count, err := Count(sql)

	return count, err
}

// QueryUserInfo query user info
func QueryUserInfo(uid uint64) (name string, err error) {
	sql := newstatement().selects("id", "name").from("user").where("id", string(uid), false).toString()
	rows, e := Query(sql)
	if nil != e {
		err = e
		return
	}
	defer rows.Close()

	var userID uint64
	var userName string
	if rows.Next() {
		err = rows.Scan(&userID, &userName)
		if nil != err {
			return
		}
	} else {
		err = errors.New("not found user")
		return
	}

	return userName, nil
}

func QueryUserIDByName(name, password string) (name string, err error) {
	sql := newstatement().selects("id", "name").from("user").where("id", string(uid), false).toString()
	rows, e := Query(sql)
	if nil != e {
		err = e
		return
	}
	defer rows.Close()

	var userID uint64
	var userName string
	if rows.Next() {
		err = rows.Scan(&userID, &userName)
		if nil != err {
			return
		}
	} else {
		err = errors.New("not found user")
		return
	}

	return userName, nil
}
