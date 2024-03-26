package productmodel

import (
	"database/sql"
	"reflect"
)

func NewTx(
	user map[string]interface {
	},
	tx TxJSON,
) bool {

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/pong_store")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	//User is the seller's user id
	if reflect.ValueOf(tx.Ia_id).IsZero() {
		result, err := db.Exec(`INSERT INTO orders (uuid,user) VALUES(?,?)`, tx.Uuid, user["userid"])
		if err != nil {
			panic(err)
		}

		lastid, err := result.LastInsertId()
		if err != nil || lastid < 1 {
			return false
		}

	} else {
		result, err := db.Exec(`INSERT INTO orders (uuid,ia_id,user) VALUES(?,?,?)`, tx.Uuid, tx.Ia_id, user["userid"])
		if err != nil {
			panic(err)
		}

		lastid, err := result.LastInsertId()
		if err != nil || lastid < 1 {
			return false
		}
	}

	return true

}
