package postgres

import (
	"database/sql"
	"errors"

	goconfig "github.com/easeq/go-config"
	"github.com/easeq/go-service/component"
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
	i      component.Initializer
	logger logger.Logger
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

	pg := &Postgres{
		Handle: newConnection(cfg.GetURI()),
		Config: cfg,
	}

	pg.i = NewInitializer(pg)
	return pg
}

// GetConfig returns the DB config
func GetConfig() *Config {
	return goconfig.NewEnvConfig(new(Config)).(*Config)
}

// Migrate runs all remaining db migrations
func (db *Postgres) Migrate() error {
	instance, err := db.instance()
	if err != nil {
		db.logger.Fatalw(
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
		db.logger.Fatalw(
			ErrMigrationLoad.Error(),
			"error", err,
		)
		return ErrMigrationLoad
	}

	if err := m.Up(); err != nil {
		db.logger.Debugw(
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

func (db *Postgres) HasInitializer() bool {
	return true
}

func (db *Postgres) Initializer() component.Initializer {
	return db.i
}
