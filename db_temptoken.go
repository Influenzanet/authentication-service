package main

import (
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func addTempTokenDB(t TempToken) (string, error) {
	ctx, cancel := getContext()
	defer cancel()

	filter := bson.M{"user_id": t.UserID, "purpose": t.Purpose, "instance_id": t.InstanceID}
	token := TempToken{}
	err := getCollection().FindOne(ctx, filter).Decode(&token)
	log.Println(err)
	log.Println(t)
	log.Println(token)
	log.Println(time.Now().Unix())
	if err == nil && !reachedExpirationTime(token.Expiration) {
		return "", errors.New("token with purpose already exists")
	}

	t.Token, err = generateUniqueTokenString()
	if err != nil {
		return "", err
	}

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

func getTempTokenFromDB(token string) (TempToken, error) {
	ctx, cancel := getContext()
	defer cancel()

	filter := bson.M{"token": token}

	t := TempToken{}
	err := getCollection().FindOne(ctx, filter).Decode(&t)
	return t, err
}

func deleteTempTokenDB(token string) error {
	// TODO: delete temporary token (defined by the token string)
	return errors.New("not implemented")
}

func deleteAllTempTokenForUserDB(instanceID string, userID string, purpose string) error {
	// TODO: delete all tokens from db with given user, instace and purpose
	return errors.New("not implemented")
}
