package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Comment struct {
    Id            primitive.ObjectID `bson:"_id" json:"id,omitempty"`
    ContentId     string             `json:"contentId,omitempty" validate:"omitempty"`
    UserId	  string             `json:"userId,omitempty" validate:"omitempty"`
    Comment	  string             `json:"comment,omitempty" validate:"required,omitempty"`
}
