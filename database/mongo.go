package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDatabase representa a estrutura para a conexão com o MongoDB
type MongoDatabase struct {
	Client     *mongo.Client
	Database   *mongo.Database
	Collection *mongo.Collection
}

// NewMongoDB cria uma nova instância de conexão com o MongoDB
func NewMongoDB(uri, dbName, collectionName string) (*MongoDatabase, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	database := client.Database(dbName)
	collection := database.Collection(collectionName)

	return &MongoDatabase{
		Client:     client,
		Database:   database,
		Collection: collection,
	}, nil
}

// Close fecha a conexão com o MongoDB
func (mdb *MongoDatabase) Close() {
	err := mdb.Client.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
