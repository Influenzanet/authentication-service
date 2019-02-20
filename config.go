package main

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

var conf config

type config struct {
	ListenPort  int `yaml:"listen_port"`
	ServiceURLs struct {
		UserManagement string `yaml:"user_management"`
	} `yaml:"service_urls"`
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
