package db

import (
	"database/sql"
	"fmt"
)

/*
	What this file is for:
	- init func: conn to db and create a db obj
	- func to ret created db obj
	-
*/

var ACTIVE_DB *sql.DB

func InitDB(user, dbname string) error {
	// default host is localhost and port is 5432, i don't need to change that.
	// dsn stands for DATA SOURCE NAME
	DSN := fmt.Sprintf(
		"user=%s dbname=%s sslmode=verify-full", 
		user, dbname,
	)
	// open db
	ACTIVE_DB, err := sql.Open("postgres", DSN)
	if err != nil{
		return fmt.Errorf("error opening db: %v", err)
	}
	// ping to test
	err = ACTIVE_DB.Ping()
	if err != nil{
		return fmt.Errorf("failed to ping db: %v", err)
	}
	return nil
}

// return db
func GetDB() *sql.DB {
	return ACTIVE_DB
}