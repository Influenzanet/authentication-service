package main

import (
	"context"

	influenzanet "github.com/influenzanet/api/dist/go"
	auth_api "github.com/influenzanet/api/dist/go/auth-service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *authServiceServer) GenerateToken(ctx context.Context, t *auth_api.TempTokenInfo) (*auth_api.TempToken, error) {
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
		tempToken.Expiration = getExpirationTime(10)
	}

	token, err := dbCreateToken(tempToken)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &auth_api.TempToken{
		Token: token,
	}, nil
}

func (s *authServiceServer) ValdateToken(ctx context.Context, t *auth_api.TempToken) (*auth_api.TempTokenInfo, error) {
	return nil, nil
}
func (s *authServiceServer) GetTokensByUser(ctx context.Context, t *auth_api.TempTokenInfo) (*auth_api.TempTokenInfos, error) {
	return nil, nil
}
func (s *authServiceServer) DeleteToken(ctx context.Context, t *auth_api.TempToken) (*influenzanet.Status, error) {
	return nil, nil
}
func (s *authServiceServer) PurgeUserTokens(ctx context.Context, t *auth_api.TempTokenInfo) (*influenzanet.Status, error) {
	return nil, nil
}
