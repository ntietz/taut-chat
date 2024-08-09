package views

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) Index(c echo.Context) error {
    keys, err := h.ts.Keys().Retrieve(context.Background())
    if err != nil {
        fmt.Println("error:", err)
    }
    fmt.Println("keys:", keys)

    return c.Render(http.StatusOK, "index.html", "World!")
}
