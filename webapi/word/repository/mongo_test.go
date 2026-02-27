package repository

import (
	"context"
	"log"
	"testing"
	"webapi/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func initMongo() (*mongo.Client, error) {
	connectionString := "mongodb://admin:admin@localhost:27017/test?authSource=admin"
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	// Defines the options for the MongoDB client
	opts := options.Client().ApplyURI(connectionString).SetServerAPIOptions(serverAPI)
	// Creates a new client and connects to the server
	client, err := mongo.Connect(opts)
	if err != nil {
		panic(err)
	}
	// Sends a ping to confirm a successful connection
	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		return nil, err
	}
	log.Println("mongo connected")
	return client, nil
}

func TestMongoWordRepository_Create(t *testing.T) {
	client, _ := initMongo()
	wordRepo := NewMongoWordRepository(client)
	word := &models.Word{
		Id:            uuid.New().String(),
		EngText:       "test",
		NativeText:    "test",
		Transcription: "test",
		Difficulty:    "test",
		CategoryId:    "test",
	}
	expWord, err := wordRepo.Create(context.TODO(), word)
	assert.NoError(t, err)
	assert.Equal(t, word, expWord)
}

func TestMongoWordRepository_GetAll(t *testing.T) {
	client, _ := initMongo()
	wordRepo := NewMongoWordRepository(client)
	words := wordRepo.GetAll(context.TODO())
	log.Println(words)
}

func TestMongoWordRepository_GetById(t *testing.T) {
	client, _ := initMongo()
	wordRepo := NewMongoWordRepository(client)
	word := wordRepo.GetById(context.TODO(), "5230165e-363a-43ce-8717-65b0f5b79e60")
	log.Println(word)
}

func TestMongoWordRepository_Delete(t *testing.T) {
	client, _ := initMongo()
	wordRepo := NewMongoWordRepository(client)
	err := wordRepo.Delete(context.TODO(), "sdfjna")
	assert.NoError(t, err)
}
