package config

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/Mattcazz/Chat-TUI/client/internal/logger"
)

const jwtFile = "clit.jwt"
const netFile = "network.conf"
const styleConfigFile = "styles.conf"

type Colors struct{
	Username string `json:"username"`
	Text string `json:"text"`
	Border string `json:"border"`
}

type Network struct{
	ServerHost string `json:"serverhost"`
	ServerPort string `json:"serverport"`
}

type Config struct{
	Jwt string
	Network Network
	Colors Colors
	SSHKeyName string `json:"ssh_key_name"`
}

func getConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		logger.Log.Panicf("Could not get user's home directory: %s", err.Error())
	}

	return filepath.Join(home, ".config", "clit")
}

var Configuration Config

func LoadConfig() {
	_, err := os.Stat(getConfigDir())
	if errors.Is(err, os.ErrNotExist) {
		os.Mkdir(getConfigDir(), 0777)
	}

	jwtFilepath := filepath.Join(getConfigDir(), jwtFile)
	_, err = os.Stat(jwtFilepath)
	if errors.Is(err, os.ErrNotExist) {
		os.Create(jwtFilepath)
		// We don't need to populate it with anything for default
	}
	data, err := os.ReadFile(jwtFilepath)
	if err != nil {
		// Can't read file
		logger.Log.Panicf("Could not read file: %s", err.Error())
	}
	Configuration.Jwt = string(data)

	networkConfigFilepath := filepath.Join(getConfigDir(), netFile)
	_, err = os.Stat(networkConfigFilepath)
	if errors.Is(err, os.ErrNotExist) {
		networkConfigDefault := make(map[string]string)
		networkConfigDefault["serverhost"] = "localhost"
		networkConfigDefault["serverport"] = "8080"
		data, err = json.Marshal(networkConfigDefault)

		os.WriteFile(networkConfigFilepath, data, 0644)
	}
	networkConfigFile, err := os.Open(networkConfigFilepath)
	if err != nil {
		// Can't read file
		logger.Log.Panicf("Could not read file: %s", err.Error())
	}
	json.NewDecoder(networkConfigFile).Decode(&Configuration.Network)

	styleConfigFilepath := filepath.Join(getConfigDir(), styleConfigFile)
	_, err = os.Stat(styleConfigFilepath)
	if errors.Is(err, os.ErrNotExist) {
		styleConfigDefault := make(map[string]string)
		styleConfigDefault["username"] = "#A32CC4"
		styleConfigDefault["text"] = "#999999"
		styleConfigDefault["border"] = "#BBBBBB"
		data, err = json.Marshal(styleConfigDefault)

		os.WriteFile(styleConfigFilepath, data, 0644)
	}
	styleConfigFile, err := os.Open(styleConfigFilepath)
	if err != nil {
		// Can't read file
		logger.Log.Panicf("Could not read file: %s", err.Error())
	}
	json.NewDecoder(styleConfigFile).Decode(&Configuration.Colors)

	// General config
	configFilepath := filepath.Join(getConfigDir(), "clit.conf")
	_, err = os.Stat(configFilepath)
	if errors.Is(err, os.ErrNotExist) {
		configDefault := make(map[string]string)
		configDefault["ssh_key_name"] = "id_ed25519"
		data, err = json.Marshal(configDefault)

		os.WriteFile(configFilepath, data, 0644)
	}
	configFile, err := os.Open(configFilepath)
	if err != nil {
		// Can't read file
		logger.Log.Panicf("Could not read file: %s", err.Error())
	}
	json.NewDecoder(configFile).Decode(&Configuration)
}

func (c *Config) SetJWT(jwt string) {
	c.Jwt = jwt

	jwtFilepath := filepath.Join(getConfigDir(), jwtFile)
	os.WriteFile(jwtFilepath, []byte(jwt), 0644)
}
