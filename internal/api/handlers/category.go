package handlers

import (
	"context"
	"net/http"
	"time"
	"todolist/internal/models"
	"todolist/internal/pkg/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type CategoryBody struct {
	Name string `json:"name"`
}

type CategoryResponse struct {
	ID uuid.UUID `json:"id"`
	CategoryBody
}

type CategoriesResponse struct {
	Categories []CategoryResponse `json:"categories"`
}

type CategoriesProvider interface {
	CreateCategory(ctx context.Context, category *models.CategoryBody) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetAll(ctx context.Context, pageIndex, recordsPerPage int) ([]models.Category, error)
}

// @Summary CreateCategory
// @Security ApiKeyAuth
// @Tags category
// @Description create category
// @ID create-category
// @Accept  json
// @Produce  json
// @Param input body CategoryBody true "category name"
// @Success 200
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/v1/category [post]
func CreateCategory(categoryProvider CategoriesProvider, timeout time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CategoryBody
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}
		category := models.CategoryBody{Name: req.Name}

		ctx := r.Context()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		err = categoryProvider.CreateCategory(ctx, &category)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}
		render.Status(r, http.StatusOK)
	}
}

// @Summary DeleteCategory
// @Security ApiKeyAuth
// @Tags category
// @Description delete category
// @ID delete-category
// @Accept  json
// @Produce  json
// @Param id   path      string  true  "Category ID (UUID)"
// @Success 200
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/v1/category/{id} [delete]
func DeleteCategory(categoryProvider CategoriesProvider, timeout time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		uuid, err := uuid.Parse(id)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid UUID"))
			return
		}

		ctx := r.Context()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		err = categoryProvider.Delete(ctx, uuid)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}
		render.Status(r, http.StatusOK)
	}
}

// @Summary GetCategories
// @Security ApiKeyAuth
// @Tags category
// @Description get all categories
// @ID get-categories
// @Accept  json
// @Produce  json
// @Param input body Pagination true "pagination info"
// @Success 200 {object} CategoriesResponse
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/v1/category/all [post]
func GetCategories(categoryProvider CategoriesProvider, timeout time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Pagination
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		ctx := r.Context()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var categories []models.Category
		categories, err = categoryProvider.GetAll(ctx, req.PageIndex, req.RecordsPerPage)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, toCategoriesResponse(categories))
	}
}

func toCategoryResponse(category models.Category) CategoryResponse {
	return CategoryResponse{
		ID: category.ID,
		CategoryBody: CategoryBody{
			Name: category.Name,
		},
	}
}

func toCategoriesResponse(categories []models.Category) CategoriesResponse {
	responses := make([]CategoryResponse, 0, len(categories))
	for _, category := range categories {
		responses = append(responses, toCategoryResponse(category))
	}
	return CategoriesResponse{
		Categories: responses,
	}
}
