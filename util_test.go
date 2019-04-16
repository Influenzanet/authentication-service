package main

import (
	"testing"
)

func TestGenerateUniqueTokenString(t *testing.T) {
	t.Run("test result", func(t *testing.T) {
		nrTest := 10000
		if testing.Short() {
			nrTest = 100
		}
		res := []string{}
		for i := 0; i <= nrTest; i++ {
			token, err := generateUniqueTokenString()
			if err != nil {
				t.Errorf("unexpected error: %s", err.Error())
				return
			}
			for _, tV := range res {
				if token == tV {
					t.Errorf("duplicated token: %s", token)
					return
				}
			}
			res = append(res, token)
		}
	})
}

func TestGetExpirationTime(t *testing.T) {
	t.Run("with negative days", func(t *testing.T) {
		t.Error("test not implemented")
	})

	t.Run("with zero days", func(t *testing.T) {
		t.Error("test not implemented")
	})

	t.Run("with positive days", func(t *testing.T) {
		t.Error("test not implemented")
	})
}

func TestReachedExpirationTime(t *testing.T) {
	t.Run("before expiration", func(t *testing.T) {
		t.Error("test not implemented")
	})

	t.Run("exactly at expiration", func(t *testing.T) {
		t.Error("test not implemented")
	})

	t.Run("after expiration", func(t *testing.T) {
		t.Error("test not implemented")
	})
}
