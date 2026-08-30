package config

import (
	"fmt"
	"net/url"
)

const (
	POSTGRES_USER     = "POSTGRES_USER"
	POSTGRES_PASSWORD = "POSTGRES_PASSWORD"
	POSTGRES_DB       = "POSTGRES_DB"
	POSTGRES_HOST     = "POSTGRES_HOST"
	POSTGRES_PORT     = "POSTGRES_PORT"
)

type PostgresConfig struct {
	User     string
	Password string
	DB       string
	Host     string
	Port     string
}

func loadPostgresConfig() (PostgresConfig, error) {
	user, err := requiredEnv(POSTGRES_USER)
	if err != nil {
		return PostgresConfig{}, err
	}
	passwd, err := requiredEnv(POSTGRES_PASSWORD)
	if err != nil {
		return PostgresConfig{}, err
	}
	db, err := requiredEnv(POSTGRES_DB)
	if err != nil {
		return PostgresConfig{}, err
	}
	return PostgresConfig{
		User:     user,
		Password: passwd,
		DB:       db,
		Host:     envOrDefault(POSTGRES_HOST, "localhost"),
		Port:     envOrDefault(POSTGRES_PORT, "5432"),
	}, nil
}

func (pc PostgresConfig) DNS() string {
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(pc.User, pc.Password),
		Host:   fmt.Sprintf("%s:%s", pc.Host, pc.Port),
		Path:   fmt.Sprintf("/%s", pc.DB),
	}
	q := u.Query()
	q.Set("sslmode", "disable")
	u.RawQuery = q.Encode()
	return u.String()
}
