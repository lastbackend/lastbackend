package middleware

import (
	"github.com/lastbackend/lastbackend/internal/pkg/storage"
	"github.com/spf13/viper"
)

const (
	logPrefix = "http:middleware"
)

type Middleware struct {
	storage storage.Storage
	viper   *viper.Viper
}

func New(stg storage.Storage, v *viper.Viper) Middleware {
	return Middleware{
		storage: stg,
		viper:   v,
	}
}
