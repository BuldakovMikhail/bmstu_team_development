package middleware

import (
	"context"
	"net/http"
	"strings"
	auth_utils "todolist/internal/pkg/authUtils"
	"todolist/internal/pkg/response"

	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

func NewJwtAuthMiddleware(loggerSrc *logrus.Logger, secretSrc string, tokenHandlerSrc auth_utils.ITokenHandler) JwtAuthMiddleware {
	return JwtAuthMiddleware{
		secret:       secretSrc,
		tokenHandler: tokenHandlerSrc,
		logger:       loggerSrc,
	}
}

type JwtAuthMiddleware struct {
	logger       *logrus.Logger
	secret       string
	tokenHandler auth_utils.ITokenHandler
}

var (
	UserIDContextKey = "contextKeyRole{}"
	RoleContextKey   = "contextKeyID{}"
)

func (m *JwtAuthMiddleware) MiddlewareFunc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			m.logger.Info("user with no token came")
			render.JSON(w, r, response.Error("Error in parsing token"))
			render.Status(r, http.StatusBadRequest)
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")

		payload, err := m.tokenHandler.ParseToken(token, m.secret)
		if err != nil {
			if err == auth_utils.ErrParsingToken {
				m.logger.Info("user with invalid jwt came")
				render.JSON(w, r, response.Error(err.Error()))
				render.Status(r, http.StatusBadRequest)
			} else {
				m.logger.Info("user with invalid jwt came")
				render.JSON(w, r, response.Error(err.Error()))
				render.Status(r, http.StatusUnauthorized)
			}
			return
		}
		ctx := context.WithValue(r.Context(), UserIDContextKey, payload.ID)

		m.logger.WithFields(
			logrus.Fields{
				"src":    "JwtAuthMiddleware.MiddleFunc",
				"userID": payload.ID,
			}).
			Info("successfully authorized")

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
