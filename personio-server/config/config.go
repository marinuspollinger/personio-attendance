package config

import "time"

type EnvConfig struct {
	LogLevel                string        `env:"LOG_LEVEL" default:"info"`
	HttpAddress             string        `env:"HTTP_ADDRESS" default:"0.0.0.0:33333"`
	CurrentTimeLoopInterval time.Duration `env:"CURRENT_TIME_LOOP_INTERVAL" default:"30s"`

	PersonioEmployeeId   int    `env:"PERSONIO_EMPLOYEE_ID" required:"true"`
	PersonioHost         string `env:"PERSONIO_HOST" required:"true"`
	PersonioClientId     string `env:"PERSONIO_CLIENT_ID" required:"true"`
	PersonioClientSecret string `env:"PERSONIO_CLIENT_SECRET" required:"true"`
}
