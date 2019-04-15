package main

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson"
)

func addTempTokenDB(t TempToken) (string, error) {
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

func getTempTokenForUserDB(instanceID string, uid string, purpose string) (TempToken, error) {
	ctx, cancel := getContext()
	defer cancel()

	filter := bson.M{"instance_id": instanceID, "user_id": uid, "purpose": purpose}

	t := TempToken{}
	err := getCollection().FindOne(ctx, filter).Decode(&t)
	return t, err
}

func getTempTokenFromDB(token string) error {
	// TODO: find one temp token by token string
	return errors.New("not implemented")
}

func deleteTempTokenDB(token string) error {
	// TODO: delete temporary token (defined by the token string)
	return errors.New("not implemented")
}

func deleteAllTempTokenForUserDB(instanceID string, userID string, purpose string) error {
	// TODO: delete all tokens from db with given user, instace and purpose
	return errors.New("not implemented")
}
