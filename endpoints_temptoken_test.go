package main

import (
	"context"
	"testing"

	auth_service "github.com/influenzanet/api/dist/go/auth-service"

	"google.golang.org/grpc/status"
)

func TestGenerateTokenEndpoint(t *testing.T) {
	s := authServiceServer{}

	testTempToken := &auth_service.TempTokenInfo{
		UserId:     "test_user_id",
		InstanceId: "test_instance_id",
		Purpose:    "test_purpose",
		Info:       "test_info",
	}

	t.Run("without payload", func(t *testing.T) {
		resp, err := s.GenerateToken(context.Background(), nil)
		if err == nil {
			t.Errorf("or response: %s", resp)
			return
		}
		if status.Convert(err).Message() != "missing argument" {
			t.Errorf("wrong error: %s", err.Error())
		}
	})

	t.Run("with empty payload", func(t *testing.T) {
		resp, err := s.GenerateToken(context.Background(), &auth_service.TempTokenInfo{})
		if err == nil {
			t.Errorf("or response: %s", resp)
			return
		}
		if status.Convert(err).Message() != "missing argument" {
			t.Errorf("wrong error: %s", err.Error())
		}
	})

	t.Run("with valid TempToken", func(t *testing.T) {
		resp, err := s.GenerateToken(context.Background(), testTempToken)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if resp.Token == "" {
			t.Errorf("wrong response: %s", resp)
		}
	})
}
