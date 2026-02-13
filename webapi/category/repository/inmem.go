package repository

import (
	"context"
	"log"
	category2 "webapi/category"
	"webapi/models"

	"github.com/hashicorp/go-memdb"
)

type InmemRepository struct {
	db *memdb.MemDB
}

func NewInmemRepository(db *memdb.MemDB) *InmemRepository {
	return &InmemRepository{db: db}
}

func (m *InmemRepository) GetAll(ctx context.Context) []*models.Category {
	var categories []*models.Category
	txn := m.db.Txn(false)
	defer txn.Abort()
	it, err := txn.Get("category", "id")
	if err != nil {
		return nil
	}
	for obj := it.Next(); obj != nil; obj = it.Next() {
		c := obj.(*models.Category)
		categories = append(categories, c)
	}
	txn.Commit()
	return categories
}

func (m *InmemRepository) GetById(ctx context.Context, id string) *models.Category {
	txn := m.db.Txn(false)
	defer txn.Abort()
	found, err := txn.First("category", "id", id)
	log.Println(found)
	if err != nil {
		log.Println(err)
		return nil
	}
	if found == nil {
		return nil
	}
	return found.(*models.Category)
}

func (m *InmemRepository) Create(ctx context.Context, category *models.Category) (*models.Category, error) {
	txn := m.db.Txn(true)
	defer txn.Abort()

	existing, err := txn.First("category", "name", category.Name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, category2.ErrCategoryAlreadyExists
	}

	err = txn.Insert("category", category)
	if err != nil {
		return nil, err
	}

	txn.Commit()
	return category, nil
}

func (m *InmemRepository) UpdateById(ctx context.Context, id string, category *models.Category) (*models.Category, error) {
	txn := m.db.Txn(true)
	defer txn.Abort()
	existing, err := txn.First("category", "id", id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, category2.ErrCategoryNotFound
	}
	err = txn.Insert("category", category)
	if err != nil {
		return nil, err
	}
	txn.Commit()
	return category, nil
}

func (m *InmemRepository) Delete(ctx context.Context, id string) error {
	txn := m.db.Txn(true)
	defer txn.Abort()
	existing, err := txn.First("category", "id", id)
	if err != nil {
		return err
	}
	if existing == nil {
		return category2.ErrCategoryNotFound
	}
	err = txn.Delete("category", existing)
	if err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (m *InmemRepository) GetAllWords(ctx context.Context, categoryId string) []*models.Word {
	var words []*models.Word
	txn := m.db.Txn(false)
	defer txn.Abort()
	it, err := txn.Get("word", "category_id", categoryId)
	if err != nil {
		return nil
	}
	for obj := it.Next(); obj != nil; obj = it.Next() {
		w := obj.(*models.Word)
		words = append(words, w)
	}
	txn.Commit()
	return words
}
