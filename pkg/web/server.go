package web

import (
	"fmt"
	"html/template"

	"github.com/labstack/echo-jwt/v4"
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
	e.Use(echojwt.WithConfig(echojwt.Config{
		Skipper: func(c echo.Context) bool {
			if c.Path() == "/login" {
				return true
			}
			return false
		},
		ErrorHandler: func(c echo.Context, err error) error {
            fmt.Println("error in jwt middleware:", err)
			return nil
		},
		ContinueOnIgnoredError: true,
		TokenLookup:            "cookie:authy",
		SigningKey:             []byte("secret"),
	}))

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
