package web

import (
	"html/template"

    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
	"github.com/ntietz/taut-chat/pkg/web/views"
)

func CreateServer() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	templateRenderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}

	e.Renderer = templateRenderer
    e.Static("/static", "static")

	e.GET("/", views.Index)

	return e
}
