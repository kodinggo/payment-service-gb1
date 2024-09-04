package helper

import (
	"context"
	"errors"
	"net/http"
	"strings"

	authPb "github.com/kodinggo/user-service-gb1/pb/auth"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

var (
	errUnauthorized = errors.New("unauthorized")
)

type JWTMiddleware struct {
	authClient authPb.JWTValidatorClient
}

func NewJWTMiddleware(authClient authPb.JWTValidatorClient) *JWTMiddleware {
	return &JWTMiddleware{
		authClient: authClient,
	}
}

// ValidateJWT is a middleware function that validates JWT tokens using a gRPC call.
func (m *JWTMiddleware) ValidateJWT(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing token"})
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token format"})
		}

		// Create a context for the gRPC call
		ctx := context.Background()

		// Call gRPC to validate the token
		req := &authPb.ValidateTokenRequest{Token: token}
		res, err := m.authClient.ValidateToken(ctx, req)
		if err != nil || res == nil || !res.Valid {
			logrus.Error(err)
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid or expired token"})
		}

		// Store the user information in the context
		c.Set("user", res.User)

		return next(c)
	}
}

func GetUserSession(c echo.Context) *authPb.User {
	user := c.Get("user")
	if user == nil {
		return nil
	}

	return user.(*authPb.User)
}
