package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	Id        primitive.ObjectID `json:"id,omitempty"`
	Username  string             `json:"username,omitempty"`
	Email     string             `json:"email,omitempty" validate:"required"`
	FirstName string             `json:"firstName,omitempty" validate:"required"`
	LastName  string             `json:"lastName,omitempty" validate:"required"`
	Password  string             `json:"password,omitempty" validate:"required"`
	Age       int                `json:"age,omitempty"`
	Gender    bool               `json:"gender,omitempty"`
	Phone     string             `json:"phone,omitempty"`
	IsActive  bool               `json:"isActive,omitempty"`
	Role      string             `json:"role,omitempty" validate:"required"`
	CreateAt  time.Time          `json:"create_at,omitempty"`
	UpdateAt  time.Time          `json:"update_at,omitempty"`
}
