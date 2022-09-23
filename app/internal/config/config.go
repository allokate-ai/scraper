package config

import (
	"fmt"

	"github.com/joho/godotenv"

	"github.com/allokate-ai/environment"
)

type ConnectionInfo struct {
	Host string
	Port int
}

type Config struct {
	Port    int
	FluentD ConnectionInfo
}

func Get() (Config, error) {
	godotenv.Load()

	port, err := environment.GetInt("PORT")
	if err != nil {
		fmt.Println(err)
		return Config{}, fmt.Errorf("invalid value (%s) for port", environment.GetValue("PORT"))
	}

	fluentdPort, err := environment.GetInt("FLUENTD_PORT")
	if err != nil {
		return Config{}, fmt.Errorf("invalid value (%s) for fluentd port", environment.GetValue("FLUENTD_PORT"))
	}

	return Config{
		Port: int(port),
		FluentD: ConnectionInfo{
			Host: environment.GetValue("FLUENTD_HOST"),
			Port: int(fluentdPort),
		},
	}, nil
}
