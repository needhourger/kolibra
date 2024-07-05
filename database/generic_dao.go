package database

import (
	"log"

	"gorm.io/gorm"
)

type GenericDAO[T any] struct {
	DB *gorm.DB
}

func NewGenericDAO[T any](db *gorm.DB) *GenericDAO[T] {
	return &GenericDAO[T]{DB: db}
}

func (g *GenericDAO[T]) GetByID(id string) (*T, error) {
	var model *T
	err := g.DB.First(model, id).Error
	return model, err
}

func (g *GenericDAO[T]) Gets(condition any) (*[]T, error) {
	var models *[]T
	err := g.DB.Where(condition).Find(models).Error
	return models, err
}

func (g *GenericDAO[T]) Get(condition any) (*T, error) {
	var model *T
	err := g.DB.Where(condition).First(model).Error
	return model, err
}

func (g *GenericDAO[T]) Exist(condition any) bool {
	var model *T
	if err := g.DB.Where(condition).First(model).Error; err == gorm.ErrRecordNotFound {
		return false
	} else if err == nil {
		return true
	} else {
		log.Printf("Database [%v] check exist error: %v", model, err)
		return false
	}
}

func (g *GenericDAO[T]) GetAll() (*[]T, error) {
	var models *[]T
	err := g.DB.Find(models).Error
	g.DB.Offset(5).Limit(20).Find(models)
	return models, err
}

func (g *GenericDAO[T]) Create(model *T) error {
	return g.DB.Create(model).Error
}

func (g *GenericDAO[T]) Update(model *T) error {
	return g.DB.Save(model).Error
}

func (g *GenericDAO[T]) DeleteByID(id string) error {
	var model *T
	return g.DB.Delete(model, id).Error
}
