package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

var ts *httptest.Server

// Define mock data for testing user-management service
var userList = []UserModel{
	UserModel{
		Email:    "test-p@test.com",
		Password: "testpassword", // is stored hashed on the real server
		ID:       1,
	},
	UserModel{
		Email:    "test-r@test.com",
		Password: "testpassword2", // is stored hashed on the real server
		ID:       2,
	},
	UserModel{
		Email:    "test-a@test.com",
		Password: "testpassword3", // is stored hashed on the real server
		ID:       3,
	},
}

var researcherList = []UserModel{
	UserModel{
		Email:    "test-r@test.com",
		Password: "testpassword2", // is stored hashed on the real server
		ID:       2,
	},
}

var adminList = []UserModel{
	UserModel{
		Email:    "test-a@test.com",
		Password: "testpassword3", // is stored hashed on the real server
		ID:       3,
	},
}

// MockLoginHandl is to emulate user-management's server response for login participant reuqests
func MockLoginParticipantHandl(context *gin.Context) {
	var creds userCredentials
	if err := context.ShouldBindJSON(&creds); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user with email address in the user list:
	var currentUser UserModel
	userFound := false
	for _, v := range userList {
		if v.Email == creds.Email {
			currentUser = v
			userFound = true
			break
		}
	}

	if !userFound {
		context.JSON(http.StatusForbidden, gin.H{"error": "wrong email or password"})
		return
	}

	// Check if password matches
	if currentUser.Password != creds.Password {
		context.JSON(http.StatusForbidden, gin.H{"error": "wrong email or password"})
		return
	}

	context.JSON(http.StatusOK, currentUser)
}

// This function is used for setup before executing the test functions
func TestMain(m *testing.M) {
	//Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	r.POST("/login-participant", MockLoginParticipantHandl)
	ts = httptest.NewServer(r)
	defer ts.Close()

	userManagementServer = ts.URL

	// Run the other tests
	os.Exit(m.Run())
}

func performRequest(r http.Handler, req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestLoginParticipant(t *testing.T) {
	r := gin.Default()
	r.POST("/v1/login/participant", loginParticipantHandl)

	/********************************************/
	/***** Check login with wrong email: *****/
	/********************************************/
	t.Run("Testing without payload", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/v1/login/participant", nil)
		w := performRequest(r, req)

		// Convert the JSON response to a map
		var response map[string]string
		if err := json.Unmarshal([]byte(w.Body.String()), &response); err != nil {
			t.Errorf("error parsing response body: %s", err.Error())
		}

		value, exists := response["error"]
		if w.Code != http.StatusBadRequest || !exists || value != "payload missing" {
			t.Errorf("status code: %d", w.Code)
			t.Errorf("response content: %s", w.Body.String())
			return
		}
	})

	/********************************************/
	/***** Check login with wrong email: *****/
	/********************************************/
	t.Run("Testing login with wrong email", func(t *testing.T) {
		loginData := &userCredentials{
			Email:    "test@test.com",
			Password: "testpassword",
		}
		loginPayload, _ := json.Marshal(loginData)

		req, _ := http.NewRequest("POST", "/v1/login/participant", bytes.NewBuffer(loginPayload))
		req.Header.Add("Content-Type", "application/json")
		w := performRequest(r, req)

		// Convert the JSON response to a map
		var response map[string]string
		if err := json.Unmarshal([]byte(w.Body.String()), &response); err != nil {
			t.Errorf("error parsing response body: %s", err.Error())
		}

		value, exists := response["error"]
		if w.Code != http.StatusForbidden || !exists || value != "wrong email or password" {
			t.Errorf("status code: %d", w.Code)
			t.Errorf("response content: %s", w.Body.String())
			return
		}
	})

	/********************************************/
	/***** Check login with wrong password: *****/
	/********************************************/
	t.Run("Testing login with wrong password", func(t *testing.T) {
		loginData := &userCredentials{
			Email:    "test-p@test.com",
			Password: "testpasswor",
		}
		loginPayload, _ := json.Marshal(loginData)

		req, _ := http.NewRequest("POST", "/v1/login/participant", bytes.NewBuffer(loginPayload))
		req.Header.Add("Content-Type", "application/json")
		w := performRequest(r, req)

		// Convert the JSON response to a map
		var response map[string]string
		if err := json.Unmarshal([]byte(w.Body.String()), &response); err != nil {
			t.Errorf("error parsing response body: %s", err.Error())
		}

		value, exists := response["error"]
		if w.Code != http.StatusForbidden || !exists || value != "wrong email or password" {
			t.Errorf("status code: %d", w.Code)
			t.Errorf("response content: %s", w.Body.String())
			return
		}
	})

	/********************************************/
	/***** Check login with correct email and password: ****/
	/********************************************/
	t.Run("Testing login with correct email and password", func(t *testing.T) {
		t.Logf("Testing login with correct email and password")
		loginData := &userCredentials{
			Email:    "test-p@test.com",
			Password: "testpassword",
		}
		loginPayload, _ := json.Marshal(loginData)

		req, _ := http.NewRequest("POST", "/v1/login/participant", bytes.NewBuffer(loginPayload))
		req.Header.Add("Content-Type", "application/json")
		w := performRequest(r, req)

		// Convert the JSON response to a map
		var response map[string]string
		if err := json.Unmarshal([]byte(w.Body.String()), &response); err != nil {
			t.Errorf("error parsing response body: %s", err.Error())
		}

		_, exists := response["token"]
		if w.Code != http.StatusOK || !exists {
			t.Errorf("status code: %d", w.Code)
			t.Errorf("response content: %s", w.Body.String())
			return
		}
	})

	/********************************************/
	/***** Check login with missing required fields: ****/
	/********************************************/
	t.Run("Testing login with missing required fields", func(t *testing.T) {
		loginData := &userCredentials{
			Email:    "",
			Password: "",
		}
		loginPayload, _ := json.Marshal(loginData)

		req, _ := http.NewRequest("POST", "/v1/login/participant", bytes.NewBuffer(loginPayload))
		req.Header.Add("Content-Type", "application/json")
		w := performRequest(r, req)

		// Convert the JSON response to a map
		var response map[string]string
		if err := json.Unmarshal([]byte(w.Body.String()), &response); err != nil {
			t.Errorf("error parsing response body: %s", err.Error())
		}

		_, exists := response["error"]
		if w.Code != http.StatusBadRequest || !exists {
			t.Errorf("status code: %d", w.Code)
			t.Errorf("response content: %s", w.Body.String())
			return
		}
	})
}

func getTokenForT() string {
	r := gin.Default()
	r.POST("/v1/login/participant", loginParticipantHandl)

	loginData := &userCredentials{
		Email:    "test-p@test.com",
		Password: "testpassword",
	}
	loginPayload, _ := json.Marshal(loginData)

	req, _ := http.NewRequest("POST", "/v1/login/participant", bytes.NewBuffer(loginPayload))
	req.Header.Add("Content-Type", "application/json")
	w := performRequest(r, req)

	// Convert the JSON response to a map
	var response map[string]string
	if err := json.Unmarshal([]byte(w.Body.String()), &response); err != nil {
		log.Fatal(err.Error())
	}

	value, _ := response["token"]
	return value
}

func TestRenewToken(t *testing.T) {
	tokenValidityPeriod = time.Second * 5
	minTokenAge = time.Second * 3

	r := gin.Default()
	r.GET("/v1/token/renew", renewTokenHandl)

	// Test without token
	t.Run("Testing renew token without token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/v1/token/renew", nil)
		// req.Header.Add("Authorization", "Bearer "+"")
		w := performRequest(r, req)

		// Convert the JSON response to a map
		var response map[string]string
		if err := json.Unmarshal([]byte(w.Body.String()), &response); err != nil {
			t.Errorf("error parsing response body: %s", err.Error())
		}

		_, exists := response["error"]
		if w.Code != http.StatusBadRequest || !exists {
			t.Errorf("status code: %d instead of %d", w.Code, http.StatusBadRequest)
			t.Errorf("response content: %s", w.Body.String())
			return
		}
	})

	// Test with empty token
	t.Run("Testing renew token with empty token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/v1/token/renew", nil)
		req.Header.Add("Authorization", "Bearer "+"")
		w := performRequest(r, req)

		// Convert the JSON response to a map
		var response map[string]string
		if err := json.Unmarshal([]byte(w.Body.String()), &response); err != nil {
			t.Errorf("error parsing response body: %s", err.Error())
		}

		_, exists := response["error"]
		if w.Code != http.StatusBadRequest || !exists {
			t.Errorf("status code: %d instead of %d", w.Code, http.StatusBadRequest)
			t.Errorf("response content: %s", w.Body.String())
			return
		}
	})

	// Test with empty token in url

	// Test with wrong token

	// Test with wrong token in url

	token := getTokenForT()

	log.Println(token)
	// Test eagerly, when min age not reached yet

	if testing.Short() {
		t.Skip("skipping renew token test in short mode, since it has to wait for token expiration.")
	}

	// Test renew after min age reached - wait 3 seconds - with header
	// Test renew after min age reached - wait 3 seconds - with url param

	// Test after expiration - wait 2 more seconds

}
