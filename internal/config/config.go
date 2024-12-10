package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const gatorfile = ".gatorconfig.json"

type Config struct {
	DBURL       string `json:"db_url"`
	CURRENTUSER string `json:"current_user_name"`
}

func Read() (*Config, error) {
	cfg := Config{}
	home, err := returnHomeDir()
	if err != nil {
		log.Printf("Error reading the home directory: %s", err)
		return nil, err
	}
	gator_config_text, err := os.ReadFile(fmt.Sprintf("%s/.gatorconfig.json", home))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(gator_config_text, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil

}

func (c *Config) Setuser(current_user string) error {
	c.CURRENTUSER = current_user
	data, err := json.Marshal(&c)
	if err != nil {
		return err
	}
	home, err := returnHomeDir()
	if err != nil {
		return err
	}
	fullPath := fmt.Sprintf("%s/%s", home, gatorfile)
	err = os.WriteFile(fullPath, data, os.ModeAppend)
	if err != nil {
		return err
	}
	log.Printf("Successfully wrote to %s", fullPath)
	return nil
}

func returnHomeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Error reading the home directory: %s", err)
		return "", err
	}
	return fmt.Sprintf("%s", home), nil
}
