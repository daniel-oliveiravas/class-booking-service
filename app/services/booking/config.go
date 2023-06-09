package main

import (
	"github.com/kelseyhightower/envconfig"
)

const configPrefix = "MEMBERS"

type Config struct {
	Host    string `split_words:"true" default:":8080" desc:"group the service belongs to"`
	GinMode string `split_words:"true" default:"release" desc:"group the service belongs to"`

	PostgresHostname       string `split_words:"true" default:"localhost" desc:"postgres hostname"`
	PostgresDatabaseName   string `split_words:"true" default:"class_booking" desc:"postgres database name to connect to"`
	PostgresDatabaseNameQA string `split_words:"true" default:"class_booking_qa" desc:"postgres database name to connect to"`
	PostgresUser           string `split_words:"true" default:"class_booking" desc:"postgres user to connect as"`
	PostgresPassword       string `split_words:"true" default:"class_booking" desc:"postgres password"`
	PostgresPort           int    `split_words:"true" default:"5432" desc:"postgres port number"`
	PostgresSSLMode        string `split_words:"true" default:"none" desc:"postgres connection ssl mode"`
}

func loadConfig() (Config, error) {
	var cfg Config
	err := envconfig.Process(configPrefix, &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
