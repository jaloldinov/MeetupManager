package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Conf struct {
	DBName      string `yaml:"db_name"`
	DBUsername  string `yaml:"db_username"`
	DBPassword  string `yaml:"db_password"`
	Port        string `yaml:"port"`
	JWTKey      string `yaml:"jwt_key"`
	DefaultLang string `yaml:"default_lang"`
}

func GetConf() *Conf {

	cfg := Conf{}

	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return &cfg
}
