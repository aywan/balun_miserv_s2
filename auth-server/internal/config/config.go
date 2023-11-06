package config

import (
	"github.com/aywan/balun_miserv_s2/shared/lib/db"
	"github.com/aywan/balun_miserv_s2/shared/lib/logger"
)

// Config full application configuration.
type Config struct {
	Mode   logger.Mode `envconfig:"mode"`
	Server Server      `envconfig:"server"`
	Db     db.Config   `envconfig:"db"`
}
