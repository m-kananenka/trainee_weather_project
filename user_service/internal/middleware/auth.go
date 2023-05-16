package middleware

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"net/http"
	"user_service/internal/errorstore"
)

func Auth(keyword string, l *logrus.Logger) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		fn := func(c echo.Context) error {
			tokenStr := c.Request().Header.Get("Authorization")

			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				return []byte(keyword), nil
			})

			if err != nil {
				l.Info("failed to authorize:", err)
				return echo.NewHTTPError(http.StatusUnauthorized, errorstore.Unauthorized(err))
			}

			idStr, err := token.Claims.GetIssuer()
			if err != nil {
				l.Info("failed to authorize:", err)
				return echo.NewHTTPError(http.StatusUnauthorized, errorstore.Unauthorized(err))
			}

			IdCtx := "userIdCtx"

			ctx := context.Background()
			ctx = context.WithValue(ctx, IdCtx, idStr)

			err = next(c)
			if err != nil {
				return err
			}
			return nil
		}

		return fn
	}
}
