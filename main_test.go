package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	"google.golang.org/grpc/status"
)

var (
	testInstanceID = strconv.FormatInt(time.Now().Unix(), 10)
)

func dropTestDB() {
	log.Println("Drop test database")
	ctx, cancel := getContext()
	defer cancel()

	// Drop Test Instance Auth DB
	err := dbClient.Database(conf.DB.DBNamePrefix + "global-infos").Drop(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

// This function is used for setup before executing the test functions
func TestMain(m *testing.M) {
	result := m.Run()
	dropTestDB()
	os.Exit(result)
}

func shouldHaveGrpcErrorStatus(err error, expectedError string) (bool, string) {
	if err == nil {
		return false, "should return an error"
	}
	st, ok := status.FromError(err)
	if !ok || st == nil {
		return false, fmt.Sprintf("unexpected error: %s", err.Error())
	}

	if len(expectedError) > 0 && st.Message() != expectedError {
		return false, fmt.Sprintf("wrong error: %s", st.Message())
	}
	return true, ""
}
