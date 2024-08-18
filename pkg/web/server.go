package web

import (
	"fmt"
	"html/template"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func CreateServer() *echo.Echo {
	h := NewHandler()
	err := CreateCollections(h.Ts)
	if err != nil {
		fmt.Println("", err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	templateRenderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}

	e.Renderer = templateRenderer
	e.Static("/static", "static")

	e.GET("/", h.Index)
	e.GET("/login", h.Login)
	e.POST("/login", h.LoginAttempt)

	return e
}
