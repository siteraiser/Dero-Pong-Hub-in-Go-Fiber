package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

func order(iaid string) FullProduct {
	db, err := sql.Open(
		"mysql",
		"root:@tcp(127.0.0.1:3306)/pong_store",
	)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	rows, err := db.Query(
		`
	SELECT *,i_addresses.id ,i_addresses.status AS ia_status 
	FROM i_addresses 
	RIGHT JOIN products ON i_addresses.product_id = products.pid
	RIGHT JOIN users ON i_addresses.user = users.userid
	WHERE i_addresses.id = ? AND i_addresses.user = products.user
		`,
		iaid,
	)
	if err != nil {
		fmt.Printf("%s", err)
	}

	cols, _ := rows.Columns()

	fmt.Println("columns", cols)

	var product FullProduct

	data := make(map[string]string)

	if rows.Next() {
		columns := make([]string, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		rows.Scan(columnPointers...)

		for i, colName := range cols {
			data[colName] = columns[i]
		}
	}

	product.
		Label = data["label"]
	product.
		Inventory = strToInt(
		data["inventory"],
	)
	product.
		Ia_inventory = strToInt(
		data["ia_inventory"],
	)
	product.
		Ia_status = strToInt(
		data["ia_status"],
	)
	product.
		Ia_comment = data["ia_comment"]
	product.
		IAddress = data["iaddr"]

	if data["username"] == "" {
		product.
			Username = data["wallet"]
	} else {
		product.
			Username = data["username"]
	}
	product.Img = ""

	if len(data["image"]) > 22 {
		teststr := data["image"]
		if teststr[0:22] == "data:image/png;base64," {
			product.Img = "<img style='max-width: 100%;' src='" +
				data["image"] +
				"'>"
		}
	}

	if data["image"] == "" {
		product.Image = false
	} else {
		product.Image = true
	}

	ask_amount, _ := strconv.
		ParseFloat(
			data["ask_amount"],
			5,
		)

	product.Ask_amount = ask_amount * .00001

	//See if is a smart contract.
	scid := data["scid"]

	if data["ia_scid"] != "" {
		scid = data["ia_scid"]
	}

	product.Scid = scid

	product.Details = strings.
		Replace(
			data["details"],
			"\n",
			"<br>",
			-1,
		)

	product.P_type = data["p_type"]

	product.Is_physical = false
	if data["p_type"] == "physical" {
		product.Is_physical = true
	}

	stock := ""
	product_inventory := false

	if product.Ia_status != 0 {
		if product.Inventory > 0 {

			product_inventory = true

			stock = "Available:" +
				strconv.
					Itoa(product.Inventory)

		} else if product.Ia_inventory > 0 {
			stock = "Available:" + strconv.Itoa(product.Ia_inventory)

		} else {
			stock = "Out of Stock"
		}
	} else {
		stock = "Item Currently Unavailable."
	}
	product.Stock = stock
	product.Product_inv = product_inventory

	hidden := ""
	if product.P_type == "physical" ||
		stock == "Out of Stock" ||
		stock == "Item Currently Unavailable." {
		hidden = "hidden"
	}

	product.IAClass = hidden

	results2, err := db.Query(
		`
		SELECT id,comment,ia_inventory,ask_amount,i_addresses.status 
		AS ia_status 
		FROM i_addresses 
		WHERE i_addresses.product_id = ? AND i_addresses.user = ?
		`, data["product_id"], data["user"])
	if err != nil {
		panic(err.Error())
	}

	//var iaddresses []interface{}

	for results2.Next() {

		var (
			id           int
			comment      string
			ia_inventory int
			ask_amount   int
			ia_status    int
			iaddress     FullIAddress
		)

		if err := results2.Scan(
			&id,
			&comment,
			&ia_inventory,
			&ask_amount,
			&ia_status,
		); err != nil {
			panic(err)
		}

		iaddress.Id = id
		iaddress.Comment = comment
		iaddress.Ia_inventory = ia_inventory
		iaddress.Ask_amount = ask_amount
		iaddress.Ia_status = ia_status

		class := "greyed_out"
		if product_inventory {
			class = ""
		}
		if iaddress.Ia_inventory > 0 {
			class = ""
		}
		if iaddress.Ia_status == 0 {
			class = "greyed_out"
		}

		selected := false
		if iaddress.Id == strToInt(iaid) {
			selected = true
		}
		iaddress.Selected = selected
		iaddress.Class = class

		if iaddress.Ia_status != 0 ||
			iaddress.Id == strToInt(iaid) {
			product.Iaddresses = append(product.Iaddresses, iaddress)
		}
	}

	return product
}
