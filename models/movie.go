package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Movie is the default data model for the crud app
type Movie struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Name        string             `bson:"name" json:"name"`
	CoverImage  string             `bson:"cover_image" json:"cover_image"`
	Description string             `bson:"description" json:"description"`
}
