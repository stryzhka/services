package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
	"webapi/models"
	"webapi/word"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoWordRepository struct {
	c     *mongo.Client
	redis word.Cache
}

func (m *MongoWordRepository) GetAll(ctx context.Context, filter word.WordFilter) []*models.Word {
	var key string
	key = fmt.Sprintf("word:%s", filter.ToBSON().String())
	cacheRes, err := m.redis.GetMany(ctx, key)
	if err != nil {
		log.Println(err)
		//return nil
	}
	if cacheRes != nil {
		var res []*models.Word
		err = json.Unmarshal(cacheRes, &res)
		if err != nil {
			log.Println(err)
			return nil
		}
		return res
	}
	if err != nil {
		log.Println(err)
	}
	coll := m.c.Database("words").Collection("words")
	count, err := coll.CountDocuments(ctx, filter.ToBSON())
	opts := options.Find().SetSort(bson.D{})

	//log.Println(filter.ToBSON().String())
	if err != nil {
		return nil
	}
	if count > 1000 {

		matchStage := bson.D{{"$match", filter.ToBSON()}}
		sampleStage := bson.D{{"$sample", bson.D{{"size", 50}}}}
		res, err := coll.Aggregate(ctx, mongo.Pipeline{matchStage, sampleStage})

		if err != nil {
			log.Println(err)
			return nil
		}
		var words []*models.Word
		if err = res.All(context.TODO(), &words); err != nil {
			log.Println(err)
		}
		_, err = m.redis.Set(ctx, key, 1000*time.Second, words)
		return words
	}

	cursor, err := coll.Find(ctx, filter.ToBSON(), opts)
	if err != nil {
		return nil
	}
	//log.Println(err)
	var words []*models.Word
	if err = cursor.All(ctx, &words); err != nil {
		return nil
	}
	_, err = m.redis.Set(ctx, key, 1000*time.Second, words)
	return words
}

func (m *MongoWordRepository) GetById(ctx context.Context, id string) *models.Word {
	key := fmt.Sprintf("word:%s", id)
	cacheRes, err := m.redis.GetOne(ctx, key)
	if err != nil {
		log.Println(err)
		//return nil
	}
	if cacheRes != nil {
		return cacheRes.(*models.Word)
	}
	coll := m.c.Database("words").Collection("words")
	sort := bson.D{{"_id", id}}
	res := coll.FindOne(ctx, sort)
	var word *models.Word
	err = res.Decode(&word)
	if err != nil {
		return nil
	}
	_, err = m.redis.Set(ctx, key, 20*time.Second, word)
	if err != nil {
		log.Println(err)
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
	err = m.redis.InvalidateAll(ctx)
	if err != nil {
		log.Println(err)
	}
	return word, err
}

func (m *MongoWordRepository) UpdateById(ctx context.Context, id string, word *models.Word) (*models.Word, error) {
	coll := m.c.Database("words").Collection("words")
	filter := bson.M{"_id": id}
	update := bson.M{"$set": word}
	_, err := coll.UpdateOne(ctx, filter, update)
	key := fmt.Sprintf("word:%s", id)
	_, err = m.redis.Set(ctx, key, 20*time.Second, word)
	if err != nil {
		log.Println(err)
	}
	err = m.redis.InvalidateAll(ctx)
	if err != nil {
		log.Println(err)
	}
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
	err = m.redis.InvalidateAll(ctx)
	if err != nil {
		log.Println(err)
	}
	return nil
}

func (m *MongoWordRepository) Confirm(ctx context.Context, id, confirmedAt string) error {
	coll := m.c.Database("words").Collection("words")

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"confirmed_status": "confirmed", "confirmed_at": confirmedAt}}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var word models.Word
	err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&word)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("word:%s", id)
	_, err = m.redis.Set(ctx, key, 20*time.Second, &word)
	if err != nil {
		log.Println(err)
	}

	err = m.redis.InvalidateAll(ctx)
	if err != nil {
		log.Println(err)
	}

	return nil
}

func NewMongoWordRepository(c *mongo.Client, redis word.Cache) *MongoWordRepository {
	return &MongoWordRepository{
		c:     c,
		redis: redis,
	}
}
