package db

import (
	mydb "cloud/db/mysql"
	"fmt"
)

// UserSignup 通过用户名及密码完成user表的注册操作
func UserSignup(username string, passwd string) bool {
	stmt, err := mydb.DBConn().Prepare("insert ignore into tbl_user(`user_name`,`user_pwd`)values (?,?)")
	if err != nil {
		fmt.Printf("Failed to create statement, err:%v\n", err)
		return false
	}
	defer stmt.Close()
	ret, err := stmt.Exec(username, passwd)
	if err != nil {
		fmt.Printf("Failed to insert, err:%v\n", err)
		return false
	}
	rows, err := ret.RowsAffected()
	if err == nil && rows > 0 {
		return true
	}
	return false
}

// UserSignin 判断密码是否一一致
func UserSignin(username string, encpwd string) bool {

	stmt, err := mydb.DBConn().Prepare("select * from tbl_user where `user_name`=? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	row, err := stmt.Query(username)
	if err != nil {
		fmt.Println(err.Error())
		return false
	} else if row == nil {
		fmt.Println("username not found:" + username)
		return false
	}
	pRows := mydb.ParseRows(row)
	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == encpwd {
		return true
	}
	return false
}

// UpdateToken 刷新用户登录的token
func UpdateToken(username string, token string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"replace into tbl_user_token(`user_name`,`user_token`)values (?,?)")
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
	Username     string `json:"Username"`
	Email        string `json:"Email"`
	Phone        string `json:"Phone"`
	SignupAt     string `json:"SignupAt"`
	LastActiveAt string `json:"LastActiveAt"`
	Status       string `json:"Status"`
}

func GetUserInfo(username string) (User, error) {
	user := User{}

	stmt, err := mydb.DBConn().Prepare("select user_name,signup_at from tbl_user where `user_name`=? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return user, err
	}
	defer stmt.Close()
	// 执行查询的操作
	err = stmt.QueryRow(username).Scan(&user.Username, &user.SignupAt)
	if err != nil {
		fmt.Println(err.Error())
		return user, err
	}
	return user, nil
}
