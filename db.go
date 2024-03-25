package main

import "database/sql"

// createDatabase creates a database if it doesn't exist
func createDatabase(db *sql.DB, dbName string) {
	_, err := db.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
	if err != nil {
		panic(err.Error())
	}
}

// connectDB opens a connection to the MySQL server
func connectDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// createTables creates required tables if they don't exist
func createTables(db *sql.DB) {
	tables := []struct {
		name  string
		query string
	}{
		{"products", productsTable},
		{"i_addresses", iAddressesTable},
		{"users", usersTable},
		{"orders", ordersTable},
	}

	for _, table := range tables {
		createTable(db, table.query)
	}
}

// createTable creates a table if it doesn't exist
func createTable(db *sql.DB, query string) {
	_, err := db.Exec(query)
	if err != nil {
		panic(err.Error())
	}
}
