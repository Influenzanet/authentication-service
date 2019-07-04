package main

import (
	"context"
	"log"

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

	token, err := generateNewToken(resp.UserId, resp.Roles, resp.InstanceId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// TODO: generate refresh token

	return &api.TokenResponse{
		AccessToken:  token,
		RefreshToken: "todo",
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

	token, err := generateNewToken(resp.UserId, resp.Roles, resp.InstanceId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// TODO: generate refresh token

	return &api.TokenResponse{
		AccessToken:  token,
		RefreshToken: "todo",
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

	// TODO: check refresh token

	// Parse and validate token
	parsedToken, ok, err := validateToken(req.AccessToken)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid token") // err
	}
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "invalid token")
	}

	// Check for too frequent requests:
	if checkTokenAgeMaturity(parsedToken.StandardClaims.IssuedAt) {
		return nil, status.Error(codes.Unavailable, "can't renew token so often")
	}

	roles := getRolesFromPayload(parsedToken.Payload)

	// Generate new token:
	newToken, err := generateNewToken(parsedToken.ID, roles, parsedToken.InstanceID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	updateReq := api.UserReference{
		UserId: parsedToken.ID,
	}
	_, err = userManagementClient.TokenRefreshed(context.Background(), &updateReq)
	if err != nil {
		st := status.Convert(err)
		log.Printf("error during token refresh: %s: %s", st.Code(), st.Message())
		return nil, status.Error(codes.Internal, st.Message())
	}

	// TODO: refresh token logic

	return &api.TokenResponse{
		AccessToken:  newToken,
		RefreshToken: "todo;",
	}, nil
}
