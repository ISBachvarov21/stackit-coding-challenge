package routes

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/ISBachvarov21/stackit-coding-challenge/models"
	"github.com/julienschmidt/httprouter"
)

type contextKey string

const userIdKey contextKey = "userId"

// writeError sends a JSON error response with the given status code.
func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": msg}); err != nil {
		log.Println("Failed to encode error response:", err)
	}
}

// withUser validates the Bearer token as a user ID and injects it into the request context.
func withUser(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			writeError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		userId := strings.TrimPrefix(auth, "Bearer ")
		if err := db.Model(&models.User{Id: userId}).WherePK().Select(); err != nil {
			writeError(w, http.StatusUnauthorized, "User with id "+userId+" not found")
			return
		}

		ctx := context.WithValue(r.Context(), userIdKey, userId)
		next(w, r.WithContext(ctx), p)
	}
}

// withOwner validates the Bearer token against the owner UUID.
func withOwner(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			writeError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		if strings.TrimPrefix(auth, "Bearer ") != ownerUUID {
			writeError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		next(w, r, p)
	}
}

func userIdFromContext(r *http.Request) string {
	userId, _ := r.Context().Value(userIdKey).(string)
	return userId
}
