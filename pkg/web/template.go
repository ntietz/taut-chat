package web

import (
	"errors"
	"fmt"
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    template := t.templates.Lookup(name)
    if template == nil {
        msg := fmt.Sprintf("missing template '%s'", name)
        return errors.New(msg)
    }
    return template.Execute(w, data)

}
