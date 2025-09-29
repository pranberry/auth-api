package db

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"regexp"
	"testing"

	"jwt-auth/models"
)

// stubConn implements the minimum driver.Conn interface needed for exercising
// sql.OpenDB while allowing us to inject ping behavior.
type stubConn struct {
	pingErr error
}

func (s *stubConn) Prepare(query string) (driver.Stmt, error) {
	return nil, errors.New("not implemented")
}

func (s *stubConn) Close() error {
	return nil
}

func (s *stubConn) Begin() (driver.Tx, error) {
	return nil, errors.New("not implemented")
}

func (s *stubConn) Ping(ctx context.Context) error {
	return s.pingErr
}

// fakeConnector provides a connector backed by stubConn instances so InitDB can
// simulate different ping outcomes.
type fakeConnector struct {
	pingErr error
}

func (c *fakeConnector) Connect(ctx context.Context) (driver.Conn, error) {
	return &stubConn{pingErr: c.pingErr}, nil
}

func (c *fakeConnector) Driver() driver.Driver {
	return &fakeDriver{pingErr: c.pingErr}
}

// fakeDriver satisfies the driver.Driver interface for code paths that call
// sql.Open directly instead of sql.OpenDB.
type fakeDriver struct {
	pingErr error
}

func (d *fakeDriver) Open(name string) (driver.Conn, error) {
	return &stubConn{pingErr: d.pingErr}, nil
}

// TestInitDBSuccess verifies that the database is opened with the expected DSN
// and a reachable connection is stored globally.
func TestInitDBSuccess(t *testing.T) {
	originalOpen := sqlOpen
	var capturedDSN string
	sqlOpen = func(name, dsn string) (*sql.DB, error) {
		capturedDSN = dsn
		return sql.OpenDB(&fakeConnector{}), nil
	}
	t.Cleanup(func() {
		sqlOpen = originalOpen
	})

	if err := InitDB("test", "testdb", "secret", "localhost"); err != nil {
		t.Fatalf("InitDB returned error: %v", err)
	}
	t.Cleanup(func() {
		if ACTIVE_DB != nil {
			ACTIVE_DB.Close()
			ACTIVE_DB = nil
		}
	})
	expected := regexp.MustCompile(`user=test dbname=testdb password=secret host=localhost sslmode=disable`)
	if !expected.MatchString(capturedDSN) {
		t.Fatalf("unexpected DSN: %s", capturedDSN)
	}
}

// TestInitDBOpenError ensures an error from sqlOpen is returned unchanged.
func TestInitDBOpenError(t *testing.T) {
	originalOpen := sqlOpen
	sqlOpen = func(name, dsn string) (*sql.DB, error) {
		return nil, errors.New("boom")
	}
	t.Cleanup(func() {
		sqlOpen = originalOpen
	})

	if err := InitDB("test", "db", "pw", "localhost"); err == nil {
		t.Fatalf("expected error when open fails")
	}
}

// TestInitDBPingError validates that ping failures are surfaced and the global
// DB handle is cleaned up.
func TestInitDBPingError(t *testing.T) {
	originalOpen := sqlOpen
	sqlOpen = func(name, dsn string) (*sql.DB, error) {
		return sql.OpenDB(&fakeConnector{pingErr: errors.New("ping failed")}), nil
	}
	t.Cleanup(func() {
		sqlOpen = originalOpen
	})

	if err := InitDB("user", "db", "pw", "localhost"); err == nil {
		t.Fatalf("expected ping failure to propagate")
	}
	if ACTIVE_DB != nil {
		ACTIVE_DB.Close()
		ACTIVE_DB = nil
	}
}

// TestGetDB confirms the accessor returns the active database pointer.
func TestGetDB(t *testing.T) {
	dbInstance := &sql.DB{}
	ACTIVE_DB = dbInstance
	if GetDB() != dbInstance {
		t.Fatalf("GetDB did not return active instance")
	}
}

// fakeStmt models the behaviour of prepared statements, letting tests inject
// row and execution outcomes.
type fakeStmt struct {
	row     rowScanner
	execErr error
	closed  bool
}

func (f *fakeStmt) QueryRow(args ...any) rowScanner {
	return f.row
}

func (f *fakeStmt) Exec(args ...any) (sql.Result, error) {
	if f.execErr != nil {
		return nil, f.execErr
	}
	return fakeResult{}, nil
}

func (f *fakeStmt) Close() error {
	f.closed = true
	return nil
}

// fakeRow makes it easy to pre-program Scan results for positive and negative
// test paths.
type fakeRow struct {
	values []any
	err    error
}

func (f fakeRow) Scan(dest ...any) error {
	if f.err != nil {
		return f.err
	}
	for i := range dest {
		if i >= len(f.values) {
			return errors.New("not enough values")
		}
		switch d := dest[i].(type) {
		case *string:
			if str, ok := f.values[i].(string); ok {
				*d = str
			}
		default:
			return errors.New("unsupported scan type")
		}
	}
	return nil
}

type fakeResult struct{}

func (fakeResult) LastInsertId() (int64, error) { return 0, nil }
func (fakeResult) RowsAffected() (int64, error) { return 1, nil }

// TestGetUserByName verifies the happy path of retrieving and scanning user
// data from the database.
func TestGetUserByName(t *testing.T) {
	originalPrepare := prepare
	row := fakeRow{values: []any{"alice", "hashed", "Earth", "127.0.0.1"}}
	stmt := &fakeStmt{row: row}
	prepare = func(db *sql.DB, query string) (statement, error) {
		return stmt, nil
	}
	ACTIVE_DB = &sql.DB{}
	t.Cleanup(func() {
		prepare = originalPrepare
	})

	user, err := GetUserByName("alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.User_Name != "alice" || user.Password != "hashed" || user.Location != "Earth" || user.IP_addr != "127.0.0.1" {
		t.Fatalf("unexpected user data: %+v", user)
	}
	if !stmt.closed {
		t.Fatalf("expected statement to be closed")
	}
}

// TestGetUserByNameError asserts that prepare failures are surfaced.
func TestGetUserByNameError(t *testing.T) {
	originalPrepare := prepare
	prepare = func(db *sql.DB, query string) (statement, error) {
		return nil, errors.New("prepare failed")
	}
	ACTIVE_DB = &sql.DB{}
	t.Cleanup(func() {
		prepare = originalPrepare
	})

	if _, err := GetUserByName("alice"); err == nil {
		t.Fatalf("expected prepare error")
	}
}

// TestGetUserByNameScanError ensures scan errors propagate to the caller.
func TestGetUserByNameScanError(t *testing.T) {
	originalPrepare := prepare
	stmt := &fakeStmt{row: fakeRow{err: errors.New("no rows")}}
	prepare = func(db *sql.DB, query string) (statement, error) {
		return stmt, nil
	}
	ACTIVE_DB = &sql.DB{}
	t.Cleanup(func() {
		prepare = originalPrepare
	})

	if _, err := GetUserByName("alice"); err == nil {
		t.Fatalf("expected scan error")
	}
}

// TestRegisterUser covers the happy path of inserting a new user record.
func TestRegisterUser(t *testing.T) {
	originalPrepare := prepare
	stmt := &fakeStmt{}
	prepare = func(db *sql.DB, query string) (statement, error) {
		return stmt, nil
	}
	ACTIVE_DB = &sql.DB{}
	t.Cleanup(func() {
		prepare = originalPrepare
	})

	user := models.ServiceUser{User_Name: "alice", Password: "hashed", Location: "Earth", IP_addr: "127.0.0.1"}
	if err := RegisterUser(user); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !stmt.closed {
		t.Fatalf("expected statement to be closed")
	}
}

// TestRegisterUserPrepareError validates errors returned during statement
// creation are surfaced.
func TestRegisterUserPrepareError(t *testing.T) {
	originalPrepare := prepare
	prepare = func(db *sql.DB, query string) (statement, error) {
		return nil, errors.New("prepare failed")
	}
	ACTIVE_DB = &sql.DB{}
	t.Cleanup(func() {
		prepare = originalPrepare
	})

	if err := RegisterUser(models.ServiceUser{}); err == nil {
		t.Fatalf("expected prepare error")
	}
}

// TestRegisterUserExecError asserts that Exec failures are passed through.
func TestRegisterUserExecError(t *testing.T) {
	originalPrepare := prepare
	stmt := &fakeStmt{execErr: errors.New("insert failed")}
	prepare = func(db *sql.DB, query string) (statement, error) {
		return stmt, nil
	}
	ACTIVE_DB = &sql.DB{}
	t.Cleanup(func() {
		prepare = originalPrepare
	})

	if err := RegisterUser(models.ServiceUser{}); err == nil {
		t.Fatalf("expected exec error")
	}
}

// TestGetSecretKey confirms the secret key is read and converted into bytes.
func TestGetSecretKey(t *testing.T) {
	originalPrepare := prepare
	stmt := &fakeStmt{row: fakeRow{values: []any{"topsecret"}}}
	prepare = func(db *sql.DB, query string) (statement, error) {
		return stmt, nil
	}
	ACTIVE_DB = &sql.DB{}
	t.Cleanup(func() {
		prepare = originalPrepare
	})

	key, err := GetSecretKey()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(key) != "topsecret" {
		t.Fatalf("unexpected secret: %s", string(key))
	}
}

// TestGetSecretKeyErrors covers the error paths for statement creation and row
// scanning when retrieving the secret.
func TestGetSecretKeyErrors(t *testing.T) {
	originalPrepare := prepare
	prepare = func(db *sql.DB, query string) (statement, error) {
		return nil, errors.New("prepare failed")
	}
	ACTIVE_DB = &sql.DB{}
	t.Cleanup(func() {
		prepare = originalPrepare
	})

	if _, err := GetSecretKey(); err == nil {
		t.Fatalf("expected prepare failure")
	}

	rowStmt := &fakeStmt{row: fakeRow{err: errors.New("not found")}}
	prepare = func(db *sql.DB, query string) (statement, error) {
		return rowStmt, nil
	}

	if _, err := GetSecretKey(); err == nil {
		t.Fatalf("expected scan failure")
	}
}
