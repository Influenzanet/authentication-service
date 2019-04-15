package main

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson"
)

func dbGetToken(uid string, purpose string) (TempToken, error) {
	ctx, cancel := getContext()
	defer cancel()

	filter := bson.M{"user_id": uid, "purpose": purpose}

	t := TempToken{}
	err := getCollection().FindOne(ctx, filter).Decode(&t)

	return t, err
}

func dbCreateToken(t TempToken) (string, error) {
	ctx, cancel := getContext()
	defer cancel()

	filter := bson.M{"user_id": t.UserID, "purpose": t.Purpose}
	token := TempToken{}
	err := getCollection().FindOne(ctx, filter).Decode(&token)
	if err == nil && !reachedExpirationTime(token.Expiration) {
		return "", errors.New("token with purpose already exists")
	}

	t.Token = randomString()

	_, err = getCollection().InsertOne(ctx, t)
	if err != nil {
		return "", err
	}

	return t.Token, nil
}

func dbDeleteToken(t TempToken) error {
	return nil
}
