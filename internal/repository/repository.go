package repository

import (
	"context"
)

type Repository[T any] interface {
	// Get all the items
	GetAll() ([]T, error)
	// Get an item by id
	GetById(id string) (T, error)
	// Create an item
	Create(item T) (T, error)
	// Update an item
	Update(id string, item T) (T, error)
	// Delete an item
	Delete(id string) error
}

func NewRepository[T any](ctx context.Context) Repository[T] {
	return &repositoryImpl[T]{
		ctx: ctx,
	}
}

type repositoryImpl[T any] struct {
	ctx   context.Context
	model T
}

func (r *repositoryImpl[T]) GetAll() ([]T, error) {
	// Implementation here
	return nil, nil
}

func (r *repositoryImpl[T]) GetById(id string) (T, error) {
	// Implementation here
	var result T
	return result, nil
}

func (r *repositoryImpl[T]) Create(item T) (T, error) {
	// Implementation here
	return item, nil
}

func (r *repositoryImpl[T]) Update(id string, item T) (T, error) {
	// Implementation here
	return item, nil
}

func (r *repositoryImpl[T]) Delete(id string) error {
	// Implementation here
	return nil
}
