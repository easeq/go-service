package db

import (
	"database/sql"
	"errors"
	"fmt"

	goconfig "github.com/easeq/go-config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	// migration file driver
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var (
	// ErrDBOpen returned when db open fails
	ErrDBOpen = errors.New("error opening database")
	// ErrDBConfigLoad returned when env config for DB results in an error
	ErrDBConfigLoad = errors.New("error loading database config")
	// ErrMigrationLoad returned when migrations cannot be loaded
	ErrMigrationLoad = errors.New("error loading new migrations")
	// ErrRunMigration returned when migrations run is unsuccessful
	ErrRunMigration = errors.New("error while running migrations")
	// ErrCreateDBInstance returned when db instance creation fails
	ErrCreateDBInstance = errors.New("error while creating a new DB instance")
)

const (
	dbStmtCreate       = "CREATE DATABASE %s"
	dbStmtCreateUser   = "CREATE USER %s WITH ENCRYPTED PASSWORD '%s'"
	dbStmtRevokeAccess = "REVOKE connect ON DATABASE %s FROM PUBLIC;"
	dbStmtAssignUser   = "GRANT ALL PRIVILEGES ON DATABASE %s TO %s"
	// Driver defines the driver name
	Driver = "postgres"
)

// Postgres contains database instance
type Postgres struct {
	Handle *sql.DB
	IsRoot bool
	*Config
}

func newConnection(uri string) *sql.DB {
	db, err := sql.Open(Driver, uri)
	if err != nil {
		panic(ErrDBOpen)
	}

	return db
}

// NewPostgres returns new connection to the postgres db
func NewPostgres() ServiceDatabase {
	cfg := GetConfig()

	return &Postgres{
		Handle: newConnection(cfg.GetURI(true)),
		Config: cfg,
	}
}

// GetConfig returns the DB config
func GetConfig() *Config {
	return goconfig.NewEnvConfig(new(Config)).(*Config)
}

// Setup creates a database, a user and assigns the created user to the database
func (db *Postgres) Setup() *ErrDatabaseSetup {
	if err := db.Create(); err != nil {
		return &ErrDatabaseSetup{err}
	}

	if err := db.RevokePublicAccess(); err != nil {
		return &ErrDatabaseSetup{err}
	}

	if err := db.CreateUser(); err != nil {
		return &ErrDatabaseSetup{err}
	}

	if err := db.AssignUser(); err != nil {
		return &ErrDatabaseSetup{err}
	}

	return nil
}

// Create creates a new database
func (db *Postgres) Create() error {
	_, err := db.Handle.Exec(fmt.Sprintf(
		dbStmtCreate,
		db.Config.Name,
	))
	return err
}

// RevokePublicAccess revokes access to the database by anyone
func (db *Postgres) RevokePublicAccess() error {
	_, err := db.Handle.Exec(fmt.Sprintf(
		dbStmtRevokeAccess,
		db.Name,
	))
	return err
}

// CreateUser creates a new database user
func (db *Postgres) CreateUser() error {
	_, err := db.Handle.Exec(fmt.Sprintf(
		dbStmtCreateUser,
		db.User,
		db.Password,
	))
	return err
}

// AssignUser grants user access to the database
func (db *Postgres) AssignUser() error {
	_, err := db.Handle.Exec(fmt.Sprintf(
		dbStmtAssignUser,
		db.Name,
		db.User,
	))
	return err
}

// Migrate runs all remaining db migrations
func (db *Postgres) Migrate() error {
	m, err := migrate.NewWithDatabaseInstance(
		db.MigrationsPath,
		db.Driver,
		db.instance(),
	)

	if err != nil {
		return ErrMigrationLoad
	}

	if err := m.Up(); err != nil {
		return ErrRunMigration
	}

	return nil
}

func (db *Postgres) instance() database.Driver {
	driverInstance, err := postgres.WithInstance(db.Handle, &postgres.Config{})
	if err != nil {
		panic(ErrCreateDBInstance)
	}

	return driverInstance
}

// Close database connection
func (db *Postgres) Close() error {
	return db.Handle.Close()
}

// UpdateHandle updates the database handle
func (db *Postgres) UpdateHandle() error {
	if err := db.Close(); err != nil {
		return err
	}

	db.Handle = newConnection(db.Config.GetURI(false))
	return nil
}
