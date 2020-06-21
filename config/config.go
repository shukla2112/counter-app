package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/gomodule/redigo/redis"

	rh "github.com/shukla2112/counter-app/redis"
)

// CliArgs : Cli args passed to this service
type CliArgs struct {
	Mode string `json:"mode"`
}

// ReadCliArgs : Read CLI args passed to this service
func ReadCliArgs() (CliArgs, string) {
	mode := flag.String("mode", "dev", "Mode for counter-app (development/staging/production)")
	flag.Parse()

	cliArgs := CliArgs{
		Mode: *mode,
	}
	jsonBytes, _ := json.MarshalIndent(cliArgs, "", "    ")
	return cliArgs, string(jsonBytes)
}

// AppConfig : Type to store App level config like db connection
type AppConfig struct {
	ConfigData Config
	Args       *CliArgs
	RedisPool  *redis.Pool
}

// Config : Holds config data
type Config struct {
	Server string `json:"server"`
	Redis  string `json:"redis"`
}

// NewAppConfig : Construct new app config object with db connections
func NewAppConfig(cliArgs *CliArgs) (appC *AppConfig, err error) {
	// Get config data
	configData, err := GetConfig(*cliArgs)
	if err != nil {
		return nil, err
	}

	log.Printf("CONFIG_1001: Connecting to redis at: %s\n", configData.Redis)
	redisPool := rh.RNewPool(configData.Redis)

	appC = &AppConfig{
		ConfigData: configData,
		Args:       cliArgs,
		RedisPool:  redisPool,
	}
	return appC, nil
}

// GetConfig : Read config file according "mode" cli argument passed in
func GetConfig(cliArgs CliArgs) (Config, error) {
	configFile := fmt.Sprintf("config/%s.json", cliArgs.Mode)
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		configFile = fmt.Sprintf("../config/%s.json", cliArgs.Mode)
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			configFile = fmt.Sprintf("../../config/%s.json", cliArgs.Mode)
		}
	}

	// Read config file
	configJSON, err := ioutil.ReadFile(configFile)
	if err != nil {
		return Config{}, err
	}

	var configData Config
	err = json.Unmarshal([]byte(configJSON), &configData)
	if err != nil {
		return Config{}, err
	}

	return configData, nil
}
