package main

import (
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
	"net/http"
	"strconv"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type HourInfo struct {
	Hour      int
	Qty       int
	Completed bool
}

type DayInfo struct {
	Name        string
	MaxQty      int
	Overflow    int
	Date        string
	MorningMash bool
	EveningMash bool
	Hours       []HourInfo
}

func main() {
	e := echo.New()

	renderer := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.gohtml")),
	}
	e.Renderer = renderer

	e.GET("/", func(c echo.Context) error {
		day := GetCurrentDay()

		return c.Render(http.StatusOK, "index", day)
	})

	//front api
	e.PUT("/api/hour/:hour", func(c echo.Context) error {
		hour := c.Param("hour")
		hourInt, _ := strconv.Atoi(hour)
		qty, _ := strconv.Atoi(c.FormValue("qty_" + hour))

		updatedHour := UpdateHour(hourInt, qty)
		overflow := GetOverflow()

		return c.Render(http.StatusOK, "completed_hour", map[string]interface{}{
			"NewHour":  updatedHour,
			"Overflow": overflow,
		})
	})

	e.Logger.Fatal(e.Start(":1323"))
}
