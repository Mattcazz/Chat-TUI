package config

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
)

const jwt_file = "clit.jwt"
const net_file = "network.conf"
const style_config_file = "styles.conf"

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
	SSH_key_name string `json:"ssh_key_name"`
}

func getConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Panic(err.Error())
	}

	return filepath.Join(home, ".config", "clit")
}

var Configuration Config

func LoadConfig() {
	_, err := os.Stat(getConfigDir())
	if errors.Is(err, os.ErrNotExist) {
		os.Mkdir(getConfigDir(), 0777)
	}

	jwt_filepath := filepath.Join(getConfigDir(), jwt_file)
	_, err = os.Stat(jwt_filepath)
	if errors.Is(err, os.ErrNotExist) {
		os.Create(jwt_filepath)
		// We don't need to populate it with anything for default
	}
	data, err := os.ReadFile(jwt_filepath)
	if err != nil {
		// TODO
	}
	Configuration.Jwt = string(data)

	network_config_filepath := filepath.Join(getConfigDir(), net_file)
	_, err = os.Stat(network_config_filepath)
	if errors.Is(err, os.ErrNotExist) {
		network_config_default := make(map[string]string)
		network_config_default["serverhost"] = "localhost"
		network_config_default["serverport"] = "8080"
		data, err = json.Marshal(network_config_default)

		os.WriteFile(network_config_filepath, data, 0644)
	}
	network_config_file, err := os.Open(network_config_filepath)
	if err != nil {
		// TODO
	}
	json.NewDecoder(network_config_file).Decode(&Configuration.Network)

	style_config_filepath := filepath.Join(getConfigDir(), style_config_file)
	_, err = os.Stat(style_config_filepath)
	if errors.Is(err, os.ErrNotExist) {
		style_config_default := make(map[string]string)
		style_config_default["username"] = "#A32CC4"
		style_config_default["text"] = "#999999"
		style_config_default["border"] = "#BBBBBB"
		data, err = json.Marshal(style_config_default)

		os.WriteFile(style_config_filepath, data, 0644)
	}
	style_config_file, err := os.Open(style_config_filepath)
	if err != nil {
		// TODO
	}
	json.NewDecoder(style_config_file).Decode(&Configuration.Colors)

	// General config
	config_filepath := filepath.Join(getConfigDir(), "clit.conf")
	_, err = os.Stat(config_filepath)
	if errors.Is(err, os.ErrNotExist) {
		config_default := make(map[string]string)
		config_default["ssh_key_name"] = "id_ed25519"
		data, err = json.Marshal(config_default)

		os.WriteFile(config_filepath, data, 0644)
	}
	config_file, err := os.Open(config_filepath)
	if err != nil {
		// TODO
	}
	json.NewDecoder(config_file).Decode(&Configuration)
}

func (c *Config) SetJWT(jwt string) {
	c.Jwt = jwt

	jwt_filepath := filepath.Join(getConfigDir(), jwt_file)
	os.WriteFile(jwt_filepath, []byte(jwt), 0644)
}
