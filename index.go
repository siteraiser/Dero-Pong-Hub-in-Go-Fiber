package main

import (
	"database/sql"
	"fmt"
	"time"
)

/*
This is our home window
*/

func index() []interface{} { // we are goin to want to interact with output of this function

	// Now open a new connection to the pong_store database
	db, err := sql.Open(
		"mysql",
		"root:@tcp(127.0.0.1:3306)/pong_store",
	)

	if err != nil {
		panic(err.Error())
	}
	// Calculate the time 'count' minutes ago
	now := time.Now().UTC()
	count := 30
	then := now.Add(
		-time.Duration(count) * time.Minute,
	)

	// Prepare the SQL query
	query := ("SELECT pid,label,details,inventory,user,image,username " +
		"FROM products " +
		"INNER JOIN users " +
		"ON products.user = users.userid " +
		"WHERE users.checkin > ? " +
		"ORDER BY products.id DESC")

	// Execute the query with the calculated time
	results,
		err := db.Query(
		query,
		then.Format(
			"2006-01-02 15:04:05",
		),
	)

	fmt.Printf("%v\n", results)

	if err != nil {
		// Handle error
		panic(err)
	}

	var products []interface{}

	for results.Next() {
		var (
			pid       int
			label     string
			details   string
			inventory int
			user      string
			image     string
			username  string
			product   Product
		)
		if err := results.Scan(
			&pid,
			&label,
			&details,
			&inventory,
			&user,
			&image,
			&username,
		); err != nil {
			panic(err)
		}

		product.Label = label

		if len(label) > 50 {
			product.Label = label[0:50]
		}

		product.Details = details

		if len(details) > 100 {
			product.Details = details[0:50] + "..."
		}

		product.Inventory = inventory
		product.User = user

		product.Img = ""
		if len(image) > 22 {
			teststr := image
			if teststr[0:22] == "data:image/png;base64," {
				product.Img = "<img class='product_image' src='" +
					image +
					"'>"
			}
		}

		if image == "" {
			product.Image = false
		} else {
			product.Image = true
		}

		product.Username = username
		if username == "" {
			product.Username = "Username Not Provided"
		}

		use_ia_inventory := false
		if product.Inventory == 0 {
			use_ia_inventory = true
		}

		product.Use_p_inv = true
		if use_ia_inventory {
			product.Use_p_inv = false
		}

		// Execute the query to retrieve information from the 'i_addresses' table
		results2,
			err := db.Query(
			`
    SELECT id, comment, ia_inventory, ask_amount, status
    FROM i_addresses
    WHERE product_id = ? AND user = ?
			`,
			pid,
			user,
		)

		if err != nil {
			panic(err.Error())
		}

		for results2.Next() {
			var (
				id           int
				comment      string
				ia_inventory int
				ask_amount   int
				status       int
				iaddress     IAddress
			)

			if err := results2.Scan(
				&id,
				&comment,
				&ia_inventory,
				&ask_amount,
				&status,
			); err != nil {
				panic(err)
			}

			iaddress.
				Id = id
			iaddress.
				Comment = comment
			iaddress.
				Ia_inventory = ia_inventory
			iaddress.
				Ask_amount = float64(ask_amount) * .00001
			iaddress.
				Status = status

			skip_ia := false
			if (use_ia_inventory &&
				iaddress.Ia_inventory < 1) ||
				iaddress.Status == 0 {
				skip_ia = true
			}

			if !use_ia_inventory {
				iaddress.Ia_inventory = 0
			}

			if !skip_ia {
				product.Iaddresses = append(
					product.Iaddresses,
					iaddress,
				)
			}
		}

		if len(product.Iaddresses) > 0 {
			products = append(products, product)
		}
	}

	return products

}
