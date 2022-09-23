package merchantmod

import "go.mongodb.org/mongo-driver/bson/primitive"

type TagAddRequest struct {
	CommodityID primitive.ObjectID `binding:"required"`
	Name        string             `binding:"required"`
}

type TagAddResponse struct {
}

type TagDelRequest struct {
	Id primitive.ObjectID `binding:"required"`
}

type TagDelResponse struct {
}
