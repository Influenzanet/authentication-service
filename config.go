package main

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

var conf config

type dbConf struct {
	CredentialsPath string `yaml:"credentials_path"`
	Address         string `yaml:"address"`
	Timeout         int    `yaml:"timeout"`
}

type dbCredentials struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type config struct {
	ListenPort  int `yaml:"listen_port"`
	ServiceURLs struct {
		UserManagement string `yaml:"user_management"`
	} `yaml:"service_urls"`
	DB dbConf `yaml:"db"`
}

func readConfig() {
	data, err := ioutil.ReadFile("./configs.yaml")
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal([]byte(data), &conf)
	if err != nil {
		log.Fatal(err)
	}
}

func readDBcredentials(path string) (dbCredentials, error) {
	var dbCreds dbCredentials
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return dbCreds, err
	}
	err = yaml.Unmarshal([]byte(data), &dbCreds)
	if err != nil {
		return dbCreds, err
	}
	return dbCreds, nil
}
