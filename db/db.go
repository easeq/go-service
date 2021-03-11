package db

// ServiceDatabase interface for database
type ServiceDatabase interface {
	// Setup database
	Setup() error
	// Migrate database schema
	Migrate() error
	// Close database connection
	Close() error
}
