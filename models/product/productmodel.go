package productmodel

import (
	"database/sql"
	"fmt"

	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// "reflect"
type Product struct {
	Id        int
	P_type    string `json:"p_type,omitempty"`
	Label     string `json:"label,omitempty"`
	Details   string `json:"details,omitempty"`
	Scid      string `json:"scid,omitempty"`
	Inventory int    `json:"inventory,omitempty"`
	Image     string `json:"image,omitempty"`
	Action    string `json:"action,omitempty"`
}

type IAddress struct {
	Id           int
	Product_id   int    `json:"product_id,omitempty"`
	Iaddr        string `json:"iaddr,omitempty"`
	Ask_amount   int    `json:"ask_amount,omitempty"`
	Comment      string `json:"comment,omitempty"`
	Ia_scid      string `json:"ia_scid,omitempty"`
	Status       int    `json:"status,omitempty"`
	Ia_inventory int    `json:"ia_inventory,omitempty"`
	Action       string `json:"action,omitempty"`
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

func SubmitProduct(user map[string]interface{}, product Product) bool {
	//_, found := reflect.TypeOf(product).FieldByName("Action")

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

func insertProduct(user map[string]interface{}, product Product) bool {

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/pong_store")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	result, err := db.Exec(
		"INSERT INTO products (pid,user,p_type,label,details,scid,inventory,image)VALUES(?,?,?,?,?,?,?,?)",
		product.Id, user["userid"], product.P_type, product.Label, product.Details, product.Scid, product.Inventory, product.Image)
	if err != nil {
		panic(err)
	}

	lastid, err := result.LastInsertId()
	if err != nil || lastid < 1 {
		return false
	}

	return true
}

func updateProduct(user map[string]interface{}, product Product) bool {

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/pong_store")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	_, err = db.Exec(
		"UPDATE products  SET p_type=?,label=?,details=?,scid=?,inventory=?,image=? WHERE pid=? AND user=?",
		product.P_type, product.Label, product.Details, product.Scid, product.Inventory, product.Image, product.Id, user["userid"])
	if err != nil {
		return false
	}
	return true
}

func deleteProduct(user map[string]interface{}, product Product) bool {
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

//Integrated Addresses

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

func SubmitIAddress(user map[string]interface{}, iaddress IAddress) bool {

	if len(getIAddressById(iaddress.Id, user)) == 0 {
		insertIAddress(user, iaddress)
	} else if iaddress.Action == "delete" {

		deleteIAddress(user, iaddress)

	} else {
		updateIAddress(user, iaddress)
	}

	return true
}
func insertIAddress(user map[string]interface{}, iaddress IAddress) bool {

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

func updateIAddress(user map[string]interface{}, iaddress IAddress) bool {

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

func deleteIAddress(user map[string]interface{}, iaddress IAddress) bool {
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

func sameAddress(iaddr string, user map[string]interface{}) bool {

	//fmt.Printf("s%", user)
	//fmt.Printf("s%", iaddr)
	waddr := B2S(user["wallet"].([]uint8))

	ia_frag := strings.TrimLeft(iaddr, "deroi1")
	wa_frag := strings.TrimLeft(waddr, "dero1")
	ia_frag = ia_frag[0:53]
	wa_frag = wa_frag[0:53]

	return ia_frag == wa_frag
}

/**/
func B2S(bs []uint8) string {
	ba := []byte{}
	for _, b := range bs {
		ba = append(ba, byte(b))
	}
	return string(ba)
}
