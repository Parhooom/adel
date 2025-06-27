package middleware

import (
	"adel/internal/service/postgres"
	"adel/internal/utils"
	"context"
	"net/http"
	"strings"
)

type UserMiddleware struct {
	UserStore postgres.UserStore
}

type contextKey string

const UserContextKey contextKey = "user"

func SetUser(r *http.Request, user *postgres.User) *http.Request {
	ctx := context.WithValue(r.Context(), UserContextKey, user)
	return r.WithContext(ctx)
}

func GetUser(r *http.Request) *postgres.User {
	user, ok := r.Context().Value(UserContextKey).(*postgres.User)
	if !ok {
		return postgres.AnonymousUser
	}

	return user
}

func (m *UserMiddleware) AuthenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			r = SetUser(r, postgres.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authHeader, " ") // Bearer <token>
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid authorization header"})
			return
		}

		token := headerParts[1]
		user, err := m.UserStore.GetUserToken(token)
		if err != nil {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid token"})
			return
		}

		if user == nil {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "token expired or invalid"})
			return
		}

		r = SetUser(r, user)
		next.ServeHTTP(w, r)
	})
}

func (m *UserMiddleware) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r)

		if user.IsAnonymous() {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "you must be logged in to access this resource"})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *UserMiddleware) RequireAdminUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r)

		if !user.IsAdmin {
			utils.WriteJSON(w, http.StatusForbidden, utils.Envelope{"error": "you must be an admin to access this resource"})
			return
		}

		next.ServeHTTP(w, r)
	})
}
