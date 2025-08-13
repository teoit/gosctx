package configs

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
)

type JWTProvider interface {
	IssueToken(ctx context.Context, id, sub string) (token string, expSecs int, err error)
	ParseToken(ctx context.Context, tokenString string) (claims *jwt.RegisteredClaims, err error)
}
