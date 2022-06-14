package mongodb

import (
	"context"
	"fmt"

	"github.com/evertontomalok/distributed-system-go/internal/domain/core/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var MongoDatabase *mongo.Database

func New() *Adapter {
	return &Adapter{Db: MongoDatabase}
}

type Adapter struct {
	Db *mongo.Database
}

func (a *Adapter) SaveEvent(ctx context.Context, internalMessage dto.BrokerInternalMessage) error {
	coll := a.Db.Collection("events")
	doc := bson.D{
		{Key: "order_id", Value: internalMessage.ID},
		{Key: "value", Value: internalMessage.Value},
		{Key: "method", Value: internalMessage.Method},
		{Key: "installments", Value: internalMessage.Installments},
		{Key: "user_id", Value: internalMessage.UserId},
		{Key: "status", Value: false},
		{Key: "metadata", Value: nil},
	}
	result, err := coll.InsertOne(ctx, doc)

	if err != nil {
		return err
	}
	fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)

	return nil
}
