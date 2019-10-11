package main

import (
	"context"
	"testing"
	"time"

	api "github.com/influenzanet/authentication-service/api"
	"google.golang.org/grpc/status"
)

func TestGenerateTempTokenEndpoint(t *testing.T) {
	s := authServiceServer{}

	testTempToken := &api.TempTokenInfo{
		UserId:     "test_user_id",
		InstanceId: testInstanceID,
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
		resp, err := s.GenerateTempToken(context.Background(), &api.TempTokenInfo{})
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
	s := authServiceServer{}

	testTempToken := TempToken{
		UserID:     "test_user_id",
		InstanceID: testInstanceID,
		Purpose:    "test_purpose_validation",
		Info:       "test_info",
		Expiration: getExpirationTime(10 * time.Second),
	}
	token, err := addTempTokenDB(testTempToken)
	if err != nil {
		t.Error(err)
		return
	}
	testTempToken.Token = token

	t.Run("without payload", func(t *testing.T) {
		resp, err := s.ValidateTempToken(context.Background(), nil)
		if err == nil {
			t.Errorf("or response: %s", resp)
			return
		}
		if status.Convert(err).Message() != "missing argument" {
			t.Errorf("wrong error: %s", err.Error())
		}
	})

	t.Run("with empty payload", func(t *testing.T) {
		resp, err := s.ValidateTempToken(context.Background(), &api.TempToken{})
		if err == nil {
			t.Errorf("or response: %s", resp)
			return
		}
		if status.Convert(err).Message() != "missing argument" {
			t.Errorf("wrong error: %s", err.Error())
		}
	})

	t.Run("with not existing token", func(t *testing.T) {
		resp, err := s.ValidateTempToken(context.Background(), &api.TempToken{
			Token: testTempToken.Token + "1",
		})
		if err == nil {
			t.Errorf("or response: %s", resp)
			return
		}
		if status.Convert(err).Message() != "mongo: no documents in result" {
			t.Errorf("wrong error: %s", err.Error())
		}
	})

	t.Run("with valid payload", func(t *testing.T) {
		resp, err := s.ValidateTempToken(context.Background(), &api.TempToken{
			Token: testTempToken.Token,
		})
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if resp.UserId != testTempToken.UserID || resp.InstanceId != testTempToken.InstanceID || resp.Purpose != testTempToken.Purpose || resp.Info != testTempToken.Info || resp.Expiration != testTempToken.Expiration {
			t.Error(resp)
			t.Error("wrong token infos")
			return
		}
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
