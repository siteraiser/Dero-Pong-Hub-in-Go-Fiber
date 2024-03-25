package main

import (
	"database/sql"
	"fmt"
)

/*
This is our home window
*/

func restricted(uuid string) int { // we are goin to want to interact with output of this function

	// Now open a new connection to the pong_store database
	db, err := sql.Open(
		"mysql",
		"root:@tcp(127.0.0.1:3306)/pong_store",
	)

	if err != nil {
		panic(err.Error())
	}

	// Prepare the SQL query
	query := ("SELECT ia_id FROM orders WHERE uuid = ? AND NOT(ia_id IS NULL)")

	// Get the integrated address associated with the order
	results,
		err := db.Query(
		query,
		uuid,
	)

	fmt.Printf("%v\n", results)

	if err != nil {
		// Handle error
		panic(err)
	}

	for results.Next() {
		var (
			ia_id int
		)
		if err := results.Scan(
			&ia_id,
		); err != nil {
			panic(err)
		}

		return ia_id
	}

	return 0

}
