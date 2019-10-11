package main

import (
	"testing"
	"time"
)

// Testing Database Interface methods
func TestDbInterfaceMethodsForTempToken(t *testing.T) {
	testTempToken := TempToken{
		UserID:     "test_user_id",
		Purpose:    "test_purpose1",
		InstanceID: testInstanceID,
		Expiration: getExpirationTime(10 * time.Second),
	}
	tokenStr := ""

	t.Run("Add temporary token to DB", func(t *testing.T) {
		ts, err := addTempTokenDB(testTempToken)
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		tokenStr = ts
	})

	t.Run("Try to add duplicate temporary token to DB", func(t *testing.T) {
		_, err := addTempTokenDB(testTempToken)
		if err == nil {
			t.Errorf("created duplicate TempToken")
			return
		}
	})

	t.Run("try to get temporary token by wrong token string", func(t *testing.T) {
		tempToken, err := getTempTokenFromDB(tokenStr + "++")
		if err == nil && tempToken.UserID != "" {
			t.Error(tempToken)
			t.Error("token should not be found")
			return
		}
	})

	t.Run("get temporary token by token string", func(t *testing.T) {
		tempToken, err := getTempTokenFromDB(tokenStr)
		if err != nil {
			t.Error("token not found by token string")
			return
		}
		if tempToken.UserID != testTempToken.UserID || tempToken.Purpose != testTempToken.Purpose || tempToken.Expiration != testTempToken.Expiration {
			t.Error("temp token does not match")
			t.Error(tempToken)
			return
		}
	})

	t.Run("try to get temporary token by wrong user id", func(t *testing.T) {
		getTempTokenForUserDB("todo", "todo", "todo")
		t.Error("test not implemented")
	})

	t.Run("try to get temporary token by wrong instace id", func(t *testing.T) {
		getTempTokenForUserDB("todo", "todo", "todo")
		t.Error("test not implemented")
	})

	t.Run("try to get temporary token by wrong purpose", func(t *testing.T) {
		getTempTokenForUserDB("todo", "todo", "todo")
		t.Error("test not implemented")
	})

	t.Run("get temporary token by user_id+instance_id", func(t *testing.T) {
		getTempTokenForUserDB("todo", "todo", "todo")
		t.Error("test not implemented")
	})

	t.Run("Try delete not existing temporary token", func(t *testing.T) {
		deleteTempTokenDB("token")
		t.Error("test not implemented")
	})

	t.Run("Delete temporary token", func(t *testing.T) {
		deleteTempTokenDB("token")
		t.Error("test not implemented")
	})

	t.Run("Delete all temporary token of a user_id with empty instance_id", func(t *testing.T) {
		deleteAllTempTokenForUserDB("", "", "")
		// TODO: should fail, instance id should not be emtpy
		t.Error("test not implemented")
	})

	t.Run("Try to delete all temporary token of a user_id with wrong id, correct instance_id", func(t *testing.T) {
		deleteAllTempTokenForUserDB("", "", "")
		t.Error("test not implemented")
	})

	t.Run("Delete all temporary token of a user_id+instance_id", func(t *testing.T) {
		deleteAllTempTokenForUserDB("", "", "")
		t.Error("test not implemented")
	})
}
