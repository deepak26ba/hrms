package config

import (
	"fmt"
	"hrms/pkg/models"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		return
	}
}

func Config() (string, error) {

	var con models.ConnectionString

	con.Host = os.Getenv("DB_HOST")
	con.User = os.Getenv("DB_USER")
	con.DBName = os.Getenv("DB_NAME")
	con.Password = os.Getenv("DB_PASSWORD")
	con.SslMode = os.Getenv("DB_SSLMODE")
	con.Port = os.Getenv("DB_PORT")

	connectionString := fmt.Sprintf("host=%s user=%s port=%s dbname=%s password=%s sslmode=%s", con.Host, con.User, con.Port, con.DBName, con.Password, con.SslMode)

	return connectionString, nil
}

func GetKey() []byte {
	return []byte(os.Getenv("JWT_SECRET_KEY"))
}

func GetPort() string {
	return os.Getenv("SERVER_PORT")
}

func GetClientID() string {
	return os.Getenv("CLIENT_ID")
}

func GetClientSecret() string {
	return os.Getenv("CLIENT_SECRET")
}

func GetMaxRetryAttempts() string {
	return os.Getenv("MAX_RETRY_ATTEMPTS")
}
