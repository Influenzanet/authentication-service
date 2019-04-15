package main

import (
	auth_api "github.com/influenzanet/api/dist/go/auth-service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TempToken is a database entry for a temporary token
type TempToken struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"token_id,omitempty"`
	Token      string             `bson:"token" json:"token"`
	Expiration int64              `bson:"expiration" json:"expiration"`
	Purpose    string             `bson:"purpose" json:"purpose"`
	UserID     string             `bson:"user_id" json:"user_id"`
	Info       string             `bson:"info" json:"info"`
	InstanceID string             `bson:"instance_id" json:"instance_id"`
}

// ToAPI converts the object from DB to API format
func (t TempToken) ToAPI() *auth_api.TempTokenInfo {
	return &auth_api.TempTokenInfo{
		Token:      t.Token,
		Expiration: t.Expiration,
		Purpose:    t.Purpose,
		UserId:     t.UserID,
		Info:       t.Info,
		InstanceId: t.InstanceID,
	}
}
