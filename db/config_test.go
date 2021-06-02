package db

import (
	"os"
	"testing"

	env "github.com/Netflix/go-env"
)

func TestUnmarshalEnv(t *testing.T) {
	tests := []struct {
		name        string
		vars        map[string]string
		emptyConfig Config
		want        Config
	}{
		{
			name:        "checkDefaultValues",
			vars:        map[string]string{},
			emptyConfig: Config{},
			want: Config{
				Driver:  "postgres",
				Host:    "localhost",
				Port:    5432,
				SSLMode: "disable",
			},
		},
		{
			name: "loadAllVarsFromEnv",
			vars: map[string]string{
				"DB_NAME":            "test-db-name",
				"DB_USER":            "test-db-username",
				"DB_PASS":            "test-db-password",
				"DB_DRIVER":          "mysql",
				"DB_HOST":            "test-db-host",
				"DB_PORT":            "1234",
				"DB_SSLMODE":         "disable",
				"DB_MIGRATIONS_PATH": "/test/path",
			},
			emptyConfig: Config{},
			want: Config{
				Name:           "test-db-name",
				User:           "test-db-username",
				Password:       "test-db-password",
				Driver:         "mysql",
				Host:           "test-db-host",
				Port:           1234,
				SSLMode:        "disable",
				MigrationsPath: "/test/path",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()

			for k, v := range tt.vars {
				os.Setenv(k, v)
			}

			envSet, err := env.EnvironToEnvSet(os.Environ())
			if err != nil {
				t.Errorf("Error loading EnvSet for %s", tt.name)
			}

			tt.emptyConfig.UnmarshalEnv(envSet)
			if tt.emptyConfig != tt.want {
				t.Errorf("Unmarshalenv %s failed, got(%v) want(%v)", tt.name, tt.emptyConfig, tt.want)
			}
		})
	}
}

func TestGetURI(t *testing.T) {
	tests := []struct {
		name   string
		vars   map[string]string
		config Config
		want   string
	}{
		{
			name: "uri",
			config: Config{
				Name:           "test-db-name",
				User:           "test-db-user",
				Password:       "test-db-password",
				Driver:         "mysql",
				Host:           "test-db-host",
				Port:           1234,
				SSLMode:        "disable",
				MigrationsPath: "/test/path",
			},
			want: "mysql://test-db-user:test-db-password@test-db-host:1234/test-db-name?sslmode=disable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.config.GetURI()
			if got != tt.want {
				t.Errorf("%s failed got(%s) want(%s)", tt.name, got, tt.want)
			}
		})
	}
}
