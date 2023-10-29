package config

// Mode - represents the application mode.
type Mode string

const (
	// DevMode - means application in development mode.
	DevMode = Mode("dev")
	// ProdMode - means application in production mode.
	ProdMode = Mode("prod")
)

// Config full application configuration.
type Config struct {
	Mode   Mode     `envconfig:"mode"`
	Server Server   `envconfig:"server"`
	Db     DbConfig `envconfig:"db"`
}
