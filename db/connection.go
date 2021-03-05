package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
)

// Connection contains database instance
type Connection struct {
	DB *sql.DB
}

// NewConnection returns new connection to the db
func NewConnection(driver string, uri string) *Connection {
	db := NewDBConnection(driver, uri)
	return &Connection{DB: db}
}

// NewDBConnection returns a new db connection
func NewDBConnection(driver string, uri string) *sql.DB {
	db, err := sql.Open(driver, uri)
	if err != nil {
		log.Fatalf("Failed to open database: %s, %s, %v", driver, uri, err)
	}

	return db
}

// SetupDatabase creates a database, a user and assigns the created user to the database
func (c *Connection) SetupDatabase(dbName string, username string, password string) error {
	if err := c.CreateDatabase(dbName); err != nil {
		return err
	}

	if err := c.RevokePublicAccess(dbName); err != nil {
		return err
	}

	if err := c.CreateDatabaseUser(username, password); err != nil {
		return err
	}

	return c.AssignUserToDB(dbName, username)
}

// CreateDatabase creates a new database
func (c *Connection) CreateDatabase(name string) error {
	_, err := c.DB.Exec(fmt.Sprintf("CREATE DATABASE %s", name))
	return err
}

// RevokePublicAccess revokes access to the database by anyone
func (c *Connection) RevokePublicAccess(name string) error {
	_, err := c.DB.Exec(fmt.Sprintf("REVOKE connect ON DATABASE %s FROM PUBLIC;", name))
	return err
}

// CreateDatabaseUser creates a new database user
func (c *Connection) CreateDatabaseUser(name string, password string) error {
	_, err := c.DB.Exec(fmt.Sprintf("CREATE USER %s WITH ENCRYPTED PASSWORD '%s'", name, password))
	return err
}

// AssignUserToDB grants user access to the database
func (c *Connection) AssignUserToDB(dbName string, username string) error {
	_, err := c.DB.Exec(fmt.Sprintf("GRANT ALL PRIVILEGES ON DATABASE %s TO %s", dbName, username))
	return err
}

// RunMigrations runs all remaining db migrations
func (c *Connection) RunMigrations(migrationsPath string, driver string) error {
	driverInstance, err := postgres.WithInstance(c.DB, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(migrationsPath, driver, driverInstance)
	if err != nil {
		return fmt.Errorf("Cannot load migrations: %v", err)
	}

	if err := m.Up(); err != nil {
		return fmt.Errorf("Migration run error: %v", err)
	}

	return nil
}
