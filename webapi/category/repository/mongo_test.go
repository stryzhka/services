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

func initMongoCategoryTest() (*mongo.Client, error) {
	connectionString := "mongodb://admin:admin@localhost:27017/test?authSource=admin"
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(connectionString).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(opts)
	if err != nil {
		panic(err)
	}
	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{bson.E{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		return nil, err
	}
	log.Println("mongo connected")
	return client, nil
}

func TestMongoCategoryRepository_Create(t *testing.T) {
	client, _ := initMongoCategoryTest()
	repo := NewMongoCategoryRepository(client)
	cat := &models.Category{
		Id:   uuid.New().String(),
		Name: "test-category",
	}
	expCat, err := repo.Create(context.TODO(), cat)
	assert.NoError(t, err)
	assert.Equal(t, cat, expCat)
}

func TestMongoCategoryRepository_GetAll(t *testing.T) {
	client, _ := initMongoCategoryTest()
	repo := NewMongoCategoryRepository(client)
	categories := repo.GetAll(context.TODO())
	log.Println(categories)
}

func TestMongoCategoryRepository_GetById(t *testing.T) {
	client, _ := initMongoCategoryTest()
	repo := NewMongoCategoryRepository(client)
	cat := repo.GetById(context.TODO(), "5230165e-363a-43ce-8717-65b0f5b79e60")
	log.Println(cat)
}

func TestMongoCategoryRepository_Delete(t *testing.T) {
	client, _ := initMongoCategoryTest()
	repo := NewMongoCategoryRepository(client)
	err := repo.Delete(context.TODO(), "sdfjna")
	assert.NoError(t, err)
}
