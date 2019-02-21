package main

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	influenzanet "github.com/Influenzanet/api/dist/go"
	auth_api "github.com/Influenzanet/api/dist/go/auth-service"
	user_api "github.com/Influenzanet/api/dist/go/user-management"
	um_mock "github.com/Influenzanet/authentication-service/mock_user_management"
)

// UserModel holds information relevant for authentication
type UserModel struct {
	Email             string   `json:"email"`
	Password          string   `json:"password"`
	ID                string   `json:"user_id"`
	Roles             []string `json:"roles"`
	AuthenticatedRole string   `json:"authenticated_role"`
}

// HasRole checks whether the user has a specified role
func (u UserModel) HasRole(role string) bool {
	for _, v := range u.Roles {
		if v == role {
			return true
		}
	}
	return false
}

// Mock user DB
var userDB = []UserModel{
	UserModel{
		Email:    "test-p1@test.com",
		Password: "testpassword", // is stored hashed on the real server
		ID:       "1",
		Roles:    []string{"PARTICIPANT"},
	},
	UserModel{
		Email:    "test-p2@test.com",
		Password: "testpassword2", // is stored hashed on the real server
		ID:       "2",
		Roles:    []string{"PARTICIPANT"},
	},
	UserModel{
		Email:    "test-p3@test.com",
		Password: "testpassword3", // is stored hashed on the real server
		ID:       "3",
		Roles:    []string{"PARTICIPANT"},
	},
	UserModel{
		Email:    "test-r1@test.com",
		Password: "testpassword4", // is stored hashed on the real server
		ID:       "4",
		Roles:    []string{"PARTICIPANT", "RESEARCHER"},
	},
	UserModel{
		Email:    "test-r2@test.com",
		Password: "testpassword5", // is stored hashed on the real server
		ID:       "5",
		Roles:    []string{"PARTICIPANT", "RESEARCHER"},
	},
	UserModel{
		Email:    "test-a1@test.com",
		Password: "testpassword6", // is stored hashed on the real server
		ID:       "6",
		Roles:    []string{"PARTICIPANT", "RESEARCHER", "ADMIN"},
	},
	UserModel{
		Email:    "test-a2@test.com",
		Password: "testpassword7", // is stored hashed on the real server
		ID:       "7",
		Roles:    []string{"PARTICIPANT", "RESEARCHER", "ADMIN"},
	},
}

// This function is used for setup before executing the test functions
func TestMain(m *testing.M) {
	// Run the other tests
	os.Exit(m.Run())
}

func TestLoginWithEmail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUserManagementClient := um_mock.NewMockUserManagementApiClient(mockCtrl)
	userManagementClient = mockUserManagementClient

	s := authServiceServer{}

	t.Run("Testing login without payload", func(t *testing.T) {
		resp, err := s.LoginWithEmail(context.Background(), nil)
		if err == nil || err.Error() != "invalid username and/or password" || resp != nil {
			t.Errorf("wrong error: %s", err.Error())
			t.Errorf("or response: %s", resp)
			return
		}
	})

	t.Run("Testing login with empty payload", func(t *testing.T) {
		req := &influenzanet.UserCredentials{}
		mockUserManagementClient.EXPECT().LoginWithEmail(
			gomock.Any(),
			gomock.Any(),
		).Return(nil, errors.New("invalid username and/or password"))

		resp, err := s.LoginWithEmail(context.Background(), req)
		if err == nil || err.Error() != "invalid username and/or password" || resp != nil {
			t.Errorf("wrong error: %s", err.Error())
			t.Errorf("or response: %s", resp)
			return
		}
	})

	t.Run("Testing login with wrong email", func(t *testing.T) {
		req := &influenzanet.UserCredentials{
			Email:      "wrong@test.com",
			Password:   "dfdfbmdpfbmd",
			InstanceId: "test-inst",
		}
		mockUserManagementClient.EXPECT().LoginWithEmail(
			gomock.Any(),
			gomock.Any(),
		).Return(nil, errors.New("invalid username and/or password"))

		resp, err := s.LoginWithEmail(context.Background(), req)
		if err == nil || err.Error() != "invalid username and/or password" || resp != nil {
			t.Errorf("wrong error: %s", err.Error())
			t.Errorf("or response: %s", resp)
			return
		}
	})

	t.Run("Testing login with wrong password", func(t *testing.T) {
		req := &influenzanet.UserCredentials{
			Email:      "test@test.com",
			Password:   "wrongpw",
			InstanceId: "test-inst",
		}
		mockUserManagementClient.EXPECT().LoginWithEmail(
			gomock.Any(),
			gomock.Any(),
		).Return(nil, errors.New("invalid username and/or password"))

		resp, err := s.LoginWithEmail(context.Background(), req)
		if err == nil || err.Error() != "invalid username and/or password" || resp != nil {
			t.Errorf("wrong error: %s", err.Error())
			t.Errorf("or response: %s", resp)
			return
		}
	})

	t.Run("Testing login with correct email and password", func(t *testing.T) {
		req := &influenzanet.UserCredentials{
			Email:      "test@test.com",
			Password:   "dfdfbmdpfbmd",
			InstanceId: "test-inst",
		}

		mockUserManagementClient.EXPECT().LoginWithEmail(
			gomock.Any(),
			gomock.Any(),
		).Return(&user_api.UserAuthInfo{
			UserId:     "testid",
			Roles:      []string{"participant"},
			InstanceId: "test-inst",
		}, nil)

		resp, err := s.LoginWithEmail(context.Background(), req)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if len(resp.Token) < 1 {
			t.Errorf("unexpected response: %s", resp)
		}
	})
}

func TestSignup(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUserManagementClient := um_mock.NewMockUserManagementApiClient(mockCtrl)
	userManagementClient = mockUserManagementClient

	s := authServiceServer{}

	t.Run("Testing signup without payload", func(t *testing.T) {
		resp, err := s.SignupWithEmail(context.Background(), nil)
		if err == nil || err.Error() != "missing arguments" || resp != nil {
			t.Errorf("wrong error: %s", err.Error())
			t.Errorf("or response: %s", resp)
			return
		}
	})

	t.Run("Testing signup with empty payload", func(t *testing.T) {
		req := &user_api.NewUser{}

		resp, err := s.SignupWithEmail(context.Background(), req)
		if err == nil || err.Error() != "missing arguments" || resp != nil {
			t.Errorf("wrong error: %s", err.Error())
			t.Errorf("or response: %s", resp)
			return
		}
	})

	t.Run("Testing signup with too short password", func(t *testing.T) {
		req := &user_api.NewUser{
			Auth: &influenzanet.UserCredentials{
				Email:      "test@test.com",
				Password:   "short",
				InstanceId: "test-inst",
			},
			Profile: &user_api.Profile{},
		}
		mockUserManagementClient.EXPECT().SignupWithEmail(
			gomock.Any(),
			gomock.Any(),
		).Return(nil, errors.New("password too weak"))

		resp, err := s.SignupWithEmail(context.Background(), req)
		if err == nil || err.Error() != "password too weak" || resp != nil {
			t.Errorf("wrong error: %s", err.Error())
			t.Errorf("or response: %s", resp)
			return
		}
	})

	t.Run("Testing signup with invalid email", func(t *testing.T) {
		req := &user_api.NewUser{
			Auth: &influenzanet.UserCredentials{
				Email:      "test-test.com",
				Password:   "short",
				InstanceId: "test-inst",
			},
			Profile: &user_api.Profile{},
		}
		mockUserManagementClient.EXPECT().SignupWithEmail(
			gomock.Any(),
			gomock.Any(),
		).Return(nil, errors.New("email not valid"))

		resp, err := s.SignupWithEmail(context.Background(), req)
		if err == nil || err.Error() != "email not valid" || resp != nil {
			t.Errorf("wrong error: %s", err.Error())
			t.Errorf("or response: %s", resp)
			return
		}
	})

	t.Run("Testing signup with existing user", func(t *testing.T) {
		req := &user_api.NewUser{
			Auth: &influenzanet.UserCredentials{
				Email:      "test@test.com",
				Password:   "short",
				InstanceId: "test-inst",
			},
			Profile: &user_api.Profile{},
		}
		mockUserManagementClient.EXPECT().SignupWithEmail(
			gomock.Any(),
			gomock.Any(),
		).Return(nil, errors.New("user already exists"))

		resp, err := s.SignupWithEmail(context.Background(), req)
		if err == nil || err.Error() != "user already exists" || resp != nil {
			t.Errorf("wrong error: %s", err.Error())
			t.Errorf("or response: %s", resp)
			return
		}
	})

	t.Run("Testing signup with valid arguments", func(t *testing.T) {
		req := &user_api.NewUser{
			Auth: &influenzanet.UserCredentials{
				Email:      "test@test.com",
				Password:   "short",
				InstanceId: "test-inst",
			},
			Profile: &user_api.Profile{},
		}

		mockUserManagementClient.EXPECT().SignupWithEmail(
			gomock.Any(),
			gomock.Any(),
		).Return(&user_api.UserAuthInfo{
			UserId:     "testid",
			Roles:      []string{"participant"},
			InstanceId: "test-inst",
		}, nil)

		resp, err := s.SignupWithEmail(context.Background(), req)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if len(resp.Token) < 1 {
			t.Errorf("unexpected response: %s", resp)
		}
	})
}

func TestValidateToken(t *testing.T) {
	tokenValidityPeriod = time.Second * 2
	minTokenAge = time.Second * 1

	s := authServiceServer{}

	t.Run("Testing token validation without payload", func(t *testing.T) {
		resp, err := s.ValidateJWT(context.Background(), nil)
		if err == nil || err.Error() != "missing arguments" || resp != nil {
			t.Errorf("wrong error: %s", err.Error())
			t.Errorf("or response: %s", resp)
			return
		}
	})

	t.Run("Testing token validation with empty payload", func(t *testing.T) {
		req := &auth_api.EncodedToken{}

		resp, err := s.ValidateJWT(context.Background(), req)
		if err == nil || err.Error() != "missing arguments" || resp != nil {
			t.Errorf("wrong error: %s", err.Error())
			t.Errorf("or response: %s", resp)
			return
		}
	})

	adminToken, err1 := generateNewToken("test-admin-id", []string{"PARTICIPANT", "ADMIN"}, "test-instance")
	userToken, err2 := generateNewToken("test-user-id", []string{"PARTICIPANT"}, "test-instance")
	if err1 != nil || err2 != nil {
		t.Errorf("unexpected error: %s or %s", err1, err2)
		return
	}

	t.Run("Test token validation with wrong token", func(t *testing.T) {
		req := &auth_api.EncodedToken{
			Token: adminToken + "x",
		}

		resp, err := s.ValidateJWT(context.Background(), req)
		if err == nil || err.Error() != "invalid token" || resp != nil {
			t.Errorf("wrong error: %s", err.Error())
			t.Errorf("or response: %s", resp)
			return
		}
	})

	t.Run("Test token validation with normal user token", func(t *testing.T) {
		req := &auth_api.EncodedToken{
			Token: userToken,
		}

		resp, err := s.ValidateJWT(context.Background(), req)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if resp == nil || resp.InstanceId != "test-instance" ||
			resp.UserId != "test-user-id" || len(resp.Roles) != 1 || resp.Roles[0] != "PARTICIPANT" {
			t.Errorf("unexpected response: %s", resp)
			return
		}
	})

	t.Run("Test token validation with admin token", func(t *testing.T) {
		req := &auth_api.EncodedToken{
			Token: adminToken,
		}

		resp, err := s.ValidateJWT(context.Background(), req)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if resp == nil || len(resp.Roles) < 2 {
			t.Errorf("unexpected response: %s", resp)
			return
		}
	})

	if testing.Short() {
		t.Skip("skipping waiting for token test in short mode, since it has to wait for token expiration.")
	}
	time.Sleep(tokenValidityPeriod + time.Second)

	t.Run("Test with expired token", func(t *testing.T) {
		req := &auth_api.EncodedToken{
			Token: adminToken,
		}
		resp, err := s.ValidateJWT(context.Background(), req)
		if err == nil || err.Error() != "invalid token" || resp != nil {
			t.Errorf("wrong error: %s", err.Error())
			t.Errorf("or response: %s", resp)
			return
		}
	})
}

func TestRenewToken(t *testing.T) {
	t.Errorf("test not implemented")
	/*
		tokenValidityPeriod = time.Second * 2
		minTokenAge = time.Second * 1

		r := gin.Default()
		r.GET("/v1/token/renew", middlewares.ExtractToken(), renewTokenHandl)

		// Test without token
		t.Run("Testing renew token without token", func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/v1/token/renew", nil)
			// req.Header.Add("Authorization", "Bearer "+"")
			w := performRequest(r, req)

			// Convert the JSON response to a map
			var response map[string]interface{}
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
			var response map[string]interface{}
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
			var response map[string]interface{}
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
			var response map[string]interface{}
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
			var response map[string]interface{}
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
			var response map[string]interface{}
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
	*/
}
