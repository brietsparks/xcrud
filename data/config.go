package data

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Vars struct {
	User string
	Password string
	Host string
	Port string
	Name string
}

// LoadEnvVars reads environment variables from a file and returns them as a Vars struct
func LoadEnvVars(filename string) (Vars, error) {
	err := godotenv.Load(filename)

	if err != nil {
		return Vars{}, fmt.Errorf("failed to load env vars: %w", err)
	}

	return Vars{
		User: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Host: os.Getenv("DB_HOST"),
		Port: os.Getenv("DB_PORT"),
		Name: os.Getenv("DB_NAME"),
	}, nil
}

// MakeUrl returns a postgres url from a Vars struct
func MakeUrl(vars Vars) string {
	return fmt.Sprintf("dbname=%s user=%s password=%s host=%s sslmode=disable",
		vars.Name,
		vars.User,
		vars.Password,
		vars.Host,
	)
}
