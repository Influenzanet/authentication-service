package main

import (
	"context"
	"log"
	"time"

	api "github.com/influenzanet/authentication-service/api"
	"github.com/influenzanet/authentication-service/models"
	"github.com/influenzanet/authentication-service/tokens"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *authServiceServer) GenerateTempToken(ctx context.Context, t *api.TempTokenInfo) (*api.TempToken, error) {
	if t == nil || t.Purpose == "" {
		return nil, status.Error(codes.InvalidArgument, "missing argument")
	}

	tempToken := models.TempToken{
		UserID:     t.UserId,
		InstanceID: t.InstanceId,
		Purpose:    t.Purpose,
		Info:       t.Info,
		Expiration: t.Expiration,
	}

	if tempToken.Expiration == 0 {
		tempToken.Expiration = tokens.GetExpirationTime(time.Hour * 24 * 10)
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
	if t == nil || t.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "missing argument")
	}

	tempToken, err := getTempTokenFromDB(t.Token)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if time.Now().Unix() > tempToken.Expiration {
		err = deleteTempTokenDB(tempToken.Token)
		log.Println(err)
		return nil, status.Error(codes.InvalidArgument, "token expired")
	}

	return tempToken.ToAPI(), nil
}

func (s *authServiceServer) GetTempTokens(ctx context.Context, t *api.TempTokenInfo) (*api.TempTokenInfos, error) {
	if t == nil || t.UserId == "" || t.InstanceId == "" {
		return nil, status.Error(codes.InvalidArgument, "missing argument")
	}

	tokens, err := getTempTokenForUserDB(t.InstanceId, t.UserId, t.Purpose)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return tokens.ToAPI(), nil
}

func (s *authServiceServer) DeleteTempToken(ctx context.Context, t *api.TempToken) (*api.Status, error) {
	if t == nil || t.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "missing argument")
	}
	if err := deleteTempTokenDB(t.Token); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &api.Status{
		Status: api.Status_NORMAL,
		Msg:    "deleted",
	}, nil
}

func (s *authServiceServer) PurgeUserTempTokens(ctx context.Context, t *api.TempTokenInfo) (*api.Status, error) {
	if t == nil || t.UserId == "" || t.InstanceId == "" {
		return nil, status.Error(codes.InvalidArgument, "missing argument")
	}
	if err := deleteAllTempTokenForUserDB(t.InstanceId, t.UserId, t.Purpose); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &api.Status{
		Status: api.Status_NORMAL,
		Msg:    "deleted",
	}, nil
}
