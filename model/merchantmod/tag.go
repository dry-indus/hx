package merchantmod

import "go.mongodb.org/mongo-driver/bson/primitive"

type TagAddRequest struct {
	CommodityID primitive.ObjectID `binding:"required" validate:"required"`
	Name        string             `binding:"required" validate:"required"`
}

type TagAddResponse struct {
}

type TagDelRequest struct {
	Id primitive.ObjectID `binding:"required" validate:"required"`
}

type TagDelResponse struct {
}
