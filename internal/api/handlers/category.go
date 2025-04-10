package handlers

import (
	"net/http"

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

// @Summary CreateCategory
// @Security ApiKeyAuth
// @Tags category
// @Description create category
// @ID create-category
// @Accept  json
// @Produce  json
// @Param input body CategoryBody true "category name"
// @Success 200
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/v1/category [post]
func CreateCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

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
func DeleteCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

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
func GetCategories() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
