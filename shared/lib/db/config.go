package db

// Config configuration for postgresql connection.
type Config struct {
	Host     string
	Port     string `default:"5432"`
	Name     string
	User     string
	Pass     string
	Trace    bool   `default:"false"`
	LogLevel string `default:"info"`
}
