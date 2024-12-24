package main

import (
	"context"
	"database/sql"
	"fmt"
)

var db *sql.DB

func InitDatabase(dbPath string) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}
	_, err = db.ExecContext(
		context.Background(),
		`CREATE TABLE IF NOT EXISTS location 
			(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				city TEXT NOT NULL,
				country TEXT NOT NULL
			)`,
	)
	if err != nil {
		return fmt.Errorf("error creating location table: %v", err)
	}
	return nil
}

func InsertLocation(city, country string) error {
	_, err := db.ExecContext(
		context.Background(),
		"INSERT INTO location (city, country) VALUES (?, ?)",
		city, country,
	)
	if err != nil {
		return fmt.Errorf("error inserting location: %v", err)
	}
	return nil
}
