package Dao

import (
	"fmt"
	"log"
	"strconv"
)

type User struct {
	id       int
	userID   int
	password string
}

//查询userID是否存在
func IfExistUserID(userID int) int {

	userIDS := strconv.Itoa(userID)
	rows, err := SqlDB.Query("select ifnull((select id from user where userID =" + userIDS + "),0)")
	check(err)

	var user int
	if rows.Next() {
		//注意这里的Scan括号中的参数顺序，和 SELECT 的字段顺序要保持一致。
		if err := rows.Scan(&user); err != nil {
			log.Fatal(err)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	rows.Close()

	return user
}

//注册
func Register(userID int, password string) {

	stmt, err := SqlDB.Prepare(`INSERT into user (userID,password) VALUES (?, ?)`)
	check(err)

	res, err := stmt.Exec(userID, password)
	check(err)

	id, err := res.LastInsertId()
	check(err)

	fmt.Println(id)
	stmt.Close()
}

//登录
func Login(userID int, password string) int {

	user := strconv.Itoa(userID)
	str := "select ifnull((select userID from user where userID = " + user + " AND password = " + "\"" + password + "\"),0)"
	fmt.Println(str)
	rows, err := SqlDB.Query(str)
	check(err)

	var id int

	if rows.Next() {
		//注意这里的Scan括号中的参数顺序，和 SELECT 的字段顺序要保持一致。
		if err := rows.Scan(&id); err != nil {
			log.Fatal(err)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	rows.Close()

	return id
}
