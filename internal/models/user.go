package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	Id        primitive.ObjectID `json:"id,omitempty"`
	Username  string             `json:"username,omitempty" unique:"true"`
	Email     string             `json:"email,omitempty" validate:"required"`
	FirstName string             `json:"first_name,omitempty" validate:"required"`
	LastName  string             `json:"last_name,omitempty" validate:"required"`
	Password  string             `json:"password,omitempty" validate:"required"`
	Age       int                `json:"age,omitempty"`
	Gender    bool               `json:"gender"`
	Phone     string             `json:"phone,omitempty"`
	IsActive  bool               `json:"is_active"`
	Role      string             `json:"role,omitempty" validate:"required"`
	CreateAt  time.Time          `json:"create_at,omitempty"`
	UpdateAt  time.Time          `json:"update_at,omitempty"`
}
