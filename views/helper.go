package views

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func Renderizar(cmp templ.Component, c echo.Context) error {
	return cmp.Render(c.Request().Context(), c.Response().Writer)
}
