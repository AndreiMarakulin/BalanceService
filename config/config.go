package config

type Config struct {
	HTTP
	PG
}

type HTTP struct {
	Port string
}

type PG struct {
	URL string
}

func New() *Config {
	return &Config{
		HTTP{
			Port: "8080",
		},
		PG{
			URL: "postgres://user:passw0rd@localhost:5432/balance",
		},
	}
}
