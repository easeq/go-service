package object_store

import (
	"errors"

	"github.com/easeq/go-service/component"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	// ErrCreatingMinioClient returned when instantiating a minio client fails
	ErrCreatingMinioClient = errors.New("error creating object-store minio client")
)

type Minio struct {
	Client *minio.Client
	Config *Config
}

// NewMinio returns a new instance of minio client and config
func NewMinio() *Minio {
	config := NewConfig()

	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		panic(ErrCreatingMinioClient)
	}

	return &Minio{Client: client, Config: config}
}

func (m *Minio) HasInitializer() bool {
	return false
}

func (m *Minio) Initializer() component.Initializer {
	return nil
}
