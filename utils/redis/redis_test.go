package redis

import (
	"os"
	"testing"

	"github.com/Netflix/go-env"
	goredis "github.com/go-redis/redis/v8"
	"github.com/go-test/deep"
)

func TestNewRedisClient(t *testing.T) {
	os.Setenv("REDIS_ADDRESS", "test-host:1234")

	tests := []struct {
		name string
		vars map[string]string
	}{
		{
			name: "defaultEnvValues",
			vars: nil,
		},
		{
			name: "withAddress",
			vars: map[string]string{
				"REDIS_ADDRESS": "test-host:1234",
			},
		},
		{
			name: "withPassword",
			vars: map[string]string{
				"REDIS_PASSWORD": "test-pass",
			},
		},
		{
			name: "withDB",
			vars: map[string]string{
				"REDIS_DB": "1",
			},
		},
		{
			name: "withAll",
			vars: map[string]string{
				"REDIS_ADDRESS":  "test-host:1234",
				"REDIS_PASSWORD": "test-pass",
				"REDIS_DB":       "1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()

			for k, v := range tt.vars {
				os.Setenv(k, v)
			}

			config := new(Config)
			envSet, err := env.EnvironToEnvSet(os.Environ())
			if err != nil {
				t.Errorf("Error loading EnvSet for %s", tt.name)
			}

			if err := config.UnmarshalEnv(envSet); err != nil {
				t.Errorf("Error unmarshaling env")
			}

			got := NewRedisClient(config)
			want := &Redis{
				Config: config,
				Client: goredis.NewClient(&goredis.Options{
					Addr:     config.Addr,
					Password: config.Password,
					DB:       config.DB,
				}),
			}

			if diff := deep.Equal(got, want); diff != nil {
				t.Errorf("Invalid redis client want(%v) got(%v) difference(%v)", *got, *want, diff)
			}
		})
	}
}
