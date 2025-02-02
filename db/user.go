package db

import (
	mydb "FILE-SERVER/db/mysql"
	"fmt"
)

// 通过用户名及密码完成user表的注册行为
func UserSignup(username string, passwd string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into tbl_user (`user_name`,`user_pwd`) values (?,?)")
	if err != nil {
		fmt.Println("Failed to insert, err : " + err.Error())
		return false
	}
	defer stmt.Close()
	ret, err := stmt.Exec(username, passwd)
	if err != nil {
		fmt.Println("Failed to insert, err : " + err.Error())
		return false
	}

	if rowsAffected, err := ret.RowsAffected(); nil == err && rowsAffected > 0 {
		return true
	}
	return false
}

// 登录接口
func UserSignin(username string, password string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"select user_pwd from tbl_user " +
			"where user_name = ? " +
			"limit 1")

	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	defer stmt.Close()

	rows, err := stmt.Query(username)
	if err != nil {
		fmt.Println(err.Error())
		return false
	} else if rows == nil {
		fmt.Println("username not found " + username)
		return false
	}

	pRows := mydb.ParseRows(rows)
	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == password {
		return true
	}
	return false
}

// 更新登录用户的token
func UpdateToken(username string, token string) bool {
	// username是unique key，之前有的话，直接覆盖即可
	stmt, err := mydb.DBConn().Prepare(
		"replace into tbl_user_token(`user_name`, `user_token`) " +
			"values(?, ?)")

	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	defer stmt.Close()

	_, err = stmt.Exec(username, token)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

type User struct {
	Username     string
	Email        string
	Phone        string
	SignupAt     string
	LastActiveAt string
	Status       int
}

func GetUserInfo(username string) (User, error) {
	user := User{}

	stmt, err := mydb.DBConn().Prepare(
		"select user_name, signup_at from tbl_user where user_name = ?" +
			" limit 1")

	if err != nil {
		fmt.Println(err.Error())
		return user, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(username).Scan(&user.Username, &user.SignupAt)
	if err != nil {
		fmt.Println(err.Error())
		return user, err
	}

	return user, nil
}
