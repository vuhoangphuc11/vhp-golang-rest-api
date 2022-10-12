package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	Id         primitive.ObjectID `json:"id,omitempty"`
	Email      string             `json:"email,omitempty" validate:"required"`
	FullName   string             `json:"fullname,omitempty" validate:"required"`
	Password   string             `json:"password,omitempty" validate:"required"`
	Age        int                `json:"age,omitempty" validate:"required"`
	Gender     bool               `json:"gender,omitempty" validate:"required"`
	Phone      string             `json:"phone,omitempty" validate:"required"`
	IsActive   bool               `json:"active,omitempty" validate:"required"`
	Role       string             `json:"role,omitempty" validate:"required"`
	CreateDate time.Time          `json:"create_date,omitempty"`
	UpdateDate time.Time          `json:"update_date,omitempty"`
}
