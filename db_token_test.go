package main

import (
	"testing"
)

// Testing Database Interface methods
func TestDbInterfaceMethods(t *testing.T) {
	testTempToken := TempToken{
		UserID:     "test_user_id",
		Purpose:    "test_purpose1",
		Expiration: getExpirationTime(10),
	}

	t.Run("Testing dbCreateToken", func(t *testing.T) {
		_, err := dbCreateToken(testTempToken)
		if err != nil {
			t.Errorf(err.Error())
			return
		}
	})

	t.Run("Testing dbCreateToken with duplicate", func(t *testing.T) {
		_, err := dbCreateToken(testTempToken)
		if err == nil {
			t.Errorf("created duplicate TempToken")
			return
		}
	})
}
