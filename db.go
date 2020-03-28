package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Collections
func collectionRefTempToken() *mongo.Collection {
	return dbClient.Database(conf.DB.DBNamePrefix + "global-infos").Collection("temp-tokens")
}

func collectionAppToken() *mongo.Collection {
	return dbClient.Database(conf.DB.DBNamePrefix + "global-infos").Collection("app-tokens")
}

// Connect to DB
func dbInit() {
	var err error
	dbClient, err = mongo.NewClient(
		options.Client().ApplyURI(conf.DB.URI),
		options.Client().SetMaxConnIdleTime(time.Duration(conf.DB.IdleConnTimeout)*time.Second),
		options.Client().SetMaxPoolSize(conf.DB.MaxPoolSize),
	)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := getContext()
	defer cancel()

	err = dbClient.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	ctx, conCancel := context.WithTimeout(context.Background(), time.Duration(conf.DB.Timeout)*time.Second)
	err = dbClient.Ping(ctx, nil)
	defer conCancel()
	if err != nil {
		log.Fatal("fail to connect to DB: " + err.Error())
	}
}
