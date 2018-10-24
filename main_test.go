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

// refactoring to role arrays and combining role dbs into one
var userDB = []UserModel{
	UserModel{
		Email:    "test-p1@test.com",
		Password: "testpassword", // is stored hashed on the real server
		ID:       1,
		Roles:    []string{"PARTICIPANT"},
	},
	UserModel{
		Email:    "test-p2@test.com",
		Password: "testpassword2", // is stored hashed on the real server
		ID:       2,
		Roles:    []string{"PARTICIPANT"},
	},
	UserModel{
		Email:    "test-p3@test.com",
		Password: "testpassword3", // is stored hashed on the real server
		ID:       3,
		Roles:    []string{"PARTICIPANT"},
	},
	UserModel{
		Email:    "test-r1@test.com",
		Password: "testpassword4", // is stored hashed on the real server
		ID:       4,
		Roles:    []string{"PARTICIPANT", "RESEARCHER"},
	},
	UserModel{
		Email:    "test-r2@test.com",
		Password: "testpassword5", // is stored hashed on the real server
		ID:       5,
		Roles:    []string{"PARTICIPANT", "RESEARCHER"},
	},
	UserModel{
		Email:    "test-a1@test.com",
		Password: "testpassword6", // is stored hashed on the real server
		ID:       6,
		Roles:    []string{"PARTICIPANT", "RESEARCHER", "ADMIN"},
	},
	UserModel{
		Email:    "test-a2@test.com",
		Password: "testpassword7", // is stored hashed on the real server
		ID:       7,
		Roles:    []string{"PARTICIPANT", "RESEARCHER", "ADMIN"},
	},
}

// MockLoginHandl is to emulate user-management server response for login requests
func MockLoginHandl(context *gin.Context) {
	var creds userCredentials
	if err := context.ShouldBindJSON(&creds); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var currentUser UserModel
	userFound := false
	for _, v := range userDB {
		if v.Email == creds.Email && v.Password == creds.Password { // checking email and password right away instead of seperately?
			currentUser = v
			userFound = true
			break
		}
	}

	if !userFound {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "wrong email or password"})
		return
	}

	if creds.Role == "" {
		creds.Role = "PARTICIPANT"
	}

	if !currentUser.HasRole(creds.Role) {
		context.JSON(http.StatusForbidden, gin.H{"error": "missing appropriate role"})
		return
	}

	responseData := &UserLoginResponse{
		ID:   currentUser.ID,
		Role: creds.Role,
	}

	context.JSON(http.StatusOK, responseData)
}

// This function is used for setup before executing the test functions
func TestMain(m *testing.M) {
	//Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	r.POST("/login", MockLoginHandl)
	ts := httptest.NewServer(r)
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
	r.POST("/v1/user/login", loginHandl)

	const testingRole = "PARTICIPANT"

	/********************************************/
	/***** Check login without payload: *****/
	/********************************************/
	t.Run("Testing without payload", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/v1/user/login", nil)
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

		req, _ := http.NewRequest("POST", "/v1/user/login", bytes.NewBuffer(loginPayload))
		req.Header.Add("Content-Type", "application/json")
		w := performRequest(r, req)

		// Convert the JSON response to a map
		var response map[string]string
		if err := json.Unmarshal([]byte(w.Body.String()), &response); err != nil {
			t.Errorf("error parsing response body: %s", err.Error())
		}

		value, exists := response["error"]
		if w.Code != http.StatusUnauthorized || !exists || value != "wrong email or password" {
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
			Email:    "test-p1@test.com",
			Password: "testpasswor",
		}
		loginPayload, _ := json.Marshal(loginData)

		req, _ := http.NewRequest("POST", "/v1/user/login", bytes.NewBuffer(loginPayload))
		req.Header.Add("Content-Type", "application/json")
		w := performRequest(r, req)

		// Convert the JSON response to a map
		var response map[string]string
		if err := json.Unmarshal([]byte(w.Body.String()), &response); err != nil {
			t.Errorf("error parsing response body: %s", err.Error())
		}

		value, exists := response["error"]
		if w.Code != http.StatusUnauthorized || !exists || value != "wrong email or password" {
			t.Errorf("status code: %d", w.Code)
			t.Errorf("response content: %s", w.Body.String())
			return
		}
	})

	/********************************************/
	/***** Check login with correct email and password without role: ****/
	/********************************************/
	t.Run("Testing login with correct email and password without role", func(t *testing.T) {
		t.Logf("Testing login with correct email and password without role")
		loginData := &userCredentials{
			Email:    "test-p1@test.com",
			Password: "testpassword",
		}
		loginPayload, _ := json.Marshal(loginData)

		req, _ := http.NewRequest("POST", "/v1/user/login", bytes.NewBuffer(loginPayload))
		req.Header.Add("Content-Type", "application/json")
		w := performRequest(r, req)

		// Convert the JSON response to a map
		var response map[string]string
		if err := json.Unmarshal([]byte(w.Body.String()), &response); err != nil {
			t.Errorf("error parsing response body: %s", err.Error())
		}

		_, tokenExists := response["token"]
		roleValue, roleExists := response["role"]
		if w.Code != http.StatusOK || !tokenExists || !roleExists || roleValue != testingRole {
			t.Errorf("status code: %d", w.Code)
			t.Errorf("response content: %s", w.Body.String())
			return
		}
	})

	/********************************************/
	/***** Check login with correct email and password with role: ****/
	/********************************************/
	t.Run("Testing login with correct email and password with role", func(t *testing.T) {
		t.Logf("Testing login with correct email and password with role")
		loginData := &userCredentials{
			Email:    "test-p1@test.com",
			Password: "testpassword",
			Role:     "PARTICIPANT",
		}
		loginPayload, _ := json.Marshal(loginData)

		req, _ := http.NewRequest("POST", "/v1/user/login", bytes.NewBuffer(loginPayload))
		req.Header.Add("Content-Type", "application/json")
		w := performRequest(r, req)

		// Convert the JSON response to a map
		var response map[string]string
		if err := json.Unmarshal([]byte(w.Body.String()), &response); err != nil {
			t.Errorf("error parsing response body: %s", err.Error())
		}

		_, tokenExists := response["token"]
		roleValue, roleExists := response["role"]
		if w.Code != http.StatusOK || !tokenExists || !roleExists || roleValue != testingRole {
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

		req, _ := http.NewRequest("POST", "/v1/user/login", bytes.NewBuffer(loginPayload))
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

func TestLoginResearcher(t *testing.T) {
	r := gin.Default()
	r.POST("/v1/user/login", loginHandl)

	const testingRole = "RESEARCHER"

	/********************************************/
	/***** Check login with correct email and password as non researcher: ****/
	/********************************************/
	t.Run("Testing login with correct email and password as non researcher", func(t *testing.T) {
		t.Logf("Testing login with correct email and password but as non researcher")
		loginData := &userCredentials{
			Email:    "test-p1@test.com",
			Password: "testpassword",
			Role:     "RESEARCHER",
		}
		loginPayload, _ := json.Marshal(loginData)

		req, _ := http.NewRequest("POST", "/v1/user/login", bytes.NewBuffer(loginPayload))
		req.Header.Add("Content-Type", "application/json")
		w := performRequest(r, req)

		// Convert the JSON response to a map
		var response map[string]string
		if err := json.Unmarshal([]byte(w.Body.String()), &response); err != nil {
			t.Errorf("error parsing response body: %s", err.Error())
		}

		value, exists := response["error"]
		if w.Code != http.StatusForbidden || !exists || value != "missing appropriate role" {
			t.Errorf("status code: %d", w.Code)
			t.Errorf("response content: %s", w.Body.String())
			return
		}
	})

	/********************************************/
	/***** Check login with correct email and password as researcher: ****/
	/********************************************/
	t.Run("Testing login with correct email and password as researcher", func(t *testing.T) {
		t.Logf("Testing login with correct email and password as researcher")
		loginData := &userCredentials{
			Email:    "test-r1@test.com",
			Password: "testpassword4",
			Role:     "RESEARCHER",
		}
		loginPayload, _ := json.Marshal(loginData)

		req, _ := http.NewRequest("POST", "/v1/user/login", bytes.NewBuffer(loginPayload))
		req.Header.Add("Content-Type", "application/json")
		w := performRequest(r, req)

		// Convert the JSON response to a map
		var response map[string]string
		if err := json.Unmarshal([]byte(w.Body.String()), &response); err != nil {
			t.Errorf("error parsing response body: %s", err.Error())
		}

		_, tokenExists := response["token"]
		roleValue, roleExists := response["role"]
		if w.Code != http.StatusOK || !tokenExists || !roleExists || roleValue != testingRole {
			t.Errorf("status code: %d", w.Code)
			t.Errorf("response content: %s", w.Body.String())
			return
		}
	})
}

func TestLoginAdmin(t *testing.T) {
	r := gin.Default()
	r.POST("/v1/user/login", loginHandl)

	const testingRole = "ADMIN"

	/********************************************/
	/***** Check login with correct email and password as non admin: ****/
	/********************************************/
	t.Run("Testing login with correct email and password as non admin", func(t *testing.T) {
		t.Logf("Testing login with correct email and password but as non admin")
		loginData := &userCredentials{
			Email:    "test-r1@test.com",
			Password: "testpassword4",
			Role:     "ADMIN",
		}
		loginPayload, _ := json.Marshal(loginData)

		req, _ := http.NewRequest("POST", "/v1/user/login", bytes.NewBuffer(loginPayload))
		req.Header.Add("Content-Type", "application/json")
		w := performRequest(r, req)

		// Convert the JSON response to a map
		var response map[string]string
		if err := json.Unmarshal([]byte(w.Body.String()), &response); err != nil {
			t.Errorf("error parsing response body: %s", err.Error())
		}

		value, exists := response["error"]
		if w.Code != http.StatusForbidden || !exists || value != "missing appropriate role" {
			t.Errorf("status code: %d", w.Code)
			t.Errorf("response content: %s", w.Body.String())
			return
		}
	})

	/********************************************/
	/***** Check login with correct email and password as admin: ****/
	/********************************************/
	t.Run("Testing login with correct email and password as admin", func(t *testing.T) {
		t.Logf("Testing login with correct email and password as admin")
		loginData := &userCredentials{
			Email:    "test-a1@test.com",
			Password: "testpassword6",
			Role:     "ADMIN",
		}
		loginPayload, _ := json.Marshal(loginData)

		req, _ := http.NewRequest("POST", "/v1/user/login", bytes.NewBuffer(loginPayload))
		req.Header.Add("Content-Type", "application/json")
		w := performRequest(r, req)

		// Convert the JSON response to a map
		var response map[string]string
		if err := json.Unmarshal([]byte(w.Body.String()), &response); err != nil {
			t.Errorf("error parsing response body: %s", err.Error())
		}

		_, tokenExists := response["token"]
		roleValue, roleExists := response["role"]
		if w.Code != http.StatusOK || !tokenExists || !roleExists || roleValue != testingRole {
			t.Errorf("status code: %d", w.Code)
			t.Errorf("response content: %s", w.Body.String())
			return
		}
	})
}

func getTokenForParticipant() string {
	r := gin.Default()
	r.POST("/v1/user/login", loginHandl)

	loginData := &userCredentials{
		Email:    "test-p1@test.com",
		Password: "testpassword",
	}
	loginPayload, _ := json.Marshal(loginData)

	req, _ := http.NewRequest("POST", "/v1/user/login", bytes.NewBuffer(loginPayload))
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
	tokenValidityPeriod = time.Second * 2
	minTokenAge = time.Second * 1

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
	t.Run("Test with empty token in url", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/v1/token/renew?token=", nil)
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

	badToken := "eydfsdfsdffsdfsdf.w45345sdfsdvcsdsdf.435345fsdf-4rwefsdfsd" // "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJyb2xlIjoicGFydGljaXBhbnQiLCJleHAiOjE1Mzk0MTc0MzAsImlhdCI6MTUzOTQxNzQyNX0.klxofJLg5J31v7hKO7TbPrceBzyYlp9kIAJuotUmY11pk08Hnn2uHtuDfdqBWVtcI_lQ-vKiikVs5icrewyQOXMzTesQXI41SZvRdEQfit1MZ5syE0a2PODRFsizaqT5vqVN04ZzX_3iPEvSBP25wMy8R4dzYaY5XcR2heJWIxaNFd3w65UDa_mNk4u3Oem7XO1Ufn_-ay98XqAUg5Zo0TI9sk2WQF57pzXAlHMVmCMNW1bP_OPra9CCQb2pUm2sKJiAgWVOBVB4lz50VoTsoJimQoTc5UpF3SCujL-Yt5mh7d7EUvDkKoSuqd5Pc8iKHs1Ix9jSmtoLpPxmCAnepA"

	// Test with wrong token
	t.Run("Testing with wrong token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/v1/token/renew", nil)
		req.Header.Add("Authorization", "Bearer "+badToken)
		w := performRequest(r, req)

		// Convert the JSON response to a map
		var response map[string]string
		if err := json.Unmarshal([]byte(w.Body.String()), &response); err != nil {
			t.Errorf("error parsing response body: %s", err.Error())
		}

		_, exists := response["error"]
		if w.Code != http.StatusUnauthorized || !exists {
			t.Errorf("status code: %d instead of %d", w.Code, http.StatusUnauthorized)
			t.Errorf("response content: %s", w.Body.String())
			return
		}
	})

	// Test with wrong token in url
	t.Run("Test with wrong token in url", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/v1/token/renew?token="+badToken, nil)
		w := performRequest(r, req)

		// Convert the JSON response to a map
		var response map[string]string
		if err := json.Unmarshal([]byte(w.Body.String()), &response); err != nil {
			t.Errorf("error parsing response body: %s", err.Error())
		}

		_, exists := response["error"]
		if w.Code != http.StatusUnauthorized || !exists {
			t.Errorf("status code: %d instead of %d", w.Code, http.StatusUnauthorized)
			t.Errorf("response content: %s", w.Body.String())
			return
		}
	})

	token := getTokenForParticipant()
	// log.Println(token)

	// Test eagerly, when min age not reached yet
	t.Run("Testing token too eagerly", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/v1/token/renew", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		w := performRequest(r, req)

		// Convert the JSON response to a map
		var response map[string]string
		if err := json.Unmarshal([]byte(w.Body.String()), &response); err != nil {
			t.Errorf("error parsing response body: %s", err.Error())
		}

		_, exists := response["error"]
		if w.Code != http.StatusTeapot || !exists {
			t.Errorf("status code: %d instead of %d", w.Code, http.StatusTeapot)
			t.Errorf("response content: %s", w.Body.String())
			return
		}
	})

	if testing.Short() {
		t.Skip("skipping renew token test in short mode, since it has to wait for token expiration.")
	}

	time.Sleep(minTokenAge)

	// Test renew after min age reached - wait 2 seconds - with header
	t.Run("Testing token with header param", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/v1/token/renew", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		w := performRequest(r, req)

		// Convert the JSON response to a map
		var response map[string]string
		if err := json.Unmarshal([]byte(w.Body.String()), &response); err != nil {
			t.Errorf("error parsing response body: %s", err.Error())
		}

		_, exists := response["token"]
		if w.Code != http.StatusOK || !exists {
			t.Errorf("status code: %d instead of %d", w.Code, http.StatusOK)
			t.Errorf("response content: %s", w.Body.String())
			return
		}
	})

	// Test renew after min age reached - wait 2 seconds - with url param
	t.Run("Testing token with url param", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/v1/token/renew?token="+token, nil)
		w := performRequest(r, req)

		// Convert the JSON response to a map
		var response map[string]string
		if err := json.Unmarshal([]byte(w.Body.String()), &response); err != nil {
			t.Errorf("error parsing response body: %s", err.Error())
		}

		_, exists := response["token"]
		if w.Code != http.StatusOK || !exists {
			t.Errorf("status code: %d instead of %d", w.Code, http.StatusOK)
			t.Errorf("response content: %s", w.Body.String())
			return
		}
	})

	time.Sleep(tokenValidityPeriod)
	// Test with expired token
	t.Run("Testing with expired token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/v1/token/renew", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		w := performRequest(r, req)

		// Convert the JSON response to a map
		var response map[string]string
		if err := json.Unmarshal([]byte(w.Body.String()), &response); err != nil {
			t.Errorf("error parsing response body: %s", err.Error())
		}

		_, exists := response["error"]
		if w.Code != http.StatusUnauthorized || !exists {
			t.Errorf("status code: %d instead of %d", w.Code, http.StatusUnauthorized)
			t.Errorf("response content: %s", w.Body.String())
			return
		}
	})

}
