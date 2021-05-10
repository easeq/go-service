package db

import "fmt"

// ErrDatabaseSetup returned when database setup fails
type ErrDatabaseSetup struct {
	value error
}

func (e *ErrDatabaseSetup) Error() string {
	return fmt.Sprintf("database setup failed: [%s]", e.value)
}

// ServiceDatabase interface for database
type ServiceDatabase interface {
	Init() error
	// Setup database
	Setup() *ErrDatabaseSetup
	// Migrate database schema
	Migrate() error
	// Update database handle
	UpdateHandle() error
	// Close database connection
	Close() error
}
