package user

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	Id            primitive.ObjectID `bson:"_id,omitempty"`
	Username      string             `bson:"username,omitempty"`
	Password      string             `bson:"password,omitempty"`
	Roles         []string           `bson:"roles,omitempty"`
	Email         string             `bson:"email,omitempty"`
	Activated     bool               `bson:"activated"`
	ActivatedDate time.Time          `bson:"activated_date,omitempty"`

	AvatarUri string `bson:"avatar_uri,omitempty"`
	Describe  string `bson:"describe,omitempty"`

	CreatedDate      time.Time `bson:"created_date,omitempty"`
	LastModifiedDate time.Time `bson:"last_modified_date,omitempty"`
}

type UserDTO struct {
	Id               string    `json:"id,omitempty"`
	Username         string    `json:"username,omitempty"`
	Roles            []string  `json:"roles,omitempty"`
	Email            string    `json:"email,omitempty"`
	Activated        bool      `json:"activated,omitempty"`
	ActivatedDate    int64     `json:"activatedDate,omitempty"`
	AvatarUri        string    `json:"avatarUri,omitempty"`
	Describe         string    `json:"describe,omitempty"`
	CreatedDate      time.Time `json:"createdDate,omitempty"`
	LastModifiedDate time.Time `json:"lastModifiedDate,omitempty"`
}
