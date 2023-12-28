package config

type EnvConfig struct {
	HttpAddress string `env:"HTTP_ADDRESS" default:"0.0.0.0:33333"`
}
