package middleware

import (
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

func Logger(l *logrus.Logger) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		fn := func(c echo.Context) error {
			c.Request()
			l.Info(c.Request().URL.String())
			err := next(c)
			if err != nil {
				return err
			}
			return nil
		}

		return fn
	}
}
