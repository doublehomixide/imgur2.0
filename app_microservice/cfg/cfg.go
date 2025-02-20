package cfg

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	MinioURL      string
	MinioUSER     string
	MinioPASSWORD string
	PsqlDBPath    string
	ServerPort    string
}

func Init() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	minioURL := os.Getenv("minioURL")
	minioUSER := os.Getenv("minioUSER")
	minioPASSWORD := os.Getenv("minioPASSWORD")
	psqlDBPath := os.Getenv("DATABASE_URL")
	serverPort := os.Getenv("PORT")
	return &Config{minioURL, minioUSER, minioPASSWORD, psqlDBPath, serverPort}
}
