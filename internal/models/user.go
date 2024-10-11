package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string             `bson:"name" json:"name"`
	Salary      float64            `bson:"salary" json:"salary"`
	Department  string             `bson:"department" json:"department"`
	Country     string             `bson:"country" json:"country"`
	Description string             `bson:"description" json:"description"`
}