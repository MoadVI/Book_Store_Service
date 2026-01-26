package middleware

import (
	"Book-Store/internal/authentication"
	"Book-Store/internal/response"
	"context"
	"net/http"
)

type key int

const UserIDKey key = 0

func AuthMiddleware(tokenSecret string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := authentication.GetCustomerToken(r.Header)
		if err != nil {
			response.RespondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		userID, err := authentication.ValidateJWT(token, tokenSecret)
		if err != nil {
			response.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserIDFromContext(ctx context.Context) int {
	if ctx == nil {
		return 0
	}
	if id, ok := ctx.Value(UserIDKey).(int); ok {
		return id
	}
	return 0
}
