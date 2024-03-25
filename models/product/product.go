package productmodel

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func SubmitProduct(user map[string]interface{}, product ProductJSON) bool {
	// reflect
	// _, found := reflect.TypeOf(Product{}).FieldByName("Action")

	if len(getProductById(product.Id, user)) == 0 {
		fmt.Printf("Inserting\n")
		insertProduct(user, product)
	} else if product.Action == "delete" {

		deleteProduct(user, product)

	} else {
		fmt.Printf("Updating\n")
		updateProduct(user, product)
	}

	return true
}

func insertProduct(user map[string]interface{}, product ProductJSON) bool {

	db, err := sql.Open(
		"mysql",
		"root:@tcp(127.0.0.1:3306)/pong_store",
	)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	result, err := db.Exec(
		"INSERT INTO products (pid,user,p_type,label,details,scid,inventory,image)VALUES(?,?,?,?,?,?,?,?)",
		product.Id,
		user["userid"],
		product.P_type,
		product.Label,
		product.Details,
		product.Scid,
		product.Inventory,
		product.Image,
	)
	if err != nil {
		panic(err)
	}

	lastid, err := result.LastInsertId()
	// fmt.Printf("%s", result)
	if err != nil || lastid < 1 {
		return false
	}

	return true
}

func getProductById(pid int, user map[string]interface{}) map[string]interface{} {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/pong_store")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM products WHERE pid = ? AND user = ?", pid, user["userid"])
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

func updateProduct(user map[string]interface{}, product ProductJSON) bool {

	db, err := sql.Open(
		"mysql",
		"root:@tcp(127.0.0.1:3306)/pong_store",
	)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	if reflect.ValueOf(product.Image).IsZero() {
		_, err = db.Exec(
			"UPDATE products  SET p_type=?,label=?,details=?,scid=?,inventory=? WHERE pid=? AND user=?",
			product.P_type,
			product.Label,
			product.Details,
			product.Scid,
			product.Inventory,
			product.Id,
			user["userid"])
	} else {
		_, err = db.Exec(
			"UPDATE products  SET p_type=?,label=?,details=?,scid=?,inventory=?,image=? WHERE pid=? AND user=?",
			product.P_type,
			product.Label,
			product.Details,
			product.Scid,
			product.Inventory,
			product.Image,
			product.Id,
			user["userid"])
	}

	if err != nil {
		return false
	}
	return true
}

func deleteProduct(user map[string]interface{}, product ProductJSON) bool {
	fmt.Printf("Deleting%v", product.Id)
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/pong_store")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	_, err = db.Exec(
		"DELETE FROM i_addresses WHERE product_id = ? AND user = ?",
		product.Id, user["userid"])
	if err != nil {
		return false
	}
	_, err = db.Exec(
		"DELETE FROM products WHERE pid = ? AND user = ?",
		product.Id, user["userid"])
	if err != nil {
		return false
	}

	return true
}
