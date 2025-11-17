package config

import (
	"log"
	"os"
)

type Config struct {
	MongoURI      string
	MongoDatabase string
	ServerPort    string
}

// Cargar congiguración desde variables de entorno
func Load() *Config {
	mongoUri := os.Getenv("MONGODB_URI")
	if mongoUri == "" {
		log.Fatal("MONGODB_URI no está configurado")
	}

	mongoDataBase := os.Getenv("MONGODB_DATABASE")
	if mongoDataBase == "" {
		mongoDataBase = "taskProcessor"
	}

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	return &Config{
		MongoURI:      mongoUri,
		MongoDatabase: mongoDataBase,
		ServerPort:    serverPort,
	}

}
