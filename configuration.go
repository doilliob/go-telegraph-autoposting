package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Configuration struct {
	TelegraphAcountName string `yaml:"telegraph_account_name"`
	MaxImageSizeMb      int    `yaml:"maximum_image_size_mb"`
	Debug               bool   `yaml:"debug"`
}

func readConfigurationFile(fileName string) Configuration {
	config := Configuration{}

	// Read File
	yamlData, err := ioutil.ReadFile(fileName)
	checkError(err)

	// Unmarshal YAML data to Config structure
	err = yaml.Unmarshal(yamlData, &config)
	checkError(err)

	return config
}
