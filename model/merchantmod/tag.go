package merchantmod

import "go.mongodb.org/mongo-driver/bson/primitive"

type TagAddRequest struct {
	CommodityID primitive.ObjectID `json:"commodityId" binding:"required" validate:"required"`
	Name        string             `json:"name" binding:"required" validate:"required"`
}

type TagAddResponse struct {
	Id primitive.ObjectID `json:"id"`
}

type TagDelRequest struct {
	Id primitive.ObjectID `binding:"required" validate:"required"`
}

type TagDelResponse struct {
}
