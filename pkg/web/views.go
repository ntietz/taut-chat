package web

import (
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

type IndexPage struct {
	Handles     []string
	CurrentUser string
	FocusedChat string
	Messages    []Message
	Query       string
}

func (h *Handler) Index(c echo.Context) error {
	usernameCookie, err := c.Cookie("username")
	if err != nil {
		fmt.Println("failed to get user")
		return c.Redirect(http.StatusFound, "/login")
	}

	focusedChatCookie, _ := c.Cookie("focusedChat")

	handles, err := ListUserHandles(h.Ts)
	if err != nil {
		fmt.Println("failed to get handles")
		return err
	}

	username := usernameCookie.Value
	focusedChat := ""
	fmt.Println("fcc?", focusedChatCookie)
	if focusedChatCookie != nil {
		focusedChat = focusedChatCookie.Value
	} else {
        focusedChat = username
    }

	messages, err := ListMessages(h.Ts, username, focusedChat)
	if err != nil {
		fmt.Println("failed to get messages", err)
		return err
	}

	viewData := IndexPage{
		Handles:     handles,
		CurrentUser: username,
		FocusedChat: focusedChat,
		Messages:    messages,
		Query:       "",
	}

	return c.Render(http.StatusOK, "index.html", viewData)
}

type SearchPage struct {
	Query       string
	Handles     []string
	CurrentUser string
	Messages    []Message
}

func (h *Handler) Search(c echo.Context) error {
	usernameCookie, err := c.Cookie("username")
	if err != nil {
		fmt.Println("failed to get user")
		return c.Redirect(http.StatusFound, "/login")
	}
	username := usernameCookie.Value

	handles, err := ListUserHandles(h.Ts)
	if err != nil {
		fmt.Println("failed to get handles")
		return err
	}

	query := ""
	err = echo.QueryParamsBinder(c).String("query", &query).BindError()
	if err != nil {
		fmt.Println("failed to get query")
		return err
	}

	messages, err := SearchMessages(h.Ts, username, query)
	if err != nil {
		fmt.Println("failed to get messages")
		return err
	}

	viewData := SearchPage{
		Query:       query,
		Handles:     handles,
		CurrentUser: username,
		Messages:    messages,
	}

	return c.Render(http.StatusOK, "search.html", viewData)
}

func (h *Handler) StartChat(c echo.Context) error {
	handle := ""
	err := echo.PathParamsBinder(c).String("handle", &handle).BindError()
	if err != nil {
		fmt.Println("failed to get user")
		return err
	}

	cookie := new(http.Cookie)
	cookie.Name = "focusedChat"
	cookie.Value = handle
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.Path = "/"
	c.SetCookie(cookie)

	return c.Redirect(http.StatusFound, "/")
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

	err := CreateUser(h.Ts, loginForm.Username)
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
