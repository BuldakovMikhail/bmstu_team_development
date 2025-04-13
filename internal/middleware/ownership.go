package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
	"todolist/internal/adapters"
	"todolist/internal/pkg/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type TaskBody struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type TaskRequest struct {
	TaskBody
	CategoryIds []uuid.UUID `json:"category_ids"`
}

type OwnershipMiddleware struct {
	userService adapters.UserAdapter
	timeout     time.Duration
}

func NewOwnershipMiddleware(service adapters.UserAdapter, timeout time.Duration) OwnershipMiddleware {
	return OwnershipMiddleware{
		userService: service,
		timeout:     timeout,
	}
}

func (m *OwnershipMiddleware) CheckCategoriesMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// saving data for further handlers
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to read request body"))
			return
		}
		r.Body.Close()

		userID, ok := r.Context().Value(UserIDContextKey).(uuid.UUID)
		if !ok {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, response.Error("Missing userID"))
			return
		}

		var req TaskRequest
		err = json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&req)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		ctx := r.Context()
		ctx, cancel := context.WithTimeout(ctx, m.timeout)
		defer cancel()

		areCategoriesOwned, err := m.userService.CheckCategoriesOwnership(ctx, userID, req.CategoryIds)

		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		if !areCategoriesOwned {
			render.Status(r, http.StatusForbidden) // Sets status in context
			render.JSON(w, r, response.Error("Unauthorized access to categories"))
			return
		}

		// setting saved data for further use
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		next.ServeHTTP(w, r)
	})
}

func (m *OwnershipMiddleware) CheckTaskMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		taskID := chi.URLParam(r, "id")
		if taskID == "" {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("Task ID required"))
			return
		}
		taskUUID, err := uuid.Parse(taskID)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid UUID"))
			return
		}

		userID, ok := r.Context().Value(UserIDContextKey).(uuid.UUID)
		if !ok {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, response.Error("Missing userID"))
			return
		}

		ctx := r.Context()
		ctx, cancel := context.WithTimeout(ctx, m.timeout)
		defer cancel()

		isTaskOwned, err := m.userService.CheckTaskOwnership(ctx, userID, taskUUID)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}
		if !isTaskOwned {
			render.Status(r, http.StatusForbidden)
			render.JSON(w, r, response.Error("Unauthorized access to task"))
			return
		}

		next.ServeHTTP(w, r)

	})
}
