package web

import (
	"fmt"
	"html/template"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func CreateServer() *echo.Echo {
	h := NewHandler()
    //DropCollections(h.Ts)
	err := CreateCollections(h.Ts)
	if err != nil {
		fmt.Println("", err)
	}

	fh := FragmentsHandler{
		Ts: h.Ts,
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	templates := template.Must(template.New("").Funcs(template.FuncMap{
		"html_unsafely": HtmlShenanigans,
	}).ParseGlob("views/*.html"))

	templateRenderer := &TemplateRenderer{
		templates: templates,
	}

	e.Renderer = templateRenderer
	e.Static("/static", "static")

	e.GET("/", h.Index)
	e.GET("/login", h.Login)
	e.GET("/search", h.Search)
	e.GET("/start-chat/:handle", h.StartChat)
	e.POST("/login", h.LoginAttempt)
	e.POST("/send-chat/:handle", h.SendChat)

	e.GET("/fragment/users", fh.UserList)
	e.GET("/fragment/chat", fh.ChatWindow)
	e.GET("/fragment/messages", fh.Messages)

	return e
}
