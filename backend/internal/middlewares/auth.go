package middlewares

import (
	"context"
	"fmt"
	"github.com/RadekKusiak71/places-app/internal/errors"
	"github.com/RadekKusiak71/places-app/internal/jwt"
	"github.com/RadekKusiak71/places-app/internal/stores"
	"github.com/RadekKusiak71/places-app/internal/utils"
	"log"
	"net/http"
	"strings"
)

type contextKey string

const UserContextKey = contextKey("user")

var (
	ExpectedHeader = "bearer"
)

func parseAndExtractToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.MissingAuthorizationHeader()
	}

	splitHeader := strings.Split(authHeader, " ")
	if len(splitHeader) != 2 {
		return "", errors.InvalidAuthorizationHeader()
	}

	tokenHeader := splitHeader[0]
	if strings.Compare(strings.ToLower(tokenHeader), ExpectedHeader) != 0 {
		return "", errors.InvalidAuthorizationHeader()
	}

	return splitHeader[1], nil
}

func AuthMiddleware(next utils.APIFunc, userStore *stores.UserStore) utils.APIFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		tokenString, err := parseAndExtractToken(r)
		if err != nil {
			return err
		}

		claims, err := jwt.ParseAndVerifyAccessToken(tokenString)
		if err != nil {
			log.Println(err.Error())
			return errors.NewAPIError(http.StatusUnauthorized, fmt.Sprintf(err.Error()))
		}

		rCtx := r.Context()

		user, err := userStore.Get(rCtx, claims.UserID)
		if err != nil {
			return errors.InvalidTokenError()
		}

		dCtx := context.WithValue(
			rCtx,
			UserContextKey,
			user,
		)

		return next(
			w,
			r.WithContext(dCtx),
		)
	}
}
