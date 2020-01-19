package main

import (
	"errors"

	"github.com/influenzanet/authentication-service/tokens"
	"go.mongodb.org/mongo-driver/bson"
)

func addTempTokenDB(t TempToken) (string, error) {
	ctx, cancel := getContext()
	defer cancel()

	filter := bson.M{"user_id": t.UserID, "purpose": t.Purpose, "instance_id": t.InstanceID}
	token := TempToken{}
	err := collectionRefTempToken().FindOne(ctx, filter).Decode(&token)

	if err == nil && !tokens.ReachedExpirationTime(token.Expiration) {
		return "", errors.New("token with purpose already exists")
	}

	t.Token, err = tokens.GenerateUniqueTokenString()
	if err != nil {
		return "", err
	}

	_, err = collectionRefTempToken().InsertOne(ctx, t)
	if err != nil {
		return "", err
	}

	return t.Token, nil
}

func getTempTokenForUserDB(instanceID string, uid string, purpose string) (tokens TempTokens, err error) {
	ctx, cancel := getContext()
	defer cancel()

	filter := bson.M{"instance_id": instanceID, "user_id": uid}
	if len(purpose) > 0 {
		filter["purpose"] = purpose
	}

	cur, err := collectionRefTempToken().Find(ctx, filter)
	if err != nil {
		return tokens, err
	}
	defer cur.Close(ctx)

	tokens = []TempToken{}
	for cur.Next(ctx) {
		var result TempToken
		err := cur.Decode(&result)
		if err != nil {
			return tokens, err
		}

		tokens = append(tokens, result)
	}
	if err := cur.Err(); err != nil {
		return tokens, err
	}
	return tokens, nil
}

func getTempTokenFromDB(token string) (TempToken, error) {
	ctx, cancel := getContext()
	defer cancel()

	filter := bson.M{"token": token}

	t := TempToken{}
	err := collectionRefTempToken().FindOne(ctx, filter).Decode(&t)
	return t, err
}

func deleteTempTokenDB(token string) error {
	ctx, cancel := getContext()
	defer cancel()

	filter := bson.M{"token": token}
	res, err := collectionRefTempToken().DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount < 1 {
		return errors.New("document not found")
	}
	return nil
}

func deleteAllTempTokenForUserDB(instanceID string, userID string, purpose string) error {
	ctx, cancel := getContext()
	defer cancel()

	filter := bson.M{"instance_id": instanceID, "user_id": userID}
	if len(purpose) > 0 {
		filter["purpose"] = purpose
	}
	res, err := collectionRefTempToken().DeleteMany(ctx, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount < 1 {
		return errors.New("document not found")
	}
	return nil
}
