package main

import (
	"testing"
	"time"
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
		resUnix := getExpirationTime(time.Hour * 24 * -5)
		resTime := time.Unix(resUnix, 0)
		expected := time.Now().AddDate(0, 0, -5)
		if resTime.Year() != expected.Year() || resTime.Month() != expected.Month() || resTime.Day() != expected.Day() {
			t.Errorf("date values don't match. result: %s, expected %s", resTime.String(), expected.String())
			return
		}
	})

	t.Run("with zero days", func(t *testing.T) {
		resUnix := getExpirationTime(time.Hour * 24 * 0)
		resTime := time.Unix(resUnix, 0)
		expected := time.Now()
		if resTime.Year() != expected.Year() || resTime.Month() != expected.Month() || resTime.Day() != expected.Day() {
			t.Errorf("date values don't match. result: %s, expected %s", resTime.String(), expected.String())
			return
		}
	})

	t.Run("with positive days", func(t *testing.T) {
		resUnix := getExpirationTime(time.Hour * 24 * 5)
		resTime := time.Unix(resUnix, 0)
		expected := time.Now().AddDate(0, 0, 5)
		if resTime.Year() != expected.Year() || resTime.Month() != expected.Month() || resTime.Day() != expected.Day() {
			t.Errorf("date values don't match. result: %s, expected %s", resTime.String(), expected.String())
			return
		}
	})
}

func TestReachedExpirationTime(t *testing.T) {
	t.Run("before expiration", func(t *testing.T) {
		exp := getExpirationTime(time.Hour * 1)
		isExp := reachedExpirationTime(exp)
		if isExp {
			t.Error("expiration should not be reached yet")
		}
	})

	t.Run("after expiration", func(t *testing.T) {
		exp := getExpirationTime(time.Hour * -1)
		isExp := reachedExpirationTime(exp)
		if !isExp {
			t.Error("should be expired now")
		}
	})
}

func TestGetRolesFromPayload(t *testing.T) {
	t.Run("with empty payload", func(t *testing.T) {
		p := map[string]string{}
		roles := getRolesFromPayload(p)
		if len(roles) > 0 {
			t.Errorf("something went wrong: %s", roles)
		}
	})

	t.Run("with missing role field", func(t *testing.T) {
		p := map[string]string{
			"test": "testRole",
		}
		roles := getRolesFromPayload(p)
		if len(roles) > 0 {
			t.Errorf("something went wrong: %s", roles)
		}
	})
	t.Run("with one role", func(t *testing.T) {
		p := map[string]string{
			"roles": "testRole1",
		}
		roles := getRolesFromPayload(p)
		if len(roles) != 1 && roles[0] != "testRole1" {
			t.Errorf("something went wrong: %s", roles)
		}
	})

	t.Run("with multiple roles", func(t *testing.T) {
		p := map[string]string{
			"roles": "testRole1,testRole2,testRole3",
		}
		roles := getRolesFromPayload(p)
		if len(roles) != 3 && roles[0] != "testRole1" {
			t.Errorf("something went wrong: %s", roles)
		}
	})
}
