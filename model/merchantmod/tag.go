package merchantmod

import "go.mongodb.org/mongo-driver/bson/primitive"

type TagAddRequest struct {
	Name string `json:"name" binding:"required" validate:"required"`
}

type TagAddResponse struct {
	Id primitive.ObjectID `json:"id"`
}

type TagDelRequest struct {
	Id primitive.ObjectID `binding:"required" validate:"required"`
}

type TagDelResponse struct {
	Id primitive.ObjectID `json:"id"`
}

type TagStatRequest struct {
	Id primitive.ObjectID `binding:"required" validate:"required"`
}

type TagStatResponse struct {
	Tag        *Tag         `json:"tag"`
	Commoditys []*Commodity `json:"commodity"`
}
