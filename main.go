package main

import (
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Open a connection to the MySQL server
	db, err := connectDB("root:@tcp(127.0.0.1:3306)/")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Create the database if it doesn't exist
	createDatabase(db, "pong_store")

	// Open a new connection to the pong_store database
	db, err = connectDB("root:@tcp(127.0.0.1:3306)/pong_store")
	if err != nil {
		panic(err.Error())
	}

	// Create tables if they don't exist
	createTables(db)

	// Initialize Fiber app
	app := initializeApp()

	// Start the server
	app.Listen(":3000")

}
