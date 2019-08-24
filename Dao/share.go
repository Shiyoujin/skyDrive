package Dao

import (
	"fmt"
	"log"
	"strconv"
)

type Share struct {
	userID   int
	photo    string
	shareUrl string
	dateNum  int
}

func AddShare(userID string, photo string, shareUrl string, ExtractedCode string, dateNum int, createdAt string) {

	aaa, _ := strconv.Atoi(userID)
	stmt, err := SqlDB.Prepare(`INSERT into share (userID,fileName,afterParams,ExtractedCode,dateNum,createdAt) VALUES (?, ?, ?, ?, ?, ?)`)
	check(err)

	res, err := stmt.Exec(aaa, photo, shareUrl, ExtractedCode, dateNum, createdAt)
	check(err)

	id, err := res.LastInsertId()

	fmt.Println(id)
	stmt.Close()
}

func GetFileName(afterparams string) (string, string, int, string) {

	var fileName string

	var ExtractedCode string

	var dateNum int

	var createdAt string

	str := "SELECT fileName,ExtractedCode,dateNum,createdAt FROM share where  afterparams =" + "\"" + afterparams + "\""
	rows, err := SqlDB.Query(str)

	check(err)
	if rows.Next() {

		//注意这里的Scan括号中的参数顺序，和 SELECT 的字段顺序要保持一致。
		if err := rows.Scan(&fileName, &ExtractedCode, &dateNum, &createdAt); err != nil {
			log.Fatal(err)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	rows.Close()
	return fileName, ExtractedCode, dateNum, createdAt
}

//取出fileName
func CheckCode(fileName string, ExtractedCode string) (int, int) {

	var id int
	var userID int

	str := "SELECT id,userID FROM share where fileName = " + fileName + " and ExtractedCode =" + ExtractedCode
	fmt.Println(str)
	rows, err := SqlDB.Query(str)

	check(err)
	if rows.Next() {

		//注意这里的Scan括号中的参数顺序，和 SELECT 的字段顺序要保持一致。
		if err := rows.Scan(&id, &userID); err != nil {
			log.Fatal(err)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	rows.Close()
	return id, userID

}

func check(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
