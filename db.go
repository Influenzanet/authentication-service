package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbClient *mongo.Client

func dbInit() {
	dbCreds, err := readDBcredentials(conf.DB.CredentialsPath)
	if err != nil {
		log.Fatal(err)
	}

	// mongodb+srv://user-management-service:<PASSWORD>@influenzanettestdbcluster-pwvbz.mongodb.net/test?retryWrites=true
	address := fmt.Sprintf(`mongodb+srv://%s:%s@%s`, dbCreds.Username, dbCreds.Password, conf.DB.Address)

	dbClient, err = mongo.NewClient(options.Client().ApplyURI(address))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conf.DB.Timeout)*time.Second)
	defer cancel()

	err = dbClient.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func getCollection() *mongo.Collection {
	return dbClient.Database("default_token").Collection("tokens")
} //

func getContext() (ctx context.Context, cancel context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}
