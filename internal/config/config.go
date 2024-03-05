package config

import (
	"flag"
	"os"
)

type Config struct {
	ServerAddress      string
	AllowTypes         []string
	DBConnect          string
	DecoderWorkerCount int
}

func NewConfig() *Config {
	serverAddress := ":8080"
	dbConnect := ""

	flag.StringVar(&serverAddress, "a", ":8080", "address and port to run server")
	flag.StringVar(&dbConnect, "db", "", "DB Connect address")
	flag.Parse()

	if envServerAddress := os.Getenv("ADDRESS"); envServerAddress != "" {
		serverAddress = envServerAddress
	}

	if envDBConnect := os.Getenv("DB_CONNECT"); envDBConnect != "" {
		dbConnect = envDBConnect
	}

	return &Config{
		ServerAddress:      serverAddress,
		DBConnect:          dbConnect,
		AllowTypes:         []string{"video/mp4", "audio/mpeg", "image/png"},
		DecoderWorkerCount: 2,
	}
}
