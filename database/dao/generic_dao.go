package dao

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
	var model T
	err := g.DB.First(&model, "id = ?", id).Error
	if err != nil {
		log.Printf("Dao get by ID error: %v", err)
	}
	return &model, err
}

func (g *GenericDAO[T]) Gets(condition any) (*[]T, error) {
	var models []T
	err := g.DB.Where(condition).Find(&models).Error
	if err != nil {
		log.Printf("Dao get multi records error: %v", err)
	}
	return &models, err
}

func (g *GenericDAO[T]) GetsByPage(condition any, order string, offset int, limit int) (*[]T, int64, error) {
	var models []T
	var total int64
	err := g.DB.Order(order).Where(condition).Count(&total).Limit(limit).Offset(offset).Find(&models).Error
	if err != nil {
		log.Printf("Dao get by page error: %v", err)
	}
	return &models, total, err
}

func (g *GenericDAO[T]) Get(condition any) (*T, error) {
	var model T
	err := g.DB.Where(condition).First(&model).Error
	if err != nil {
		log.Printf("Dao get error: %v", err)
	}
	return &model, err
}

func (g *GenericDAO[T]) Exist(condition any) bool {
	var model T
	if err := g.DB.Where(condition).First(&model).Error; err == gorm.ErrRecordNotFound {
		return false
	} else if err == nil {
		return true
	} else {
		log.Printf("Database [%v] check exist error: %v", model, err)
		return false
	}
}

func (g *GenericDAO[T]) ExistByID(id string) bool {
	var model T
	if err := g.DB.First(&model, "id = ?", id).Error; err == gorm.ErrRecordNotFound {
		return false
	} else if err == nil {
		return true
	} else {
		log.Printf("Database [%v] check exist by id error: %v", model, err)
		return false
	}
}

func (g *GenericDAO[T]) GetAll() (*[]T, error) {
	var models []T
	err := g.DB.Find(&models).Error
	if err != nil {
		log.Printf("Dao get all error: %v", err)
	}
	return &models, err
}

func (g *GenericDAO[T]) Create(model *T) error {
	err := g.DB.Create(model).Error
	if err != nil {
		log.Printf("Dao create by ID error: %v", err)
	}
	return err
}

func (g *GenericDAO[T]) Update(model *T) error {
	err := g.DB.Save(model).Error
	if err != nil {
		log.Printf("Dao update by ID error: %v", err)
	}
	return err
}

func (g *GenericDAO[T]) DeleteByID(id string) error {
	var model T
	err := g.DB.Delete(&model, "id = ?", id).Error
	if err != nil {
		log.Printf("Dao delete by ID error: %v", err)
	}
	return err
}

func (g *GenericDAO[T]) Deletes(condition any) error {
	var model T
	err := g.DB.Where(condition).Delete(&model).Error
	return err
}

func (g *GenericDAO[T]) UnscopedDeletes(condition any) error {
	var model T
	err := g.DB.Unscoped().Where(condition).Delete(&model).Error
	return err
}
