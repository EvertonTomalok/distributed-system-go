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
		{Key: "steps", Value: internalMessage.Steps},
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

func (a *Adapter) UpdateEventStep(ctx context.Context, orderId string, step dto.EventSteps) error {

	col, err := a.getCol()

	if err != nil {
		log.Fatal(err)
		return err
	}

	steps := []dto.EventSteps{step}

	match := bson.M{"order_id": orderId}
	change := bson.M{"$push": bson.M{"steps": bson.M{"$each": steps}}}

	result, err := col.UpdateOne(ctx, match, change)

	log.Printf("%+v | %+v\n\n", result, err)
	if err != nil {
		return err
	}

	return nil
}
