package db

import (
	"database/sql"
	"errors"
	"log"

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
	// ErrMigrationLoad returned when migrations cannot be loaded
	ErrMigrationLoad = errors.New("error loading new migrations")
	// ErrCreateDBInstance returned when db instance creation fails
	ErrCreateDBInstance = errors.New("error while creating a new DB instance")
)

const (
	// Driver defines the driver name
	Driver = "postgres"
)

// Postgres contains database instance
type Postgres struct {
	Handle *sql.DB
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
func NewPostgres() *Postgres {
	cfg := GetConfig()

	return &Postgres{
		Handle: newConnection(cfg.GetURI()),
		Config: cfg,
	}
}

// GetConfig returns the DB config
func GetConfig() *Config {
	return goconfig.NewEnvConfig(new(Config)).(*Config)
}

// Init database
func (db *Postgres) Init() error {
	// Run migrations
	if err := db.Migrate(); err != nil {
		log.Println(err)
	}

	return nil
}

// Migrate runs all remaining db migrations
func (db *Postgres) Migrate() error {
	instance, err := db.instance()
	if err != nil {
		return ErrCreateDBInstance
	}

	m, err := migrate.NewWithDatabaseInstance(
		db.MigrationsPath,
		db.Driver,
		instance,
	)

	if err != nil {
		return ErrMigrationLoad
	}

	if err := m.Up(); err != nil {
		return err
	}

	return nil
}

func (db *Postgres) instance() (database.Driver, error) {
	driverInstance, err := postgres.WithInstance(db.Handle, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	return driverInstance, nil
}

// Close database connection
func (db *Postgres) Close() error {
	return db.Handle.Close()
}
