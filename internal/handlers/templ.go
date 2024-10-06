package handlers

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

// render renders templ a component using the provided echo.Context.
func render(c echo.Context, component templ.Component) error {
	return component.Render(c.Request().Context(), c.Response())
}
