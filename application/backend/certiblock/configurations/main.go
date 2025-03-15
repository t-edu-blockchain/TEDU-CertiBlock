package configurations

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type DatabaseConfigurations struct {
	Address      string
	Username     string
	Password     string
	DatabaseName string
}

type Configurations struct {
	Database  DatabaseConfigurations
	SecretKey string
}

func envMustBeSet(env string) string {
	value := os.Getenv(env)
	if value == "" {
		panic(fmt.Sprintf("Environment variable %s must be set", env))
	}
	return value
}

func envOptional(env string, defaultValue string) string {
	value := os.Getenv(env)
	if value == "" {
		return defaultValue
	}
	return value
}

func Load() (*Configurations, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	secretKey := envMustBeSet("APP_SECRET")

	if len(secretKey) < 32 {
		return nil, fmt.Errorf("APP_SECRET must be at least 32 characters long")
	}

	return &Configurations{
		Database: DatabaseConfigurations{
			Address:      envOptional("DATABASE_ADDRESS", "localhost:3306"),
			Username:     envMustBeSet("DATABASE_USERNAME"),
			Password:     envMustBeSet("DATABASE_PASSWORD"),
			DatabaseName: envMustBeSet("DATABASE_NAME"),
		},
		SecretKey: secretKey,
	}, nil
}
