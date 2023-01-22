package user_handlers

import (
	"fmt"
	"net/http"
	"prmeet/internal/auth"
	"prmeet/internal/er"
	"prmeet/internal/natsio"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func IndexPage(c echo.Context) error {
	feed := auth.LoadAppData(c)
	return c.Render(http.StatusOK, "default.go.html", feed)
}

// show login form
func UserLogin(c echo.Context) error {

	feed := auth.LoadAppData(c)

	if feed["user"] != nil {
		sess, err := session.Get("session", c)
		er.ErrorPrint(err)
		sess.Values["success"] = "You are logged in."
		err = sess.Save(c.Request(), c.Response())
		if err != nil {
			fmt.Println(err)
		}
		return c.Redirect(http.StatusMovedPermanently, "default.go.html")
	}

	return c.Render(http.StatusOK, "login.go.html", feed)
}

// show do login request
func UserDoLogin(c echo.Context) error {

	feed := auth.LoadAppData(c)

	email := c.Request().PostFormValue("email")
	password := c.Request().PostFormValue("password")

	err := validation.Validate(email,
		validation.Required, // not empty
		is.Email,            // is a valid email
	)

	if err != nil {
		feed["err"] = err
		return c.Render(http.StatusBadRequest, "login.go.html", feed)
	}

	request := make(map[string]any)
	request["email"] = email

	response := natsio.AskPipe("users.login.getByEmail", request)

	if response["err"] != nil {
		fmt.Printf("\nUser not found: %v \nand the error is: %v", email, response["err"])
		feed["err"] = "User not found."
		return c.Render(http.StatusNotFound, "login.go.html", feed)
	}

	user := response["user"].(map[string]interface{})

	ok := auth.PasswordCheckHash(password, user["password"].(string))

	if !ok {
		feed["err"] = "Wrong password."
		return c.Render(http.StatusNotFound, "login.go.html", feed)
	}

	sess, err := session.Get("session", c)
	if err != nil {
		fmt.Println(err)
	}
	sess.Values["user"] = user

	fmt.Printf("The sess.Values user is: %v", sess.Values["user"])

	err = sess.Save(c.Request(), c.Response())
	er.ErrorPrint(err)

	feed["user"] = user
	feed["success"] = "Successfully logged in."

	return c.Render(http.StatusOK, "default.go.html", feed)
}

func UserLogout(c echo.Context) error {

	feed := auth.LoadAppData(c)

	sess, err := session.Get("session", c)
	er.ErrorPrint(err)

	if sess.Values["user"] == nil {
		feed["success"] = "You did not login."
		return c.Render(http.StatusBadRequest, "default.go.html", feed)
	}

	sess.Values["user"] = nil

	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		fmt.Println(err)
	}

	feed["success"] = "Successfully logged out."

	return c.Render(http.StatusOK, "login.go.html", feed)

}
