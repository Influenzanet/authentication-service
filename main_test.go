package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

var ts *httptest.Server

func MockLoginHandl(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"msg": "ok"})
}

// This function is used for setup before executing the test functions
func TestMain(m *testing.M) {
	//Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	r.POST("/login", MockLoginHandl)
	ts = httptest.NewServer(r)
	defer ts.Close()

	userManagementServer = ts.URL

	// Run the other tests
	os.Exit(m.Run())
}

func TestLoginParticipant(t *testing.T) {
	w := httptest.NewRecorder()

	r := gin.Default()

	r.POST("/v1/login/participant", loginParticipantHandl)

	loginData := &userCredentials{
		Email:    "test@test.de",
		Password: "testpassword",
	}

	loginPayload, err := json.Marshal(loginData)

	req, _ := http.NewRequest("POST", "/v1/login/participant", bytes.NewBuffer(loginPayload))
	req.Header.Add("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("login failed")
		t.Fail()
	}

	p, err := ioutil.ReadAll(w.Body)
	if err != nil || strings.Index(string(p), "<title>Successful registration &amp; Login</title>") < 0 {
		t.Errorf("body wrong")
		t.Fail()
	}
}
