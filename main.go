package main

import (
	"encoding/gob"
	"html/template"
	"io"
	mu "prmeet/internal/utils/my_utils"
	myrouter "prmeet/router"

	"github.com/gobuffalo/envy"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	envy.Load()
	e := echo.New()
	e.Static("/public", "views/public")
	e.Use(session.Middleware(mu.GetCookieStore()))
	//now we can use it in the session
	gob.Register(map[string]interface{}{})
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("views/layout/**/*.html")),
	}
	e.Renderer = renderer
	myrouter.SetupRoutes(e)
	e.Logger.Fatal(e.Start("localhost:8001"))
}
