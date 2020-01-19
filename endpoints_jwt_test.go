package main

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api "github.com/influenzanet/authentication-service/api"
	api_mock "github.com/influenzanet/authentication-service/mocks"
	"github.com/influenzanet/authentication-service/tokens"
)

func TestLoginWithEmail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUserManagementClient := api_mock.NewMockUserManagementApiClient(mockCtrl)
	clients.userManagement = mockUserManagementClient

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
		req := &api.UserCredentials{}
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
		req := &api.UserCredentials{
			Email:      "wrong@test.com",
			Password:   "dfdfbmdpfbmd",
			InstanceId: testInstanceID,
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
		req := &api.UserCredentials{
			Email:      "test@test.com",
			Password:   "wrongpw",
			InstanceId: testInstanceID,
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
		req := &api.UserCredentials{
			Email:      "test@test.com",
			Password:   "dfdfbmdpfbmd",
			InstanceId: testInstanceID,
		}

		mockUserManagementClient.EXPECT().LoginWithEmail(
			gomock.Any(),
			gomock.Any(),
		).Return(&api.UserAuthInfo{
			UserId:     "testid",
			Roles:      []string{"participant"},
			InstanceId: testInstanceID,
		}, nil)
		mockUserManagementClient.EXPECT().TokenRefreshed(
			gomock.Any(),
			gomock.Any(),
		).Return(&api.Status{}, nil)

		resp, err := s.LoginWithEmail(context.Background(), req)
		if err != nil {
			st, _ := status.FromError(err)
			t.Errorf("unexpected error: %s", st.Message())
			return
		}
		if len(resp.AccessToken) < 1 || len(resp.RefreshToken) < 1 {
			t.Errorf("unexpected response: %s", resp)
		}
	})
}

func TestSignup(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUserManagementClient := api_mock.NewMockUserManagementApiClient(mockCtrl)
	clients.userManagement = mockUserManagementClient

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

	t.Run("with empty payload", func(t *testing.T) {
		req := &api.UserCredentials{}
		mockUserManagementClient.EXPECT().SignupWithEmail(
			gomock.Any(),
			gomock.Any(),
		).Return(nil, status.Error(codes.InvalidArgument, "missing arguments"))
		resp, err := s.SignupWithEmail(context.Background(), req)
		st, ok := status.FromError(err)
		if !ok || st == nil || st.Message() != "missing arguments" || resp != nil {
			t.Errorf("wrong error: %s", err.Error())
			t.Errorf("or response: %s", resp)
			return
		}
	})

	t.Run("with too short password", func(t *testing.T) {
		req := &api.UserCredentials{
			Email:      "test@test.com",
			Password:   "short",
			InstanceId: testInstanceID,
		}
		mockUserManagementClient.EXPECT().SignupWithEmail(
			gomock.Any(),
			gomock.Any(),
		).Return(nil, status.Error(codes.InvalidArgument, "password too weak"))

		resp, err := s.SignupWithEmail(context.Background(), req)
		st, ok := status.FromError(err)
		if !ok || st == nil || st.Message() != "password too weak" || resp != nil {
			t.Errorf("wrong error: %s", err.Error())
			t.Errorf("or response: %s", resp)
			return
		}
	})

	t.Run("with invalid email", func(t *testing.T) {
		req := &api.UserCredentials{
			Email:      "test-test.com",
			Password:   "short",
			InstanceId: testInstanceID,
		}
		mockUserManagementClient.EXPECT().SignupWithEmail(
			gomock.Any(),
			gomock.Any(),
		).Return(nil, status.Error(codes.InvalidArgument, "email not valid"))

		resp, err := s.SignupWithEmail(context.Background(), req)
		st, ok := status.FromError(err)
		if !ok || st == nil || st.Message() != "email not valid" || resp != nil {
			t.Errorf("wrong error: %s", err.Error())
			t.Errorf("or response: %s", resp)
			return
		}
	})

	t.Run("with existing user", func(t *testing.T) {
		req := &api.UserCredentials{
			Email:      "test@test.com",
			Password:   "short",
			InstanceId: testInstanceID,
		}
		mockUserManagementClient.EXPECT().SignupWithEmail(
			gomock.Any(),
			gomock.Any(),
		).Return(nil, status.Error(codes.Internal, "user already exists"))

		resp, err := s.SignupWithEmail(context.Background(), req)
		st, ok := status.FromError(err)
		if !ok || st == nil || st.Message() != "user already exists" || resp != nil {
			t.Errorf("wrong error: %s", err.Error())
			t.Errorf("or response: %s", resp)
			return
		}
	})

	t.Run("with valid arguments", func(t *testing.T) {
		req := &api.UserCredentials{
			Email:      "test@test.com",
			Password:   "short",
			InstanceId: testInstanceID,
		}

		mockUserManagementClient.EXPECT().SignupWithEmail(
			gomock.Any(),
			gomock.Any(),
		).Return(&api.UserAuthInfo{
			UserId:     "testid",
			Roles:      []string{"participant"},
			InstanceId: "test-inst",
		}, nil)
		mockUserManagementClient.EXPECT().TokenRefreshed(
			gomock.Any(),
			gomock.Any(),
		).Return(&api.Status{}, nil)

		resp, err := s.SignupWithEmail(context.Background(), req)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if len(resp.AccessToken) < 1 || len(resp.RefreshToken) < 1 {
			t.Errorf("unexpected response: %s", resp)
		}
	})
}

func TestValidateJWT(t *testing.T) {
	conf.JWT.TokenExpiryInterval = time.Second * 2
	conf.JWT.TokenMinimumAgeMin = time.Second * 1

	s := authServiceServer{}

	t.Run("without payload", func(t *testing.T) {
		resp, err := s.ValidateJWT(context.Background(), nil)
		st, ok := status.FromError(err)
		if !ok || st == nil || st.Message() != "missing arguments" || resp != nil {
			t.Errorf("wrong error: %s", err.Error())
			t.Errorf("or response: %s", resp)
			return
		}
	})

	t.Run("with empty payload", func(t *testing.T) {
		req := &api.JWTRequest{}

		resp, err := s.ValidateJWT(context.Background(), req)
		st, ok := status.FromError(err)
		if !ok || st == nil || st.Message() != "missing arguments" || resp != nil {
			t.Errorf("wrong error: %s", err.Error())
			t.Errorf("or response: %s", resp)
			return
		}
	})

	adminToken, err1 := tokens.GenerateNewToken("test-admin-id", []string{"PARTICIPANT", "ADMIN"}, testInstanceID, conf.JWT.TokenExpiryInterval)
	userToken, err2 := tokens.GenerateNewToken("test-user-id", []string{"PARTICIPANT"}, testInstanceID, conf.JWT.TokenExpiryInterval)
	if err1 != nil || err2 != nil {
		t.Errorf("unexpected error: %s or %s", err1, err2)
		return
	}

	t.Run("with wrong token", func(t *testing.T) {
		req := &api.JWTRequest{
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

	t.Run("with normal user token", func(t *testing.T) {
		req := &api.JWTRequest{
			Token: userToken,
		}

		resp, err := s.ValidateJWT(context.Background(), req)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		roles := tokens.GetRolesFromPayload(resp.Payload)
		if resp == nil || resp.InstanceId != testInstanceID || resp.Id != "test-user-id" || len(roles) != 1 || roles[0] != "PARTICIPANT" {
			t.Errorf("unexpected response: %s", resp)
			return
		}
	})

	t.Run("with admin token", func(t *testing.T) {
		req := &api.JWTRequest{
			Token: adminToken,
		}

		resp, err := s.ValidateJWT(context.Background(), req)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		roles := tokens.GetRolesFromPayload(resp.Payload)
		if resp == nil || len(roles) < 2 {
			t.Errorf("unexpected response: %s", resp)
			return
		}
	})

	if testing.Short() {
		t.Skip("skipping waiting for token test in short mode, since it has to wait for token expiration.")
	}
	time.Sleep(conf.JWT.TokenExpiryInterval + time.Second)

	t.Run("with expired token", func(t *testing.T) {
		req := &api.JWTRequest{
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

func TestRenewJWT(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUserManagementClient := api_mock.NewMockUserManagementApiClient(mockCtrl)
	clients.userManagement = mockUserManagementClient

	conf.JWT.TokenExpiryInterval = time.Second * 2
	conf.JWT.TokenMinimumAgeMin = time.Second * 1

	userToken, err := tokens.GenerateNewToken("test-user-id", []string{"PARTICIPANT"}, testInstanceID, conf.JWT.TokenExpiryInterval)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	refreshToken := "TEST-REFRESH-TOKEN-STRING"

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

	t.Run("with empty token", func(t *testing.T) {
		req := &api.RefreshJWTRequest{}

		resp, err := s.RenewJWT(context.Background(), req)
		st, ok := status.FromError(err)
		if !ok || st == nil || st.Message() != "missing arguments" || resp != nil {
			t.Errorf("wrong error: %s", err.Error())
			t.Errorf("or response: %s", resp)
			return
		}
	})

	t.Run("with wrong access token", func(t *testing.T) {
		req := &api.RefreshJWTRequest{
			AccessToken:  userToken + "x",
			RefreshToken: refreshToken,
		}
		resp, err := s.RenewJWT(context.Background(), req)
		st, ok := status.FromError(err)
		if !ok || st == nil || st.Message() != "wrong access token" || resp != nil {
			t.Errorf("wrong error: %s", err.Error())
			t.Errorf("or response: %s", resp)
			return
		}
	})

	// Test eagerly, when min age not reached yet
	t.Run("too eagerly", func(t *testing.T) {
		req := &api.RefreshJWTRequest{
			AccessToken:  userToken,
			RefreshToken: refreshToken,
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

	time.Sleep(conf.JWT.TokenMinimumAgeMin)

	t.Run("with wrong refresh token", func(t *testing.T) {
		req := &api.RefreshJWTRequest{
			AccessToken:  userToken,
			RefreshToken: userToken + "x",
		}
		mockUserManagementClient.EXPECT().CheckRefreshToken(
			gomock.Any(),
			gomock.Any(),
		).Return(nil, errors.New("wrong refresh token"))

		_, err := s.RenewJWT(context.Background(), req)
		if err == nil {
			t.Error("should fails with error")
			return
		}
		st, _ := status.FromError(err)
		if st.Message() != "wrong refresh token" {
			t.Errorf("unexpected error: %s", st.Message())
			return
		}
	})

	// Test renew after min age reached - wait 2 seconds
	t.Run("with normal tokens", func(t *testing.T) {
		req := &api.RefreshJWTRequest{
			AccessToken:  userToken,
			RefreshToken: refreshToken,
		}

		mockUserManagementClient.EXPECT().CheckRefreshToken(
			gomock.Any(),
			gomock.Any(),
		).Return(&api.Status{}, nil)
		mockUserManagementClient.EXPECT().TokenRefreshed(
			gomock.Any(),
			gomock.Any(),
		).Return(&api.Status{}, nil)

		resp, err := s.RenewJWT(context.Background(), req)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if resp == nil {
			t.Error("response is missing")
			return
		}
		if len(resp.AccessToken) < 10 || len(resp.RefreshToken) < 10 {
			t.Errorf("unexpected response: %s", resp)
			return
		}
	})

	time.Sleep(conf.JWT.TokenExpiryInterval)

	// Test with expired token
	t.Run("with expired token", func(t *testing.T) {
		req := &api.RefreshJWTRequest{
			AccessToken:  userToken,
			RefreshToken: refreshToken,
		}
		mockUserManagementClient.EXPECT().CheckRefreshToken(
			gomock.Any(),
			gomock.Any(),
		).Return(&api.Status{}, nil)
		mockUserManagementClient.EXPECT().TokenRefreshed(
			gomock.Any(),
			gomock.Any(),
		).Return(&api.Status{}, nil)

		resp, err := s.RenewJWT(context.Background(), req)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if resp == nil {
			t.Error("response is missing")
			return
		}
		if len(resp.AccessToken) < 10 || len(resp.RefreshToken) < 10 {
			t.Errorf("unexpected response: %s", resp)
			return
		}
	})
}
