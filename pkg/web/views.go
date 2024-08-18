package web

import (
	//"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/typesense/typesense-go/v2/typesense"
)

type Handler struct {
	Ts *typesense.Client
}

func NewHandler() *Handler {
	ts := typesense.NewClient(
		typesense.WithServer("http://localhost:8108"),
		typesense.WithAPIKey("1667b96f-da3c-40f9-a3b5-8b461a78ed68"))
	return &Handler{
		ts,
	}
}

func (h *Handler) Index(c echo.Context) error {
	usernameCookie, err := c.Cookie("username")
	if err != nil {
        fmt.Println("failed to get user")
		return c.Redirect(http.StatusFound, "/login")
	}

    username := usernameCookie.Value

	return c.Render(http.StatusOK, "index.html", username)
}

type LoginForm struct {
	Username string `form:"username" validate:"required"`
}

func (h *Handler) Login(c echo.Context) error {
	loginForm := LoginForm{
		Username: "",
	}
	return c.Render(http.StatusOK, "login.html", loginForm)
}

func (h *Handler) LoginAttempt(c echo.Context) error {
	loginForm := LoginForm{
		Username: "",
	}

	if err := c.Bind(&loginForm); err != nil {
		return c.Render(http.StatusOK, "login.html", loginForm)
	}
	if loginForm.Username == "" {
		fmt.Println("Username must be non-empty; username?", loginForm.Username)
		return c.Render(http.StatusOK, "login.html", loginForm)
	}

	_, err := CreateUser(h.Ts, loginForm.Username)
	if err != nil {
		return err
	}

    cookie := new(http.Cookie)
    cookie.Name = "username"
    cookie.Value = loginForm.Username
    cookie.Expires = time.Now().Add(24 * time.Hour)
    c.SetCookie(cookie)

    return c.Redirect(http.StatusFound, "/")
}
