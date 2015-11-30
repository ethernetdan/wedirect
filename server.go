package main // import "github.com/ethernetdan/populardns"

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/pearkes/cloudflare"
)

const defaultConfigFile = "config.json"

var client cloudflare.Client

func init() {
	configFile := defaultConfigFile
	if configFileEnv := os.Getenv("CONFIG_FILE"); len(configFileEnv) != 0 {
		configFile = configFileEnv
	}

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Panicf("Could not read configuration file: %v", err)
	}

	err = json.Unmarshal(data, &client)
	if err != nil {
		log.Panicf("Could not unmarshal configuration struct from file `%s`: %v", configFile, err)
	}
}
