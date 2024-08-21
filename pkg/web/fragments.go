package web

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/typesense/typesense-go/v2/typesense"
)

type FragmentsHandler struct {
    Ts *typesense.Client
}

type UserListFragment struct {
    Handles []string
}

func (h *FragmentsHandler) UserList(c echo.Context) error {
    handles, err := ListUserHandles(h.Ts)
	if err != nil {
        fmt.Println("failed to get handles")
        return err
    }

    viewData := UserListFragment{
    	Handles:     handles,
    }

	return c.Render(http.StatusOK, "fragment_user_list.html", viewData)
}

type ChatWindowFragment struct {
    FocusedChat string
}

func (h *FragmentsHandler) ChatWindow(c echo.Context) error {
	focusedChatCookie, _ := c.Cookie("focusedChat")
	focusedChat := ""
	if focusedChatCookie != nil {
		focusedChat = focusedChatCookie.Value
	}

    viewData := ChatWindowFragment{
        FocusedChat: focusedChat,
    }

	return c.Render(http.StatusOK, "fragment_chat_window.html", viewData)
}

type MessagesFragment struct {
    Messages []Message
}

func (h *FragmentsHandler) Messages(c echo.Context) error {
	usernameCookie, err := c.Cookie("username")
	if err != nil {
		fmt.Println("failed to get user")
		return c.Redirect(http.StatusFound, "/login")
	}

	focusedChatCookie, _ := c.Cookie("focusedChat")

	username := usernameCookie.Value
	focusedChat := ""
    fmt.Println("fcc?", focusedChatCookie)
	if focusedChatCookie != nil {
		focusedChat = focusedChatCookie.Value
	}

    viewData := MessagesFragment{
        Messages: make([]Message, 0),
    }

	messages, err := ListMessages(h.Ts, username, focusedChat)
	if err != nil {
		fmt.Println("failed to get messages", err)
		return err
	}

    viewData.Messages = messages

	return c.Render(http.StatusOK, "fragment_messages.html", viewData)
}

type ChatForm struct {
    Message string `form:"message" validate:"required"`
}

type SendChatFragment struct {
	CurrentUser string
	FocusedChat string
	Messages    []Message
}

func (h *Handler) SendChat(c echo.Context) error {
	handle := ""
	err := echo.PathParamsBinder(c).String("handle", &handle).BindError()
	if err != nil {
		fmt.Println("failed to get user")
		return err
	}
    chatForm := ChatForm {
        Message: "",
    }
    if err := c.Bind(&chatForm); err != nil {
        return err
    }
    fmt.Println("chatForm?", chatForm)

	usernameCookie, err := c.Cookie("username")
	if err != nil {
		fmt.Println("failed to get user")
		return c.Redirect(http.StatusFound, "/login")
	}

	focusedChatCookie, _ := c.Cookie("focusedChat")

	username := usernameCookie.Value
	focusedChat := ""
    fmt.Println("fcc?", focusedChatCookie)
	if focusedChatCookie != nil {
		focusedChat = focusedChatCookie.Value
	}

    err = CreateMessage(h.Ts, username, focusedChat, chatForm.Message)
    if err != nil {
		fmt.Println("failed to save message", err)
		return err
    }

	messages, err := ListMessages(h.Ts, username, focusedChat)
	if err != nil {
		fmt.Println("failed to get messages", err)
		return err
	}


    viewData := SendChatFragment{
    	CurrentUser: username,
    	FocusedChat: focusedChat,
    	Messages:    messages,
    }


	return c.Render(http.StatusOK, "fragment_chat_window.html", viewData)
}


