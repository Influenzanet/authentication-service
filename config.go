package main

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	yaml "gopkg.in/yaml.v2"
)

var conf config

type dbConf struct {
	CredentialsPath string `yaml:"credentials_path"`
	Address         string `yaml:"address"`
	Timeout         int    `yaml:"timeout"`
	DBNamePrefix string
}

type dbCredentials struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type jwtConf struct {
	SecretKey           string
	TokenExpiryInterval time.Duration
	TokenMinimumAgeMin  time.Duration
}

type config struct {
	ListenPort  int `yaml:"listen_port"`
	ServiceURLs struct {
		UserManagement string `yaml:"user_management"`
	} `yaml:"service_urls"`
	DB  dbConf `yaml:"db"`
	JWT jwtConf
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

	conf.DB.DBNamePrefix = "INF"


	// Get Token Attributes
	accessTokenExpiration, err := strconv.Atoi(os.Getenv("TOKEN_EXPIRATION_MIN"))
	if err != nil {
		log.Fatal(err.Error())
	}
	conf.JWT.TokenExpiryInterval = time.Minute * time.Duration(accessTokenExpiration)

	tokenMinAge, err := strconv.Atoi(os.Getenv("TOKEN_MINIMUM_AGE_MIN"))
	if err != nil {
		log.Fatal(err.Error())
	}
	conf.JWT.TokenMinimumAgeMin = time.Minute * time.Duration(tokenMinAge)
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
