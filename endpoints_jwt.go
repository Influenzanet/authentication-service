package main

import (
	"context"
	"log"
	"strings"
	"time"

	"google.golang.org/grpc/codes"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/status"

	api "github.com/influenzanet/authentication-service/api"
)

func (s *authServiceServer) Status(ctx context.Context, _ *empty.Empty) (*api.Status, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (s *authServiceServer) LoginWithEmail(ctx context.Context, req *api.UserCredentials) (*api.TokenResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid username and/or password")
	}
	resp, err := userManagementClient.LoginWithEmail(context.Background(), req)
	if err != nil {
		log.Printf("error during login with email: %s", err.Error())
		return nil, status.Error(codes.InvalidArgument, "invalid username and/or password")
	}

	// generate tokens
	token, err := generateNewToken(resp.UserId, resp.Roles, resp.InstanceId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	rt, err := generateUniqueTokenString()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// submit token to user management
	_, err = userManagementClient.TokenRefreshed(context.Background(), &api.UserReference{
		InstanceId: resp.InstanceId,
		Token:      rt,
		UserId:     resp.UserId,
	})
	if err != nil {
		st := status.Convert(err)
		log.Printf("error during signup with email: %s: %s", st.Code(), st.Message())
		return nil, status.Error(codes.Internal, st.Message())
	}

	return &api.TokenResponse{
		AccessToken:  token,
		RefreshToken: rt,
		ExpiresIn:    int32(conf.JWT.TokenExpiryInterval / time.Minute),
	}, nil
}

func (s *authServiceServer) SignupWithEmail(ctx context.Context, req *api.UserCredentials) (*api.TokenResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "missing arguments")
	}
	resp, err := userManagementClient.SignupWithEmail(context.Background(), req)
	if err != nil {
		st := status.Convert(err)
		log.Printf("error during signup with email: %s: %s", st.Code(), st.Message())
		return nil, status.Error(codes.Internal, st.Message())
	}

	// generate tokens
	token, err := generateNewToken(resp.UserId, resp.Roles, resp.InstanceId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	rt, err := generateUniqueTokenString()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// submit refresh token to user management
	_, err = userManagementClient.TokenRefreshed(context.Background(), &api.UserReference{
		InstanceId: resp.InstanceId,
		Token:      rt,
		UserId:     resp.UserId,
	})
	if err != nil {
		st := status.Convert(err)
		log.Printf("error during signup with email: %s: %s", st.Code(), st.Message())
		return nil, status.Error(codes.Internal, st.Message())
	}

	return &api.TokenResponse{
		AccessToken:  token,
		RefreshToken: rt,
		ExpiresIn:    int32(conf.JWT.TokenExpiryInterval / time.Minute),
	}, nil
}

func (s *authServiceServer) ValidateJWT(ctx context.Context, req *api.JWTRequest) (*api.TokenInfos, error) {
	if req == nil || req.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "missing arguments")
	}
	// Parse and validate token
	parsedToken, ok, err := validateToken(req.Token)
	if err != nil || !ok {
		return nil, status.Error(codes.InvalidArgument, "invalid token")
	}

	return &api.TokenInfos{
		Id:         parsedToken.ID,
		InstanceId: parsedToken.InstanceID,
		IssuedAt:   parsedToken.IssuedAt,
		Payload:    parsedToken.Payload,
	}, nil
}

func (s *authServiceServer) RenewJWT(ctx context.Context, req *api.RefreshJWTRequest) (*api.TokenResponse, error) {
	if req == nil || req.AccessToken == "" || req.RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "missing arguments")
	}

	// Parse and validate token
	parsedToken, _, err := validateToken(req.AccessToken)
	if err != nil && !strings.Contains(err.Error(), "token is expired by") {
		return nil, status.Error(codes.PermissionDenied, "wrong access token")
	}

	// Check for too frequent requests:
	if checkTokenAgeMaturity(parsedToken.StandardClaims.IssuedAt) {
		return nil, status.Error(codes.Unavailable, "can't renew token so often")
	}

	// check refresh token from user management
	_, err = userManagementClient.CheckRefreshToken(context.Background(), &api.UserReference{
		Token:      req.RefreshToken,
		UserId:     parsedToken.ID,
		InstanceId: parsedToken.InstanceID,
	})
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "wrong refresh token") // err
	}

	roles := getRolesFromPayload(parsedToken.Payload)

	// Generate new access token:
	newToken, err := generateNewToken(parsedToken.ID, roles, parsedToken.InstanceID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	newRefreshToken, err := generateUniqueTokenString()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// submit refresh token to user management
	_, err = userManagementClient.TokenRefreshed(context.Background(), &api.UserReference{
		UserId:     parsedToken.ID,
		InstanceId: parsedToken.InstanceID,
		Token:      newRefreshToken,
	})
	if err != nil {
		st := status.Convert(err)
		log.Printf("error during token refresh: %s: %s", st.Code(), st.Message())
		return nil, status.Error(codes.Internal, st.Message())
	}

	return &api.TokenResponse{
		AccessToken:  newToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int32(conf.JWT.TokenExpiryInterval / time.Minute),
	}, nil
}
