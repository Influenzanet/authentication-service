package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc/codes"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/status"

	influenzanet "github.com/influenzanet/api/dist/go"
	auth_api "github.com/influenzanet/api/dist/go/auth-service"
	user_api "github.com/influenzanet/api/dist/go/user-management"
)

func checkTokenAgeMaturity(issuedAt int64) bool {
	return time.Now().Unix() < time.Unix(issuedAt, 0).Add(minTokenAge).Unix()
}

func (s *authServiceServer) Status(ctx context.Context, _ *empty.Empty) (*influenzanet.Status, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (s *authServiceServer) LoginWithEmail(ctx context.Context, req *influenzanet.UserCredentials) (*auth_api.EncodedToken, error) {
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

	return &auth_api.EncodedToken{
		Token: token,
	}, nil
}

func (s *authServiceServer) SignupWithEmail(ctx context.Context, req *user_api.NewUser) (*auth_api.EncodedToken, error) {
	if req == nil || req.Auth == nil || req.Profile == nil {
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

	return &auth_api.EncodedToken{
		Token: token,
	}, nil
}

func (s *authServiceServer) ValidateJWT(ctx context.Context, req *auth_api.EncodedToken) (*influenzanet.ParsedToken, error) {
	if req == nil || req.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "missing arguments")
	}
	// Parse and validate token
	parsedToken, ok, oldKey, err := validateToken(req.Token)
	if err != nil || !ok {
		return nil, status.Error(codes.InvalidArgument, "invalid token")
	}

	// TODO: what to do when token is generated by old key
	if oldKey {
		log.Println("handling old keys for token validation is not implemented")
	}

	return &influenzanet.ParsedToken{
		UserId:     parsedToken.UserID,
		Roles:      parsedToken.Roles,
		InstanceId: parsedToken.InstanceID,
		IssuedAt:   parsedToken.IssuedAt,
	}, nil
}

func (s *authServiceServer) RenewJWT(ctx context.Context, req *auth_api.EncodedToken) (*auth_api.EncodedToken, error) {
	if req == nil || req.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "missing arguments")
	}

	// Parse and validate token
	parsedToken, ok, oldKey, err := validateToken(req.Token)
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

	// TODO: what to do when token is generated by old key
	if oldKey {
		log.Println("handling old keys for token renewal is not implemented")
	}

	// Generate new token:
	newToken, err := generateNewToken(parsedToken.UserID, parsedToken.Roles, parsedToken.InstanceID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &auth_api.EncodedToken{
		Token: newToken,
	}, nil
}
