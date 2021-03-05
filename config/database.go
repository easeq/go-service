package config

import (
	"fmt"
)

// Database configuration
type Database struct {
	Name         string `env:"DB_NAME"`
	User         string `env:"DB_USER"`
	Password     string `env:"DB_PASS"`
	Driver       string `env:"DB_DRIVER,default=postgres"`
	Host         string `env:"DB_HOST,default=localhost"`
	Port         int    `env:"DB_PORT,default=5432"`
	SSLMode      string `env:"DB_SSL_MODE,default=disable"`
	RootUser     string `env:"DB_ROOT_USER"`
	RootPassword string `env:"DB_ROOT_PASS"`
}

// GetUsername returns the required db username
func (db *Database) GetUsername(root bool) string {
	if root == true {
		return db.RootUser
	}

	return db.User
}

// GetPassword returns the required db password
func (db *Database) GetPassword(root bool) string {
	if root == true {
		return db.RootPassword
	}

	return db.Password
}

// GetURI generates and returns the database URI from the provided config
func (db *Database) GetURI(includeDbName bool, asRoot bool) string {
	dbName := db.Name
	if includeDbName != true {
		dbName = ""
	}

	return fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=%s",
		db.Driver,
		db.GetUsername(asRoot),
		db.GetPassword(asRoot),
		db.Host,
		db.Port,
		dbName,
		db.SSLMode,
	)
}
