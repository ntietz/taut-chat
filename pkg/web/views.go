package web

import (
	//"context"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
	_, ok := c.Get("user").(*jwt.Token)
	if !ok {
        fmt.Println("failed to get user")
		return c.Redirect(http.StatusFound, "/login")
	}

	return c.Render(http.StatusOK, "index.html", "World!")
}

type LoginForm struct {
	Username string `form:"username" validate:"required"`
	Password string `form:"password" validate:"required"`
}

func (h *Handler) Login(c echo.Context) error {
	loginForm := LoginForm{
		Username: "",
		Password: "",
	}
	return c.Render(http.StatusOK, "login.html", loginForm)
}

func (h *Handler) LoginAttempt(c echo.Context) error {
	loginForm := LoginForm{
		Username: "",
		Password: "",
	}

	if err := c.Bind(&loginForm); err != nil {
		return c.Render(http.StatusOK, "login.html", loginForm)
	}
	if loginForm.Username == "" || loginForm.Password == "" {
		fmt.Println("Username/password must be non-empty; username?", loginForm.Username, "password?", loginForm.Password)
		return c.Render(http.StatusOK, "login.html", loginForm)
	}

	loginSuccess, err := CheckLogin(h.Ts, loginForm.Username, loginForm.Password)
	if err != nil {
		return err
	}

	if loginSuccess {
		// TODO: set cookie

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": loginForm.Username,
			// this is a demo example, so no expiration, YOLO
		})

		// Sign and get the complete encoded token as a string using the secret
		tokenString, err := token.SignedString([]byte("secret"))
		// if the signing fails, just fail closed onto login screen... YOLO
		if err != nil {
            fmt.Println("failed to sign", err)
			return c.Redirect(http.StatusFound, "/login")
		}

		cookie := new(http.Cookie)
		cookie.Name = "authy"
		cookie.Value = tokenString
		cookie.Expires = time.Now().Add(24 * time.Hour)
		c.SetCookie(cookie)


		return c.Redirect(http.StatusFound, "/")
	} else {
		return c.Render(http.StatusOK, "login.html", loginForm)
	}
}
