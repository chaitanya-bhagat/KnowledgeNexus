package config

const (
	APP_HOST = "APP_HOST"
	APP_PORT = "APP_PORT"
)

type ServerConfig struct {
	Host string
	Port string
}

func loadServerConfig() (ServerConfig, error) {

	return ServerConfig{
		Host: envOrDefault(APP_HOST, "0.0.0.0"),
		Port: envOrDefault(APP_PORT, "8080"),
	}, nil
}

func (sc ServerConfig) Address() string {
	return sc.Host + ":" + sc.Port
}
