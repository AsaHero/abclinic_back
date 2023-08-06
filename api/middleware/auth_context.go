package middleware

import (
	"context"
	"net/http"

	tokenpkg "github.com/AsaHero/abclinic/internal/pkg/token"
)

func AuthContext(jwtsecret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var (
				token    string
				authData = make(map[string]string)
			)

			if token = r.Header.Get("Authorization"); len(token) > 10 {
				token = token[7:]
			}

			claims, err := tokenpkg.ParseJwtToken(token, jwtsecret)
			if err == nil && len(claims) != 0 {
				for k, v := range claims {
					if value, ok := v.(string); ok {
						authData[k] = value
					}
				}

				next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), CtxKeyAuthData, authData)))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
