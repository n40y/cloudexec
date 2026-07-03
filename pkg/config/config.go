package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type AWSConfig struct {
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	Region    string `yaml:"region"`
}

type Config struct {
	AWS AWSConfig `yaml:"aws"`
}

// LoadConfig charge le fichier de configuration s'il existe
func LoadConfig(filePath string) (*Config, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, nil // Aucun fichier de config, pas une erreur (utilisation des flags)
	}

	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
