package main

import (
	"context"
	"time"

	api "github.com/influenzanet/authentication-service/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *authServiceServer) GenerateTempToken(ctx context.Context, t *api.TempTokenInfo) (*api.TempToken, error) {
	if t == nil || t.UserId == "" || t.InstanceId == "" || t.Purpose == "" {
		return nil, status.Error(codes.InvalidArgument, "missing argument")
	}

	tempToken := TempToken{
		UserID:     t.UserId,
		InstanceID: t.InstanceId,
		Purpose:    t.Purpose,
		Info:       t.Info,
		Expiration: t.Expiration,
	}

	if tempToken.Expiration == 0 {
		tempToken.Expiration = getExpirationTime(time.Hour * 24 * 10)
	}

	token, err := addTempTokenDB(tempToken)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &api.TempToken{
		Token: token,
	}, nil
}

func (s *authServiceServer) ValidateTempToken(ctx context.Context, t *api.TempToken) (*api.TempTokenInfo, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (s *authServiceServer) GetTempTokens(ctx context.Context, t *api.TempTokenInfo) (*api.TempTokenInfos, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (s *authServiceServer) DeleteTempToken(ctx context.Context, t *api.TempToken) (*api.Status, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (s *authServiceServer) PurgeUserTempTokens(ctx context.Context, t *api.TempTokenInfo) (*api.Status, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
