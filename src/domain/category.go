package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Category struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name string             `json:"name"`
}
