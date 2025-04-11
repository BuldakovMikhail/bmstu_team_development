package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"
	"todolist/internal/models"
	auth_utils "todolist/internal/pkg/authUtils"
	"todolist/internal/pkg/response"

	"github.com/go-chi/render"
)

type UserInfo struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Token struct {
	Token string `json:"token"`
}

type AuthProvider interface {
	SignIn(candidate *models.UserAuth) (tokenStr string, err error)
	SignUp(candidate *models.UserAuth) error
}

func FromUserInfo(userDTO UserInfo) *models.UserAuth {
	return &models.UserAuth{
		Name:     userDTO.Name,
		Password: userDTO.Password,
	}
}

// @Summary SignIn
// @Tags user
// @Description SignIn
// @ID sign-in
// @Accept  json
// @Produce  json
// @Param input body UserInfo true "user's name and password"
// @Success 200 {object} Token
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/v1/sign-in [post]
func SignIn(authProvider AuthProvider, timeout time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UserInfo
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		ctx := r.Context()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		tokenStr, err := authProvider.SignIn(FromUserInfo(req))
		if err != nil {
			if errors.Is(err, auth_utils.ErrInvalidToken) {
				render.Status(r, http.StatusBadRequest)
			} else if errors.Is(err, models.ErrUserNotFound) {
				render.Status(r, http.StatusNotFound)
			}
			render.JSON(w, r, response.Error(err.Error()))
			return
		}
		render.JSON(w, r, Token{Token: tokenStr})
	}
}

// @Summary SignUp
// @Tags user
// @Description SignUp
// @ID sign-up
// @Accept  json
// @Produce  json
// @Param input body UserInfo true "user's name and password"
// @Success 200 {object} Token
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/v1/sign-up [post]
func SignUp(authProvider AuthProvider, timeout time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UserInfo
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}
		err = authProvider.SignUp(FromUserInfo(req))
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}
		render.Status(r, http.StatusOK)
	}
}
