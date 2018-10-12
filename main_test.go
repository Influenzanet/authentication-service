package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

var ts *httptest.Server

// Define mock data for testing user-management service
var userList = []userModel{
	userModel{
		Email:    "test-p@test.com",
		Password: "testpassword", // is stored hashed on the real server
		ID:       1,
	},
	userModel{
		Email:    "test-r@test.com",
		Password: "testpassword2", // is stored hashed on the real server
		ID:       2,
	},
	userModel{
		Email:    "test-a@test.com",
		Password: "testpassword3", // is stored hashed on the real server
		ID:       3,
	},
}

var researcherList = []userModel{
	userModel{
		Email:    "test-r@test.com",
		Password: "testpassword2", // is stored hashed on the real server
		ID:       2,
	},
}

var adminList = []userModel{
	userModel{
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
	var currentUser userModel
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

func performRequest(r http.Handler, method, path string, payload *bytes.Buffer) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, payload)
	if payload != nil {
		req.Header.Add("Content-Type", "application/json")
	}
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
	t.Run("Testing login with wrong email", func(t *testing.T) {
		loginData := &userCredentials{
			Email:    "test@test.com",
			Password: "testpassword",
		}
		loginPayload, _ := json.Marshal(loginData)

		w := performRequest(r, "POST", "/v1/login/participant", bytes.NewBuffer(loginPayload))

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

		w := performRequest(r, "POST", "/v1/login/participant", bytes.NewBuffer(loginPayload))

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

		w := performRequest(r, "POST", "/v1/login/participant", bytes.NewBuffer(loginPayload))

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

		w := performRequest(r, "POST", "/v1/login/participant", bytes.NewBuffer(loginPayload))

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

func TestRefreshToken(t *testing.T) {
	tokenValidityPeriod = time.Second * 10
	minTokenAge = time.Second * 3
}
