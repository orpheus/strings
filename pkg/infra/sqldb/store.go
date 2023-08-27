package sqldb

import (
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"path/filepath"
)

// DriverPostgres imported via github.com/jackc/pgx/v5/stdlib
const DriverPostgres = "pgx"

var (
	ErrTransactionInProgress = errors.New("transaction already in progress")
	ErrTransactionNotStarted = errors.New("db operation marked as atomic but not run within a transaction")

	PostgresMigrationsDirAbs = filepath.Join("internal", "infra", "sqldb", "migrations", "postgres")
)

type ConnConfig struct {
	Host   string
	Port   string
	User   string
	Pass   string
	Dbname string
}

type Store struct {
	Db *sqlx.DB
	Tx *sqlx.Tx
}

// GetAtomicExecutor checks if the required behavior is Atomic (transaction-based)
// and if it is, and no transaction exist, returns an error.
// If a transaction is in place, it will return the *sql.Tx object
// instead of the *sql.DB object as a QueryAble interface which Repositories
// can use a sql executor.
//
// TL;DR: Determines the sql object to use for queries.
func (s *Store) GetAtomicExecutor(isAtomic bool) (QueryAble, error) {
	if isAtomic && s.Tx == nil {
		return nil, ErrTransactionNotStarted
	}

	var x QueryAble
	x = s.Db
	if s.Tx != nil {
		x = s.Tx
	}

	return x, nil
}

func (s *Store) GetExecutor() QueryAble {
	var x QueryAble
	x = s.Db
	if s.Tx != nil {
		x = s.Tx
	}

	return x
}

func (s *Store) Close() error {
	return s.Db.Close()
}

// NewStore opens a new connection to a database using the given driver
// and a connection url called a Data Source Name (DSN).
//
// sqlx: http://jmoiron.github.io/sqlx/#gettingStarted
func NewStore(driverName string, config ConnConfig) (*Store, error) {
	fmt.Printf("Connecting to %s:<password_hidden>@%s:%s/%s\n", config.User, config.Host, config.Port, config.Dbname)

	dsn := createPostgresDsn(config)
	db, err := sqlx.Open(driverName, dsn)
	if err != nil {
		return nil, err
	}
	return &Store{db, nil}, db.Ping()
}

func createPostgresDsn(config ConnConfig) string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", config.User, config.Pass, config.Host, config.Port, config.Dbname)
}

// NewMockStore based off: https://github.com/jmoiron/sqlx/issues/204
func NewMockStore() (*Store, sqlmock.Sqlmock) {
	mockDB, mock, _ := sqlmock.New()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	return &Store{
		Db: sqlxDB,
		Tx: nil,
	}, mock
}
