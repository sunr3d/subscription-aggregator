package config

import "time"

type Config struct {
  HTTPPort    string        `envconfig:"HTTP_PORT" default:"8080"`
  HTTPTimeout time.Duration `envconfig:"HTTP_TIMEOUT" default:"30s"`
  LogLevel    string        `envconfig:"LOG_LEVEL" default:"info"`
  Postgres    PostgresConfig `envconfig:"POSTGRES"`
}

type PostgresConfig struct {
	Host     string `envconfig:"HOST" default:"localhost"`
	Port     string `envconfig:"PORT" default:"5432"`
	User     string `envconfig:"USER" default:"postgres"`
	Password string `envconfig:"PASSWORD" default:"postgres"`
	DBName   string `envconfig:"DB_NAME" default:"postgres"`
	SSLMode  string `envconfig:"SSL_MODE" default:"disable"`
}