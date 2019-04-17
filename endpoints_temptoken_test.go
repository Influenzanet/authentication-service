package main

import (
	"context"
	"testing"

	auth_service "github.com/influenzanet/api/dist/go/auth-service"

	"google.golang.org/grpc/status"
)

func TestGenerateTempTokenEndpoint(t *testing.T) {
	s := authServiceServer{}

	testTempToken := &auth_service.TempTokenInfo{
		UserId:     "test_user_id",
		InstanceId: "test_instance_id",
		Purpose:    "test_purpose",
		Info:       "test_info",
	}

	t.Run("without payload", func(t *testing.T) {
		resp, err := s.GenerateTempToken(context.Background(), nil)
		if err == nil {
			t.Errorf("or response: %s", resp)
			return
		}
		if status.Convert(err).Message() != "missing argument" {
			t.Errorf("wrong error: %s", err.Error())
		}
	})

	t.Run("with empty payload", func(t *testing.T) {
		resp, err := s.GenerateTempToken(context.Background(), &auth_service.TempTokenInfo{})
		if err == nil {
			t.Errorf("or response: %s", resp)
			return
		}
		if status.Convert(err).Message() != "missing argument" {
			t.Errorf("wrong error: %s", err.Error())
		}
	})

	t.Run("with valid TempToken", func(t *testing.T) {
		resp, err := s.GenerateTempToken(context.Background(), testTempToken)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if resp.Token == "" {
			t.Errorf("wrong response: %s", resp)
		}
	})
}

func TestValidateTempTokenEndpoint(t *testing.T) {
	// TODO: create test temp token

	t.Run("without payload", func(t *testing.T) {
		t.Error("test not implemented")
	})

	t.Run("with empty payload", func(t *testing.T) {
		t.Error("test not implemented")
	})

	t.Run("with not existing token", func(t *testing.T) {
		t.Error("test not implemented")
	})

	t.Run("with valid payload", func(t *testing.T) {
		t.Error("test not implemented")
	})
}

func TestGetTempTokensEndpoint(t *testing.T) {
	// TODO: create test tokens
	t.Run("without payload", func(t *testing.T) {
		t.Error("test not implemented")
	})

	t.Run("with empty payload", func(t *testing.T) {
		t.Error("test not implemented")
	})

	t.Run("get by user_id + instace_id", func(t *testing.T) {
		t.Error("test not implemented")
	})

	t.Run("get by user_id + instace_id + type", func(t *testing.T) {
		t.Error("test not implemented")
	})

	t.Run("get by type", func(t *testing.T) {
		t.Error("test not implemented")
	})
}

func TestDeleteTempTokenEndpoint(t *testing.T) {
	t.Run("without payload", func(t *testing.T) {
		t.Error("test not implemented")
	})

	t.Run("with empty payload", func(t *testing.T) {
		t.Error("test not implemented")
	})

	t.Run("with not existing token", func(t *testing.T) {
		t.Error("test not implemented")
	})

	t.Run("with existing token", func(t *testing.T) {
		t.Error("test not implemented")
	})
}

func TestPurgeUserTempTokensEndpoint(t *testing.T) {
	t.Run("without payload", func(t *testing.T) {
		t.Error("test not implemented")
	})

	t.Run("with empty payload", func(t *testing.T) {
		t.Error("test not implemented")
	})

	t.Run("with not exisiting user_id/instance_id combination", func(t *testing.T) {
		t.Error("test not implemented")
	})

	t.Run("with exisiting user_id/instance_id combination", func(t *testing.T) {
		t.Error("test not implemented")
	})
}
