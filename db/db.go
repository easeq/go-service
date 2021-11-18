package db

import "github.com/easeq/go-service/component"

const (
	DATABASE = "database"
)

// ServiceDatabase interface for database
type ServiceDatabase interface {
	component.Component
	// Migrate database schema
	Migrate() error
}
