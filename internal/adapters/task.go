package adapters

import (
	"context"
	"todolist/internal/models"

	"github.com/pkg/errors"

	"github.com/google/uuid"
)

type TaskRepository interface {
	CreateTask(ctx context.Context, body *models.TaskBody) error
	Update(ctx context.Context, id uuid.UUID, body *models.TaskBody) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.TaskFullInfo, error)
	GetAll(ctx context.Context, pageIndex, recordsPerPage int) ([]models.TaskShortInfo, error)
	Delete(ctx context.Context, id uuid.UUID) error
	ToggleDone(ctx context.Context, id uuid.UUID) error
}

type TaskAdapter struct {
	repository TaskRepository
}

func NewTaskAdapter(repository TaskRepository) *TaskAdapter {
	return &TaskAdapter{repository: repository}
}

func (t *TaskAdapter) CreateTask(ctx context.Context, body *models.TaskBody) error {
	err := t.repository.CreateTask(ctx, body)
	if err != nil {
		return errors.Wrap(err, "failed to create task")
	}

	return nil
}

func (t *TaskAdapter) Update(ctx context.Context, id uuid.UUID, body *models.TaskBody) error {
	err := t.repository.Update(ctx, id, body)
	if err != nil {
		return errors.Wrapf(err, "failed to update task with id: %s", id)
	}
	return nil
}

func (t *TaskAdapter) GetByID(ctx context.Context, id uuid.UUID) (*models.TaskFullInfo, error) {
	task, err := t.repository.GetByID(ctx, id)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get task by id: %s", id)
	}
	return task, nil
}

func (t *TaskAdapter) GetAll(ctx context.Context, pageIndex, recordsPerPage int) ([]models.TaskShortInfo, error) {
	tasks, err := t.repository.GetAll(ctx, pageIndex, recordsPerPage)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all tasks")
	}
	return tasks, nil
}

func (t *TaskAdapter) Delete(ctx context.Context, id uuid.UUID) error {
	err := t.repository.Delete(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "failed to delete task with id: %s", id)
	}
	return nil
}

func (t *TaskAdapter) ToggleDone(ctx context.Context, id uuid.UUID) error {
	err := t.repository.ToggleDone(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "failed to toggle task done status with id: %s", id)
	}
	return nil
}
