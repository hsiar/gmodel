package rds

type Config struct {
	Host     string
	Port     int
	Password string
	Db       int
}

func DefaultConfig() *Config {
	return &Config{
		Host:     "127.0.0.1",
		Port:     6379,
		Password: "",
		Db:       0,
	}
}
