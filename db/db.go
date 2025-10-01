package db

import (
	"auth-api/models"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

/*
	What this file is for:
	- init func: conn to db and create a db obj
	- func to ret created db obj
*/

var (
	// ACTIVE_DB holds the shared database connection used by the
	// application.
	ACTIVE_DB *sql.DB
	// sqlOpen and prepare are overridable for tests to inject fakes.
	sqlOpen = sql.Open
	prepare = defaultPrepare
)

type rowScanner interface {
	Scan(dest ...any) error
}

type statement interface {
	QueryRow(args ...any) rowScanner
	Exec(args ...any) (sql.Result, error)
	Close() error
}

type sqlStmt struct {
	stmt *sql.Stmt
}

func (s *sqlStmt) QueryRow(args ...any) rowScanner {
	return s.stmt.QueryRow(args...)
}

func (s *sqlStmt) Exec(args ...any) (sql.Result, error) {
	return s.stmt.Exec(args...)
}

func (s *sqlStmt) Close() error {
	return s.stmt.Close()
}

func defaultPrepare(db *sql.DB, query string) (statement, error) {
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	return &sqlStmt{stmt: stmt}, nil
}

// InitDB configures the global database connection pool using the provided
// credentials.
func InitDB(user, dbname, password, host string) error {
	// default host is localhost and port is 5432, i don't need to change that.
	// dsn stands for DATA SOURCE NAME
	DSN := fmt.Sprintf(
		"user=%s dbname=%s password=%v host=%s sslmode=disable",
		user, dbname, password, host,
	)
	// open db
	// added this because the := operator with un-explicitly declared err was leading to a shadowed ACTIve_db
	// basically, it would create a local ACTIVE_DB, which dies after the func returns
	// global ACTIVE_DB would stay unassigned
	var err error
	ACTIVE_DB, err = sqlOpen("postgres", DSN)
	if err != nil {
		return fmt.Errorf("error opening db: %v", err)
	}
	// ping to test
	err = ACTIVE_DB.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping db: %v", err)
	}
	return nil
}

// GetDB returns the active database connection pool.
func GetDB() *sql.DB {
	return ACTIVE_DB
}

// GetUserByName retrieves a user record from the USERS table using the supplied
// username.
func GetUserByName(username string) (*models.ServiceUser, error) {

	db := GetDB()
	stmt, err := prepare(db, "SELECT username, password, location, ip_addr FROM USERS WHERE USERNAME = $1")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close() // very interesting. when you defer it closes the statement when the func ends

	// rows must be queried and scanned into a struct
	// QueryRow returns a non-nil value, always. if scan turns up no data, that is, if there a no rows, then you get an error
	// you gotta scan it into your struct
	var user_data models.ServiceUser
	row := stmt.QueryRow(username)
	err = row.Scan(
		&user_data.Username, &user_data.Password, &user_data.Location, &user_data.IP_addr,
	)

	if err != nil {
		return nil, fmt.Errorf("no user found: %v", err)
	}
	return &user_data, nil

}

// RegisterUser inserts a new user record into the USERS table.
func RegisterUser(newUser models.ServiceUser) error {
	db := GetDB()
	stmt, err := prepare(db, "INSERT INTO USERS (username, password, location, ip_addr) values ($1, $2, $3, $4)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(newUser.Username, newUser.Password, newUser.Location, newUser.IP_addr)
	if err != nil {
		return fmt.Errorf("failed to save user to db: %v", err)
	}
	return nil
}

// GetSecretKey fetches the signing secret for JWT issuance from the secrets
// table.
func GetSecretKey() ([]byte, error) {
	db := GetDB()
	stmt, err := prepare(db, "SELECT SECRET_KEY FROM secrets where project_name = 'go-auth-api'")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	var secretKey string
	err = stmt.QueryRow().Scan(&secretKey)

	if err != nil {
		return nil, fmt.Errorf("secret key not found: %v", err)
	}

	return []byte(secretKey), nil
}
