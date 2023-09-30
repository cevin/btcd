package web

import "github.com/labstack/echo/v4"

func DefaultNotFound(c echo.Context) error {
	return c.JSON(404, errOutput{404, "undefined endpoint"})
}
