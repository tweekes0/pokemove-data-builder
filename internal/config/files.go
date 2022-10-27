package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

const (
	configFolder = "./conf"
)

var (
	DBFile = configFile("db.toml")
)

type DBConfig struct {
	DBName     string `toml:"DatabaseName"`
	DBUser     string `toml:"DatabaseUser"`
	DBPassword string `toml:"DatabasePassword"`
	SSLEnabled string `toml:"SSLEnabled"`
}

func ReadDBConfig() (*DBConfig, error) {
	f, err := ioutil.ReadFile(DBFile)
	if err != nil {
		return nil, err
	}

	var c DBConfig
	if err = toml.Unmarshal(f, &c); err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *DBConfig) GetDBN() string {
	return fmt.Sprintf(
		"dbname=%v user=%v password=%v sslmode=%v",
		c.DBName, c.DBUser, c.DBPassword, c.SSLEnabled,
	)
}

func configFile(filename string) string {
	return filepath.Join(configFolder, filename)
}
