package db

import (
	"database/sql"
	"fmt"
	"jwt-auth/models"
	_ "github.com/lib/pq"
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
		"user=%s dbname=%s sslmode=disable", 
		user, dbname,
	)
	// open db
	// added this because the := operator with un-explicitly declared err was leading to a shadowed ACTIve_db
	// basically, it would create a local ACTIVE_DB, which dies after the func returns
	// global ACTIVE_DB would stay unassigned
	var err error 
	ACTIVE_DB, err = sql.Open("postgres", DSN)
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

/*
	- get the user record from the user's table
*/
func GetUserByName(username string) (*models.ServiceUser, error) {

	db := GetDB()
	stmt, err := db.Prepare("SELECT username, password, location, ip_addr FROM USERS WHERE USERNAME = $1")	
	if err != nil{
		return nil, fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close() // very interesting. when you defer it closes the statement when the func ends

	// rows must be queried and scanned into a struct
	// QueryRow returns a non-nil value, always. if scan turns up no data, that is, if there a no rows, then you get an error
	// you gotta scan it into your struct
	var user_data models.ServiceUser
	row := stmt.QueryRow(username)
	err = row.Scan(
		&user_data.User_Name, &user_data.Password, &user_data.Location, &user_data.IP_addr,
	)

	if err != nil{
		return nil, fmt.Errorf("no user found: %v", err)
	}
	return &user_data, nil
	

}

func RegisterUser(newUser models.ServiceUser) error {
	db := GetDB()
	stmt, err := db.Prepare("INSERT INTO USERS (username, password, location, ip_addr) values ($1, $2, $3, $4)")
	if err != nil{
		return fmt.Errorf("failed to prepare statement: %v", err)		
	}
	defer stmt.Close()

	_, err = stmt.Exec(newUser.User_Name, newUser.Password, newUser.Location, newUser.IP_addr)
	if err != nil{
		return fmt.Errorf("failed to save user to db: %v", err)
	}
	return nil
}

func GetSecretKey() ([]byte, error) {
	db := GetDB()
	stmt, err := db.Prepare("SELECT SECRET_KEY FROM secrets where project_name = 'go-jwt-auth'")
	if err != nil{
		return nil, fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	var secretKey string
	err = stmt.QueryRow().Scan(&secretKey)

	if err != nil{
		return nil, fmt.Errorf("secret key not found: %v", err)
	}

	return []byte(secretKey), nil
}