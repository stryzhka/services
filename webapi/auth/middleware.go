package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type ctxKey string

const userCtxKey ctxKey = "user"

type Claims struct {
	jwt.RegisteredClaims
	UserId string `json:"uid"`
}

type Config struct {
	Secret   []byte
	Issuer   string
	Audience string
}

func Middleware(cfg Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				writeUnauthorized(w, "missing authorization header")
				return
			}
			raw, ok := strings.CutPrefix(header, "Bearer ")
			if !ok {
				writeUnauthorized(w, "expected Bearer scheme")
				return
			}

			parserOpts := []jwt.ParserOption{
				jwt.WithValidMethods([]string{"HS256"}),
				jwt.WithExpirationRequired(),
			}
			if cfg.Issuer != "" {
				parserOpts = append(parserOpts, jwt.WithIssuer(cfg.Issuer))
			}
			if cfg.Audience != "" {
				parserOpts = append(parserOpts, jwt.WithAudience(cfg.Audience))
			}
			parser := jwt.NewParser(parserOpts...)

			claims := &Claims{}
			token, err := parser.ParseWithClaims(raw, claims, func(t *jwt.Token) (any, error) {
				return cfg.Secret, nil
			})
			if err != nil || !token.Valid {
				writeUnauthorized(w, "invalid token")
				return
			}

			ctx := context.WithValue(r.Context(), userCtxKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ClaimsFrom(ctx context.Context) (*Claims, bool) {
	c, ok := ctx.Value(userCtxKey).(*Claims)
	return c, ok
}

func writeUnauthorized(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte(`{"error":"` + msg + `"}`))
}
