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

func HtmlShenanigans(s string) template.HTML {
	return template.HTML(s)
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	templ := t.templates.Lookup(name)
	if templ == nil {
		msg := fmt.Sprintf("missing template '%s'", name)
		return errors.New(msg)
	}
	return templ.Funcs(template.FuncMap{
        "html_unsafely": HtmlShenanigans,
    }).Execute(w, data)

}
