package usermodel

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

func CheckUser(user string, pass string) map[string]interface{} {
	userRecord := getUserByUUID(pass)

	if len(userRecord) != 0 && userRecord["username"] != user {
		updateUsername(user, pass)
	}
	return userRecord
}

func getUserByUUID(uuid string) map[string]interface{} {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/pong_store")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM users WHERE uuid = ?", uuid)
	if err != nil {
		panic(err)
	}
	columns, err := rows.Columns()
	if err != nil {
		panic(err)
	}

	result := make(map[string]interface{})
	values := make([]interface{}, len(columns))
	for i := range values {
		values[i] = new(interface{})
	}

	for rows.Next() {
		err := rows.Scan(values...)
		if err != nil {
			panic(err)
		}
		for i, col := range columns {
			result[col] = *(values[i].(*interface{}))
		}
	}

	return result
}
func getUserByWallet(wallet string) map[string]interface{} {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/pong_store")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM users WHERE wallet = ?", wallet)
	if err != nil {
		panic(err)
	}
	columns, err := rows.Columns()
	if err != nil {
		panic(err)
	}

	result := make(map[string]interface{})
	values := make([]interface{}, len(columns))
	for i := range values {
		values[i] = new(interface{})
	}

	for rows.Next() {
		err := rows.Scan(values...)
		if err != nil {
			panic(err)
		}
		for i, col := range columns {
			result[col] = *(values[i].(*interface{}))
		}
	}

	return result
}
func InsertUser(username string, wallet string) string {
	if len(getUserByWallet(wallet)) != 0 {
		return "failed"
	}
	fmt.Printf("New User: %s\n", username)

	var uid = uuid.New().String()

	currentTime := time.Now().UTC()

	fmt.Printf("time: %s\n", currentTime.Format("2006-01-02 15:04:05"))

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/pong_store")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	result, err := db.Exec("INSERT INTO users (username,wallet,uuid,status,checkin)VALUES(?,?,?,?,?)", username, wallet, uid, 1, currentTime.Format("2006-01-02 15:04:05"))
	if err != nil {
		panic(err)
	}

	lastid, err := result.LastInsertId()
	if err != nil || lastid < 1 {
		return "failed"
	}

	return string(uid)
}

func updateUsername(username string, uuid string) {

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/pong_store")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	_, err = db.Exec("UPDATE users SET username=? WHERE uuid=?", username, uuid)
	if err != nil {
		panic(err)
	}

}

func CheckIn(uuid string) {

	currentTime := time.Now().UTC()

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/pong_store")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	_, err = db.Exec("UPDATE users SET checkin=? WHERE uuid=?", currentTime.Format("2006-01-02 15:04:05"), uuid)
	if err != nil {
		panic(err)
	}

}
