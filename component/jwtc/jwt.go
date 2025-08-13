package jwtc

import (
	"context"
	"flag"
	"fmt"

	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"github.com/teoit/gosctx"
)

/**
 * Default exp
 * token: 1h
 * refresh 7day
 */
const (
	defaultSecret                      = "123456-Xyz@2023-TaOn-Chuasangsoi"
	defaultExpireTokenInSeconds        = 60 * 60 * 24
	defaultExpireRefreshTokenInSeconds = 60 * 60 * 24 * 7

	ACCESS_TOKEN  = "ACCESS_TOKEN"
	REFRESH_TOKEN = "REFRESH_TOKEN"
)

var (
	ErrSecretKeyNotValid     = errors.New("secret key must be in 32 bytes")
	ErrTokenLifeTimeTooShort = errors.New("token life time too short")

	TokenTypes = map[string]bool{
		ACCESS_TOKEN:  true,
		REFRESH_TOKEN: true,
	}
)

type JWTProvider interface {
	IssueToken(ctx context.Context, id, sub string, tokenType *string) (token string, expSecs int, err error)
	ParseToken(ctx context.Context, tokenString string) (claims *jwt.RegisteredClaims, err error)
	GetExpireToken(tokenType *string) int
}

type jwtx struct {
	id                          string
	secret                      string
	expireTokenInSeconds        int
	expireRefreshTokenInSeconds int
}

func NewJWT(id string) *jwtx {
	return &jwtx{id: id}
}

func (j *jwtx) ID() string {
	return j.id
}

func (j *jwtx) InitFlags() {
	flag.StringVar(
		&j.secret,
		"jwt-secret",
		defaultSecret,
		"Secret key to sign JWT",
	)
	flag.IntVar(
		&j.expireTokenInSeconds,
		"jwt-exp-secs",
		defaultExpireTokenInSeconds,
		"Number of seconds token will expired",
	)

	flag.IntVar(
		&j.expireRefreshTokenInSeconds,
		"jwt-exp-refesh-token-secs",
		defaultExpireRefreshTokenInSeconds,
		"Number of seconds token will expired",
	)
}

func (j *jwtx) Activate(_ gosctx.ServiceContext) error {
	if len(j.secret) < 32 {
		return errors.WithStack(ErrSecretKeyNotValid)
	}

	if j.expireTokenInSeconds <= 60 {
		return errors.WithStack(ErrTokenLifeTimeTooShort)
	}

	return nil
}

func (j *jwtx) Stop() error {
	return nil
}

func (j *jwtx) IssueToken(ctx context.Context, id, sub string, tokenType *string) (token string, expSecs int, err error) {
	now := time.Now().UTC()
	expToken := j.GetExpireToken(tokenType)

	claims := jwt.RegisteredClaims{
		Subject:   sub,
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Second * time.Duration(expToken))),
		NotBefore: jwt.NewNumericDate(now),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        id,
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenSignedStr, err := t.SignedString([]byte(j.secret))

	if err != nil {
		return "", 0, errors.WithStack(err)
	}

	return tokenSignedStr, expToken, nil
}

func (j *jwtx) ParseToken(ctx context.Context, tokenString string) (claims *jwt.RegisteredClaims, err error) {
	var rc jwt.RegisteredClaims

	token, err := jwt.ParseWithClaims(tokenString, &rc, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(j.secret), nil
	})

	if !token.Valid {
		return nil, errors.WithStack(err)
	}

	return &rc, nil
}

func (j *jwtx) GetExpireToken(tokenType *string) int {
	if tokenType == nil || *tokenType == ACCESS_TOKEN {
		return j.expireTokenInSeconds
	}

	if *tokenType == REFRESH_TOKEN {
		return j.expireRefreshTokenInSeconds
	}

	return j.expireTokenInSeconds
}
