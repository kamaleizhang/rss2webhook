package r2w

import (
	"encoding/json"
	"os"

	"github.com/jmoiron/sqlx"
)

type ConfigManager interface {
	GetConfigs() ([]Config, error)
}

type ConfigManagerImpl struct {
	repo               *sqlx.DB
	loadFromConfigFile string
}

func NewConfigManager(repo *sqlx.DB, loadFromConfigFile string) ConfigManager {
	return &ConfigManagerImpl{
		repo:               repo,
		loadFromConfigFile: loadFromConfigFile,
	}
}

func (c *ConfigManagerImpl) GetConfigs() ([]Config, error) {
	if c.loadFromConfigFile != "" {
		return c.loadFromFile()
	}
	return c.loadFromDB()
}

func (c *ConfigManagerImpl) loadFromFile() ([]Config, error) {
	data, err := os.ReadFile(c.loadFromConfigFile)
	if err != nil {
		return nil, err
	}

	var configs []Config
	if err := json.Unmarshal(data, &configs); err != nil {
		return nil, err
	}

	return configs, nil
}

func (c *ConfigManagerImpl) loadFromDB() ([]Config, error) {
	var configs []Config
	err := c.repo.Select(&configs, "SELECT domain, rss_url, target_webhook FROM configs")
	if err != nil {
		return nil, err
	}
	return configs, nil
}
