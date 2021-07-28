package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseHost                string
	DatabasePort                string
	DatabaseUser                string
	DatabaseName                string
	DatabasePassword            string
	DatabaseMaxConnection       int
	JwtTokenExpiration          int
	JwtSecret                   string
	Email                       string
	EmailPassword               string
	VerificationTokenExpiration int
}

var Conf Config

// LoadConfiguration Load the configuration for the database, security and Email
func LoadConfiguration() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	jwtTokenExpiration, err := strconv.Atoi(os.Getenv("JWT_TOKEN_EXPIRATION"))
	if err != nil {
		panic(err)
	}

	databaseMaxConnection, err := strconv.Atoi(os.Getenv("DATABASE_MAX_CONNECTION"))
	if err != nil {
		panic(err)
	}

	verificationTokenExpiration, err := strconv.Atoi(os.Getenv("VERIFICATION_TOKEN_EXPIRATION"))
	if err != nil {
		panic(err)
	}

	Conf = Config{
		DatabaseHost:                os.Getenv("DATABASE_HOST"),
		DatabasePort:                os.Getenv("DATABASE_PORT"),
		DatabaseUser:                os.Getenv("DATABASE_USER"),
		DatabaseName:                os.Getenv("DATABASE_NAME"),
		DatabasePassword:            os.Getenv("DATABASE_PASSWORD"),
		DatabaseMaxConnection:       databaseMaxConnection,
		JwtTokenExpiration:          jwtTokenExpiration,
		JwtSecret:                   os.Getenv("JWT_SECRET"),
		Email:                       os.Getenv("EMAIL"),
		EmailPassword:               os.Getenv("EMAIL_PASSWORD"),
		VerificationTokenExpiration: verificationTokenExpiration,
	}
}
