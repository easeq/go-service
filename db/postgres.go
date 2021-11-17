package db

import (
	"context"
	"database/sql"
	"errors"

	goconfig "github.com/easeq/go-config"
	"github.com/easeq/go-service/logger"
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
	// ErrDBMigrationFailed returned when db migration fails
	ErrDBMigrationFailed = errors.New("db migration failed")
)

const (
	// Driver defines the driver name
	Driver = "postgres"
)

// Postgres contains database instance
type Postgres struct {
	Handle *sql.DB
	logger logger.Logger
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
		db.logger.Debugf("%s: %s", ErrDBMigrationFailed, err)
	}

	return nil
}

// Migrate runs all remaining db migrations
func (db *Postgres) Migrate() error {
	instance, err := db.instance()
	if err != nil {
		db.logger.Errorw(
			ErrCreateDBInstance.Error(),
			"error", err,
		)
		return ErrCreateDBInstance
	}

	m, err := migrate.NewWithDatabaseInstance(
		db.MigrationsPath,
		db.Driver,
		instance,
	)

	if err != nil {
		db.logger.Errorw(
			ErrMigrationLoad.Error(),
			"error", err,
		)
		return ErrMigrationLoad
	}

	if err := m.Up(); err != nil {
		db.logger.Errorw(
			"Migration UP failed",
			"error", err,
		)
		return err
	}

	return nil
}

func (db *Postgres) instance() (database.Driver, error) {
	driverInstance, err := postgres.WithInstance(db.Handle, &postgres.Config{})
	if err != nil {
		db.logger.Errorw(
			"Postgres driverInstance creation failed",
			"error", err,
		)
		return nil, err
	}

	return driverInstance, nil
}

// Close database connection
func (db *Postgres) Close() error {
	db.logger.Infow(
		"Closing database connection",
	)
	return db.Handle.Close()
}

// AddDependency adds necessary service components as dependencies
func (db *Postgres) AddDependency(dep interface{}) error {
	switch v := dep.(type) {
	case logger.Logger:
		db.logger = v
	}

	return nil
}

// Dependencies returns the string names of service components
// that are required as dependencies for this component
func (db *Postgres) Dependencies() []string {
	return []string{"logger"}
}

// CanRun returns true if the component has anything to Run
func (db *Postgres) CanRun() bool {
	return true
}

// Run start the service component
func (db *Postgres) Run(ctx context.Context) error {
	// Run migrations
	if err := db.Migrate(); err != nil {
		db.logger.Errorw(
			"Database migration failed",
			"error", err,
		)
	}

	return nil
}
