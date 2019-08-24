package Dao

import (
	"fmt"
	"log"
	"strconv"
)

type File struct {
	id       int
	userID   int
	fileName string
	public   int
}

//获取public用户权限的校验
func BoolPermit(userID int, fileName string) int {

	var public int

	userIDS := strconv.Itoa(userID)
	rows, err := SqlDB.Query("SELECT public FROM file where userID = " + userIDS + " and fileName =" + "\"" + fileName + "\"")

	check(err)
	if rows.Next() {

		//注意这里的Scan括号中的参数顺序，和 SELECT 的字段顺序要保持一致。
		if err := rows.Scan(&public); err != nil {
			log.Fatal(err)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	rows.Close()
	return public
}

func CFilePermit(userID int, fileName string, public int) int {
	stmt, err := SqlDB.Prepare("UPDATE file SET public = ? WHERE userID=? and fileName = ? ")
	check(err)

	result, err := stmt.Exec(public, userID, fileName)
	num, err := result.RowsAffected()

	return int(num)
}

func AddFile(userID string, fileName string, pp int) {
	aaa, _ := strconv.Atoi(userID)
	stmt, err := SqlDB.Prepare(`INSERT into file (userID,fileName,public) VALUES (?,?,?)`)
	check(err)

	res, err := stmt.Exec(aaa, fileName, pp)
	check(err)

	id, err := res.LastInsertId()

	fmt.Println(id)
	stmt.Close()

}
