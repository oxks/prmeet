package user

import (
	"errors"
	"net/http"
	"prmeet/internal/auth"
	"prmeet/internal/er"
	"prmeet/internal/myvalidation"
	"prmeet/internal/natsio"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func UserSignup(c echo.Context) error {

	feed := make(map[string]any)

	sess, err := session.Get("session", c)
	er.ErrorPrint(err)

	if sess.Values["user"] != nil {
		feed["user"] = sess.Values["user"]
		return c.Render(http.StatusOK, "default.go.html", feed)
	}
	feed["user"] = sess.Values["user"]

	return c.Render(http.StatusOK, "signup.go.html", feed)
}
func UserDoSignup(c echo.Context) error {
	feed := auth.LoadAppData(c)

	//if is signed up successfully
	sess, err := session.Get("session", c)
	er.ErrorPrint(err)

	if sess.Values["user"] != nil {
		feed["success"] = "You are signed up."
		return c.Render(http.StatusContinue, "default.go.html", feed)
	}
	u := myvalidation.SignupUserParams{}
	u.Email = c.FormValue("email")
	u.Nickname = c.FormValue("nickname")
	u.Password = c.FormValue("password")

	ers := []error{}
	err = u.Validate()
	// u.Validate() returns lowercase properties, so below the workaround
	if err != nil {
		err = er.FixValidationResult(err)

		feed["err"] = append(ers, err)
		return c.Render(http.StatusBadRequest, "signup.go.html", feed)
	}

	// validate password match
	if c.FormValue("password") != c.FormValue("repeat_password") {
		feed["err"] = append(ers, errors.New("please retype password to match"))
		return c.Render(http.StatusBadRequest, "signup.go.html", feed)
	}

	u.Password, err = auth.PasswordHash(u.Password)
	er.ErrorPrint(err)

	request := make(map[string]any)
	request["user"] = u

	response := natsio.AskPipe("users.signup", request)
	if response["err"] != nil || response["user"] == nil {
		feed["err"] = "Nickname and email must be unique, please try again."
		return c.Render(http.StatusBadRequest, "signup.go.html", feed)
	}

	user := response["user"].(map[string]interface{})

	auth.SaveToSession(c, "user", user)
	feed["user"] = user
	feed["err"] = ers
	return c.Render(http.StatusOK, "default.go.html", feed)
}
