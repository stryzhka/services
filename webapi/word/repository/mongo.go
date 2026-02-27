package repository

import (
	"context"
	"log"
	"webapi/models"
	"webapi/word"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoWordRepository struct {
	c *mongo.Client
}

func (m *MongoWordRepository) GetAll(ctx context.Context) []*models.Word {
	coll := m.c.Database("words").Collection("words")
	count, err := coll.CountDocuments(ctx, bson.D{})
	if err != nil {
		return nil
	}
	if count > 1000 {
		stage := bson.D{{"$sample", bson.D{{"size", 50}}}}
		res, err := coll.Aggregate(ctx, mongo.Pipeline{stage}, options.Aggregate())
		if err != nil {
			log.Println(err)
			return nil
		}
		var words []*models.Word
		if err = res.All(context.TODO(), &words); err != nil {
			log.Println(err)
		}
		err = res.Decode(&words)
		if err != nil {
			log.Println(err)
		}
		return words
	}
	cursor, err := coll.Find(ctx, bson.D{})
	if err != nil {
		return nil
	}
	//log.Println(err)
	var words []*models.Word
	if err = cursor.All(ctx, &words); err != nil {
		return nil
	}
	//log.Println(err)
	return words
}

func (m *MongoWordRepository) GetById(ctx context.Context, id string) *models.Word {
	coll := m.c.Database("words").Collection("words")
	sort := bson.D{{"_id", id}}
	res := coll.FindOne(ctx, sort)
	var word *models.Word
	err := res.Decode(&word)
	if err != nil {
		return nil
	}
	//log.Println(err)

	return word
}

func (m *MongoWordRepository) Create(ctx context.Context, word *models.Word) (*models.Word, error) {
	coll := m.c.Database("words").Collection("words")
	_, err := coll.InsertOne(ctx, word)
	if err != nil {
		return nil, err
	}
	return word, err
}

func (m *MongoWordRepository) UpdateById(ctx context.Context, id string, word *models.Word) (*models.Word, error) {
	coll := m.c.Database("words").Collection("words")
	filter := bson.M{"_id": id}
	update := bson.M{"$set": word}
	_, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return word, nil
}

func (m *MongoWordRepository) Delete(ctx context.Context, id string) error {
	coll := m.c.Database("words").Collection("words")
	result, err := coll.DeleteOne(ctx, bson.D{{"_id", id}})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return word.ErrWordNotFound
	}
	return nil
}

func NewMongoWordRepository(c *mongo.Client) *MongoWordRepository {
	return &MongoWordRepository{c}
}
