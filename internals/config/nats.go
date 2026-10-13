package config

const (
	NATS_URL = "NATS_URL"
)

type NatsConfig struct {
	URL string
}

func loadNatsConfig() (NatsConfig, error) {
	url, err := requiredEnv(NATS_URL)
	if err != nil {
		return NatsConfig{}, err
	}
	return NatsConfig{
		URL: url,
	}, nil
}
