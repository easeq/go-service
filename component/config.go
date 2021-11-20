package component

import (
	"github.com/Netflix/go-env"
)

// NewConfig creates a new instance of component Config
func NewConfig(v interface{}) error {
	_, err := env.UnmarshalFromEnviron(v)
	if err != nil {
		return err
	}

	return nil
}
