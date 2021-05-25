package db

// ServiceDatabase interface for database
type ServiceDatabase interface {
	Init() error
	// Migrate database schema
	Migrate() error
	// Close database connection
	Close() error
}
