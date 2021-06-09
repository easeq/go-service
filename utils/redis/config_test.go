package redis

import (
	"os"
	"reflect"
	"testing"
	"time"

	env "github.com/Netflix/go-env"
)

func TimeDuration(format string) time.Duration {
	v, _ := time.ParseDuration(format)
	return v
}

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name string
		vars map[string]string
		want Config
	}{
		{
			name: "defaultConfig",
			vars: map[string]string{},
			want: Config{
				Network:            "tcp",
				Addr:               "localhost:6379",
				Username:           "",
				Password:           "",
				DB:                 0,
				MaxRetries:         3,
				MinRetryBackoff:    TimeDuration("8ms"),
				MaxRetryBackoff:    TimeDuration("512ms"),
				DialTimeout:        TimeDuration("5s"),
				ReadTimeout:        TimeDuration("3s"),
				WriteTimeout:       TimeDuration(""),
				PoolSize:           10,
				MinIdleConns:       0,
				MaxConnAge:         0,
				PoolTimeout:        0,
				IdleTimeout:        0,
				IdleCheckFrequency: TimeDuration("60s"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()

			for k, v := range tt.vars {
				os.Setenv(k, v)
			}

			config := NewConfig()

			if !reflect.DeepEqual(*config, tt.want) {
				t.Errorf("Unmarshalenv %s failed, got(%v) want(%v)", tt.name, config, tt.want)
			}
		})
	}
}

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
				Network:            "tcp",
				Addr:               "localhost:6379",
				Username:           "",
				Password:           "",
				DB:                 0,
				MaxRetries:         3,
				MinRetryBackoff:    TimeDuration("8ms"),
				MaxRetryBackoff:    TimeDuration("512ms"),
				DialTimeout:        TimeDuration("5s"),
				ReadTimeout:        TimeDuration("3s"),
				WriteTimeout:       TimeDuration(""),
				PoolSize:           10,
				MinIdleConns:       0,
				MaxConnAge:         0,
				PoolTimeout:        0,
				IdleTimeout:        0,
				IdleCheckFrequency: TimeDuration("60s"),
			},
		},
		{
			name: "loadAllVarsFromEnv",
			vars: map[string]string{
				"REDIS_ADDRESS":      "test-host:1234",
				"REDIS_USERNAME":     "test-username",
				"REDIS_MAX_RETRIES":  "5",
				"REDIS_DIAL_TIMEOUT": "10s",
			},
			emptyConfig: Config{},
			want: Config{
				Network:            "tcp",
				Addr:               "test-host:1234",
				Username:           "test-username",
				Password:           "",
				DB:                 0,
				MaxRetries:         5,
				MinRetryBackoff:    TimeDuration("8ms"),
				MaxRetryBackoff:    TimeDuration("512ms"),
				DialTimeout:        TimeDuration("10s"),
				ReadTimeout:        TimeDuration("3s"),
				WriteTimeout:       TimeDuration(""),
				PoolSize:           10,
				MinIdleConns:       0,
				MaxConnAge:         0,
				PoolTimeout:        0,
				IdleTimeout:        0,
				IdleCheckFrequency: TimeDuration("60s"),
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

			if err := tt.emptyConfig.UnmarshalEnv(envSet); err != nil {
				t.Errorf("Error unmarshaling env")
			}

			if !reflect.DeepEqual(tt.emptyConfig, tt.want) {
				t.Errorf("Unmarshalenv %s failed, got(%v) want(%v)", tt.name, tt.emptyConfig, tt.want)
			}
		})
	}
}
