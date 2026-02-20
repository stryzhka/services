package repository

import (
	"context"
	"webapi/category"
	"webapi/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoCategoryRepository struct {
	c *mongo.Client
}

func (m *MongoCategoryRepository) GetAll(ctx context.Context) []*models.Category {
	coll := m.c.Database("words").Collection("categories")
	cursor, err := coll.Find(ctx, bson.D{})
	if err != nil {
		return nil
	}
	var categories []*models.Category
	if err = cursor.All(ctx, &categories); err != nil {
		return nil
	}
	return categories
}

func (m *MongoCategoryRepository) GetById(ctx context.Context, id string) *models.Category {
	coll := m.c.Database("words").Collection("categories")
	filter := bson.D{bson.E{Key: "_id", Value: id}}
	res := coll.FindOne(ctx, filter)
	var cat *models.Category
	err := res.Decode(&cat)
	if err != nil {
		return nil
	}
	return cat
}

func (m *MongoCategoryRepository) GetAllWords(ctx context.Context, categoryId string) []*models.Word {
	coll := m.c.Database("words").Collection("words")
	filter := bson.D{bson.E{Key: "category_id", Value: categoryId}}
	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		return nil
	}
	var words []*models.Word
	if err = cursor.All(ctx, &words); err != nil {
		return nil
	}
	return words
}

func (m *MongoCategoryRepository) Create(ctx context.Context, category *models.Category) (*models.Category, error) {
	coll := m.c.Database("words").Collection("categories")
	_, err := coll.InsertOne(ctx, category)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (m *MongoCategoryRepository) UpdateById(ctx context.Context, id string, category *models.Category) (*models.Category, error) {
	coll := m.c.Database("words").Collection("categories")
	filter := bson.M{"_id": id}
	update := bson.M{"$set": category}
	_, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (m *MongoCategoryRepository) Delete(ctx context.Context, id string) error {
	coll := m.c.Database("words").Collection("categories")
	result, err := coll.DeleteOne(ctx, bson.D{bson.E{Key: "_id", Value: id}})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return category.ErrCategoryNotFound
	}
	return nil
}

func NewMongoCategoryRepository(c *mongo.Client) *MongoCategoryRepository {
	return &MongoCategoryRepository{c}
}
