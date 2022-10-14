package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Token struct {
	tokenId      primitive.ObjectID `json:"tokenId,omitempty"`
	accessToken  string             `json:"accessToken,omitempty"`
	refreshToken string             `json:"refreshToken,omitempty"`
	expireDate   string             `json:"expireDate,omitempty"`
}
