package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Token struct {
	tokenId      primitive.ObjectID `json:"token_id,omitempty"`
	accessToken  string             `json:"access_token,omitempty"`
	refreshToken string             `json:"refresh_token,omitempty"`
	expireDate   string             `json:"expire_date,omitempty"`
}
