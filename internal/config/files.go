package config

import (
	"fmt"
	"os"
)

type DBConfig struct {
	DBHost     string 
	DBName     string 
	DBUser     string 
	DBPassword string 
	SSLEnabled string 
}

func LoadDBConfig() (*DBConfig, error) {
	var c DBConfig

	c.DBName = os.Getenv("POSTGRES_DB")
	c.DBHost = os.Getenv("DATABASE_HOST")
	c.DBUser = os.Getenv("POSTGRES_USER")
	c.DBPassword = os.Getenv("POSTGRES_PASSWORD")
	c.SSLEnabled = "disable"

	return &c, nil
}

func (c *DBConfig) GetDBN() string {
	return fmt.Sprintf(
		"dbname=%v user=%v password=%v sslmode=%v host=%v",
		c.DBName, c.DBUser, c.DBPassword, c.SSLEnabled, c.DBHost,
	)
}
