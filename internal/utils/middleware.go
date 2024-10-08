package utils

import (
	"errors"

	"github.com/golang-jwt/jwt/v4"
	"github.com/kodinggo/payment-service-gb1/internal/config"
	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

var (
	errUnauthorized = errors.New("unauthorized")
)

type JWTClaims struct {
	jwt.RegisteredClaims
	ID   string
	Role string
}

func JWTConfig() echojwt.Config {
	c := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(JWTClaims)
		},
		SigningKey: []byte(config.GetJwtSecret()),
	}

	return c
}

func UserClaims(c echo.Context) (*JWTClaims, error) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JWTClaims)

	if claims == nil {
		logrus.WithField("ctx", Dump(c)).Error(errUnauthorized)
		return nil, errUnauthorized
	}

	return claims, nil
}
