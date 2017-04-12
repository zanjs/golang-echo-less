package main

import (
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	e.Static("/static", "static")

	e.GET("/", index)

	e.POST("/users", users)

	e.POST("/save", save)

	e.Logger.Fatal(e.Start(":1212"))
}

func index(c echo.Context) error {
	cookieVal := &http.Cookie{
		Name:     "uid",
		Value:    "1212",
		Path:     "/",
		HttpOnly: false,
	}
	c.SetCookie(cookieVal)

	WriteCookie(c)

	return c.String(http.StatusOK, "ok")
}

func save(c echo.Context) error {
	// Get name
	name := c.FormValue("name")
	// Get avatar
	avatar, err := c.FormFile("avatar")
	if err != nil {
		return err
	}

	// Source
	src, err := avatar.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	dst, err := os.Create(avatar.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return c.HTML(http.StatusOK, "<b>Thank you! "+name+"</b>")
}

func users(c echo.Context) error {

	// User is
	type User struct {
		Name  string `json:"name" xml:"name" form:"name" query:"name"`
		Email string `json:"email" xml:"email" form:"email" query:"name"`
	}

	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, u)
	// or
	// return c.XML(http.StatusCreated, u)
}
