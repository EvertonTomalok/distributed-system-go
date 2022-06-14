package mongodb

import (
	"context"
	"fmt"
	"log"

	application "github.com/evertontomalok/distributed-system-go/internal/app"
	"github.com/evertontomalok/distributed-system-go/internal/domain/core/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func New(config application.Config) *Adapter {
	return &Adapter{Host: config.Mongodb.Host}
}

type Adapter struct {
	Host string
}

func (a *Adapter) getCol() (*mongo.Collection, error) {
	var err error
	session, err := mongo.NewClient(options.Client().ApplyURI(a.Host))
	if err != nil {
		log.Fatal(err)
		return &mongo.Collection{}, err
	}
	session.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
		return &mongo.Collection{}, err
	}

	var DB = session.Database("event_source")
	return DB.Collection("events"), nil
}

func (a *Adapter) SaveEvent(ctx context.Context, internalMessage dto.BrokerInternalMessage) error {

	value, err := primitive.ParseDecimal128(internalMessage.Value.String())
	if err != nil {
		return err
	}

	doc := bson.D{
		{Key: "order_id", Value: internalMessage.ID},
		{Key: "value", Value: value},
		{Key: "method", Value: internalMessage.Method},
		{Key: "installments", Value: internalMessage.Installments},
		{Key: "user_id", Value: internalMessage.UserId},
		{Key: "status", Value: false},
		{Key: "metadata", Value: nil},
	}

	col, err := a.getCol()

	if err != nil {
		log.Fatal(err)
		return err
	}

	result, err := col.InsertOne(ctx, doc)

	if err != nil {
		return err
	}
	fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)

	return nil
}
