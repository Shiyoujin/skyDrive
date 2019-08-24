package Dao

import (
	"fmt"
	"log"
	"strconv"
)

func OneUp(userID string, fileName string) {

	aaa, _ := strconv.Atoi(userID)
	stmt, err := SqlDB.Prepare(`INSERT into one (userID,fileName) VALUES (?, ?)`)
	check(err)

	res, err := stmt.Exec(aaa, fileName)
	check(err)

	id, err := res.LastInsertId()

	fmt.Println(id)
	stmt.Close()
}

func CheckOne(ownerID int, fileName string) int {

	var id int

	on := strconv.Itoa(ownerID)

	str := "select ifnull((select userID from one where userID = " + on + " AND fileName = " + "\"" + fileName + "\"),0)"

	rows, _ := SqlDB.Query(str)

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

func DeleteOne(userID int, fileName string) {

	stmt, err := SqlDB.Prepare("delete from one where userID=? and fileName = ?")
	check(err)

	res, err := stmt.Exec(userID, fileName)
	check(err)

	affect, err := res.RowsAffected()
	check(err)

	fmt.Println(affect)
}
