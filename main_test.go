package main

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"google.golang.org/grpc/status"

	influenzanet "github.com/influenzanet/api/dist/go"
	auth_api "github.com/influenzanet/api/dist/go/auth-service"
	user_api "github.com/influenzanet/api/dist/go/user-management"
	um_mock "github.com/influenzanet/authentication-service/mock_user_management"
)

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
		st, ok := status.FromError(err)
		if !ok || err == nil || st.Message() != "invalid username and/or password" || resp != nil {
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
		st, ok := status.FromError(err)
		if !ok || st == nil || st.Message() != "invalid username and/or password" || resp != nil {
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
		st, ok := status.FromError(err)
		if !ok || st == nil || st.Message() != "invalid username and/or password" || resp != nil {
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
		st, ok := status.FromError(err)
		if !ok || st == nil || st.Message() != "invalid username and/or password" || resp != nil {
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
			st, _ := status.FromError(err)
			t.Errorf("unexpected error: %s", st.Message())
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
		st, ok := status.FromError(err)
		if !ok || st == nil || st.Message() != "missing arguments" || resp != nil {
			t.Errorf("wrong error: %s", err.Error())
			t.Errorf("or response: %s", resp)
			return
		}
	})

	t.Run("Testing signup with empty payload", func(t *testing.T) {
		req := &user_api.NewUser{}

		resp, err := s.SignupWithEmail(context.Background(), req)
		st, ok := status.FromError(err)
		if !ok || st == nil || st.Message() != "missing arguments" || resp != nil {
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
		st, ok := status.FromError(err)
		if !ok || st == nil || st.Message() != "password too weak" || resp != nil {
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
		st, ok := status.FromError(err)
		if !ok || st == nil || st.Message() != "email not valid" || resp != nil {
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
		st, ok := status.FromError(err)
		if !ok || st == nil || st.Message() != "user already exists" || resp != nil {
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
		st, ok := status.FromError(err)
		if !ok || st == nil || st.Message() != "missing arguments" || resp != nil {
			t.Errorf("wrong error: %s", err.Error())
			t.Errorf("or response: %s", resp)
			return
		}
	})

	t.Run("Testing token validation with empty payload", func(t *testing.T) {
		req := &auth_api.EncodedToken{}

		resp, err := s.ValidateJWT(context.Background(), req)
		st, ok := status.FromError(err)
		if !ok || st == nil || st.Message() != "missing arguments" || resp != nil {
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
		st, ok := status.FromError(err)
		if !ok || st == nil || st.Message() != "invalid token" || resp != nil {
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
		st, ok := status.FromError(err)
		if !ok || st == nil || st.Message() != "invalid token" || resp != nil {
			t.Errorf("wrong error: %s", err.Error())
			t.Errorf("or response: %s", resp)
			return
		}
	})
}

func TestRenewToken(t *testing.T) {
	tokenValidityPeriod = time.Second * 2
	minTokenAge = time.Second * 1

	userToken, err := generateNewToken("test-user-id", []string{"PARTICIPANT"}, "test-instance")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	s := authServiceServer{}

	t.Run("Testing token refresh without token", func(t *testing.T) {
		resp, err := s.RenewJWT(context.Background(), nil)
		st, ok := status.FromError(err)
		if !ok || st == nil || st.Message() != "missing arguments" || resp != nil {
			t.Errorf("wrong error: %s", err.Error())
			t.Errorf("or response: %s", resp)
			return
		}
	})

	t.Run("Test token refresh with empty token", func(t *testing.T) {
		req := &auth_api.EncodedToken{}

		resp, err := s.RenewJWT(context.Background(), req)
		st, ok := status.FromError(err)
		if !ok || st == nil || st.Message() != "missing arguments" || resp != nil {
			t.Errorf("wrong error: %s", err.Error())
			t.Errorf("or response: %s", resp)
			return
		}
	})

	t.Run("Test token refresh with wrong token", func(t *testing.T) {
		req := &auth_api.EncodedToken{
			Token: userToken + "x",
		}

		resp, err := s.RenewJWT(context.Background(), req)
		st, ok := status.FromError(err)
		if !ok || st == nil || st.Message() != "invalid token" || resp != nil {
			t.Errorf("wrong error: %s", err.Error())
			t.Errorf("or response: %s", resp)
			return
		}
	})

	// Test eagerly, when min age not reached yet
	t.Run("Testing token too eagerly", func(t *testing.T) {
		req := &auth_api.EncodedToken{
			Token: userToken,
		}

		resp, err := s.RenewJWT(context.Background(), req)
		st, ok := status.FromError(err)
		if !ok || st == nil || st.Message() != "can't renew token so often" || resp != nil {
			t.Errorf("wrong error: %s", err.Error())
			t.Errorf("or response: %s", resp)
			return
		}
	})

	if testing.Short() {
		t.Skip("skipping renew token test in short mode, since it has to wait for token expiration.")
	}

	time.Sleep(minTokenAge)

	// Test renew after min age reached - wait 2 seconds
	t.Run("Testing token refresh", func(t *testing.T) {
		req := &auth_api.EncodedToken{
			Token: userToken,
		}

		resp, err := s.RenewJWT(context.Background(), req)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if resp == nil || len(resp.Token) < 10 {
			t.Errorf("unexpected response: %s", resp)
			return
		}
	})

	time.Sleep(tokenValidityPeriod)
	// Test with expired token
	t.Run("Testing with expired token", func(t *testing.T) {
		req := &auth_api.EncodedToken{
			Token: userToken,
		}
		resp, err := s.RenewJWT(context.Background(), req)
		st, ok := status.FromError(err)
		if !ok || st == nil || st.Message() != "invalid token" || resp != nil {
			t.Errorf("wrong error: %s", err.Error())
			t.Errorf("or response: %s", resp)
			return
		}
	})
}
