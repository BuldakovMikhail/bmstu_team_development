package adapters

import (
	"context"
	"todolist/internal/api/handlers"
	"todolist/internal/models"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type CategoryRepository interface {
	CreateCategory(ctx context.Context, body *models.CategoryBody) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetAll(ctx context.Context, pageIndex, recordsPerPage int) ([]models.Category, error)
}

type CategoryAdapter struct {
	repository CategoryRepository
}

func NewCategoryAdapter(repository CategoryRepository) handlers.CategoriesProvider {
	return &CategoryAdapter{repository: repository}
}

func (c *CategoryAdapter) CreateCategory(ctx context.Context, body *models.CategoryBody) error {
	err := c.repository.CreateCategory(ctx, body)
	if err != nil {
		return errors.Wrap(err, "failed to create category")
	}

	return nil
}

func (c *CategoryAdapter) Delete(ctx context.Context, id uuid.UUID) error {
	err := c.repository.Delete(ctx, id)
	if err != nil {
		return errors.Wrap(err, "failed to delete category")
	}
	return nil
}

func (c *CategoryAdapter) GetAll(ctx context.Context, pageIndex, recordsPerPage int) ([]models.Category, error) {
	categories, err := c.repository.GetAll(ctx, pageIndex, recordsPerPage)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all categories")
	}
	return categories, nil
}
