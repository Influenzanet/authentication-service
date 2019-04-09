package main

import (
	"go.mongodb.org/mongo-driver/bson"
)

func dbGetTokenByUserID(uid string) (TempToken, error) {
	ctx, cancel := getContext()
	defer cancel()

	filter := bson.M{"user_id": uid}

	t := TempToken{}
	err := getCollection().FindOne(ctx, filter).Decode(&t)

	return t, err
}

func dbCreateToken(t TempToken) (string, error) {
	ctx, cancel := getContext()
	defer cancel()

	t.Token = randomString()

	_, err := getCollection().InsertOne(ctx, t)
	if err != nil {
		return "", err
	}

	return t.Token, nil
}
