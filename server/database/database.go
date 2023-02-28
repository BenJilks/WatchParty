package main

import "database/sql"

type Video struct {
}

func OpenDatabase() (*sql.DB, error) {
	return sql.Open("sqlite3", "file:db.sqlite")
}
