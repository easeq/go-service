package gateway

import (
	"os"
	"testing"

	env "github.com/Netflix/go-env"
	"github.com/stretchr/testify/assert"
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
			want:        Config{Port: 8080, Metadata: Metadata{GrpcPort: 9090}},
		},
		{
			name: "loadAllVarsFromEnv",
			vars: map[string]string{
				"HTTP_HOST":        "localhost",
				"HTTP_PORT":        "8085",
				"GRPC_HOST":        "grpc-host",
				"GRPC_PORT":        "9095",
				"HTTP_CONSUL_TAGS": "primary,secondary",
			},
			emptyConfig: Config{},
			want:        Config{"localhost", 8085, "primary,secondary", Metadata{"grpc-host", 9095}},
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

func TestAddress(t *testing.T) {
	tests := []struct {
		name   string
		vars   map[string]string
		config Config
		want   string
	}{
		{
			name:   "getAddressWithHostAndPort",
			config: Config{Host: "localhost", Port: 8085},
			want:   "localhost:8085",
		},
		{
			name:   "getAddressWithOnlyPort",
			config: Config{Port: 8085},
			want:   ":8085",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.config.Address()
			if got != tt.want {
				t.Errorf("%s failed got(%s) want(%s)", tt.name, got, tt.want)
			}
		})
	}
}

func TestGetTags(t *testing.T) {
	tests := []struct {
		name   string
		vars   map[string]string
		config Config
		want   []string
	}{
		{
			name: "getConsulTagsFromEnv",
			vars: map[string]string{
				"HTTP_CONSUL_TAGS": "primary,secondary",
			},
			config: Config{},
			want:   []string{"primary", "secondary"},
		}, {
			name:   "getConsulTagsSetManually",
			config: Config{Tags: "tag-3,tag-4"},
			want:   []string{"tag-3", "tag-4"},
		}, {
			name:   "getEmptyConsulTags",
			config: Config{},
			want:   []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()

			if tt.vars != nil {
				for k, v := range tt.vars {
					os.Setenv(k, v)
				}

				envSet, err := env.EnvironToEnvSet(os.Environ())
				if err != nil {
					t.Errorf("Error loading EnvSet for %s", tt.name)
				}

				tt.config.UnmarshalEnv(envSet)
			}

			got := tt.config.GetTags()
			if !assert.Equal(t, got, tt.want) {
				t.Errorf("%s failed got(%s) want(%s)", tt.name, got, tt.want)
			}
		})
	}
}