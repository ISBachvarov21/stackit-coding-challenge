package routes

import (
	"context"
	"net/http"
	"strings"

	"github.com/ISBachvarov21/stackit-coding-challenge/models"
	"github.com/julienschmidt/httprouter"
)

type contextKey string

const userIdKey contextKey = "userId"

// withUser validates the Bearer token as a user ID and injects it into the request context.
func withUser(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userId := strings.TrimPrefix(auth, "Bearer ")
		if err := db.Model(&models.User{Id: userId}).WherePK().Select(); err != nil {
			http.Error(w, "User with id "+userId+" not found", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIdKey, userId)
		next(w, r.WithContext(ctx), p)
	}
}

func userIdFromContext(r *http.Request) string {
	userId, _ := r.Context().Value(userIdKey).(string)
	return userId
}
