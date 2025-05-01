package db

import (
	"database/sql"
	"fmt"
	"jwt-auth/models"
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


/*
	WHATS NEXT?
	- query/select db for an existing user. register and login both need to do this
		- the data the gets returned, need to get marshalled into JSON
		- the data has to fit the existing structure of MasterUserDB ( map[string]ServiceUser )
		- basically a swap-in
		- what would be a good time to load the database into masteruserdb for the first time?
			- this seems like an important question!
			- or should it even be done? just load one item per query. no need to hold it all in mem, dumbo!
	- create new/insert user function
	- create new/insert jwt function
	- also, modify tests user/jwt to do testing with db

*/

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
	}else{
		return &user_data, nil
	}

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