package productmodel

import "database/sql"

//Integrated Addresses

func SubmitIAddress(user map[string]interface{}, iaddress IAddressJSON) bool {

	if len(getIAddressById(iaddress.Id, user)) == 0 {
		insertIAddress(user, iaddress)
	} else if iaddress.Action == "delete" {

		deleteIAddress(user, iaddress)

	} else {
		updateIAddress(user, iaddress)
	}

	return true
}
func insertIAddress(user map[string]interface{}, iaddress IAddressJSON) bool {

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/pong_store")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	result, err := db.Exec(
		"INSERT INTO i_addresses (iaddr_id,user,product_id,iaddr,ask_amount,comment,ia_scid,status,ia_inventory)VALUES(?,?,?,?,?,?,?,?,?)",
		iaddress.Id, user["userid"], iaddress.Product_id, iaddress.Iaddr, iaddress.Ask_amount, iaddress.Comment, iaddress.Ia_scid, iaddress.Status, iaddress.Ia_inventory)
	if err != nil {
		panic(err)
	}

	lastid, err := result.LastInsertId()
	if err != nil || lastid < 1 {
		return false
	}

	return true
}

func getIAddressById(iaid int, user map[string]interface{}) map[string]interface{} {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/pong_store")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM i_addresses WHERE iaddr_id = ? AND user = ?", iaid, user["userid"])
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

func updateIAddress(user map[string]interface{}, iaddress IAddressJSON) bool {

	if !sameAddress(iaddress.Iaddr, user) {
		return false
	}

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/pong_store")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	_, err = db.Exec(
		"UPDATE i_addresses  SET ia_scid=?,status=?,ia_inventory=? WHERE iaddr=?",
		iaddress.Ia_scid, iaddress.Status, iaddress.Ia_inventory, iaddress.Iaddr)
	return err == nil
}

func deleteIAddress(user map[string]interface{}, iaddress IAddressJSON) bool {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/pong_store")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	_, err = db.Exec(
		"DELETE FROM i_addresses WHERE iaddr_id = ? AND user = ?",
		iaddress.Id, user["userid"])
	return err == nil
}
