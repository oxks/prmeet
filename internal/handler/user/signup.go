package user_handlers

import (
	"errors"
	"net/http"
	me "prmeet/internal/utils/my_errors"
	mu "prmeet/internal/utils/my_utils"
	"prmeet/internal/utils/my_validation"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func UserSignup(c echo.Context) error {

	feed := mu.M{}

	sess, err := session.Get("session", c)
	me.ErrorPrint(err)

	if sess.Values["user"] != nil {
		feed["user"] = sess.Values["user"]
		return c.Render(http.StatusOK, "default.go.html", feed)
	}
	feed["user"] = sess.Values["user"]

	return c.Render(http.StatusOK, "signup.go.html", feed)
}
func UserDoSignup(c echo.Context) error {
	feed := mu.LoadAppData(c)

	//if is signed up successfully
	sess, err := session.Get("session", c)
	me.ErrorPrint(err)

	if sess.Values["user"] != nil {
		feed["success"] = "You are signed up."
		return c.Render(http.StatusContinue, "default.go.html", feed)
	}
	u := my_validation.SignupUserParams{}
	u.Email = c.FormValue("email")
	u.Nickname = c.FormValue("nickname")
	u.Password = c.FormValue("password")

	er := []error{}
	err = u.Validate()
	// u.Validate() returns lowercase properties, so below the workaround
	if err != nil {
		err = me.FixValidationResult(err)

		feed["err"] = append(er, err)
		return c.Render(http.StatusBadRequest, "signup.go.html", feed)
	}

	// validate password match
	if c.FormValue("password") != c.FormValue("repeat_password") {
		feed["err"] = append(er, errors.New("Please retype password to match!"))
		return c.Render(http.StatusBadRequest, "signup.go.html", feed)
	}

	u.Password, err = mu.PasswordHash(u.Password)
	me.ErrorPrint(err)

	request := mu.M{}
	request["user"] = u

	response := mu.AskPipe("users.signup", request)
	if response["err"] != nil || response["user"] == nil {
		feed["err"] = "Nickname and email must be unique, please try again."
		return c.Render(http.StatusBadRequest, "signup.go.html", feed)
	}

	user := response["user"].(map[string]interface{})

	mu.SaveToSession(c, "user", user)
	feed["user"] = user
	feed["err"] = er
	return c.Render(http.StatusOK, "default.go.html", feed)
}
