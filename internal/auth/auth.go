package auth

import (
	"errors"
	"fmt"
	"net/http"
	"prmeet/internal/er"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func PasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func PasswordCheckHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// fills in feed["user"] even if user is nil
func LoadAppData(c echo.Context) map[string]any {

	feed := make(map[string]any)
	sess, err := session.Get("session", c)
	er.ErrorPrint(err)
	feed["user"] = sess.Values["user"]
	if sess.Values["err"] != nil {
		feed["err"] = sess.Values["err"]
		sess.Values["err"] = nil
	}
	if sess.Values["success"] != nil {
		feed["success"] = sess.Values["success"]
		sess.Values["success"] = nil
	}
	err = sess.Save(c.Request(), c.Response())

	if err != nil {
		fmt.Println(err)
	}

	return feed
}

// authentication middleware
func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		feed := LoadAppData(c)
		if feed["user"] == nil {
			feed["err"] = errors.New("please login first")
			return c.Render(http.StatusBadRequest, "login.go.html", feed)
		}
		return next(c)
	}
}

func SaveToSession(c echo.Context, key string, data interface{}) bool {
	sess, err := session.Get("session", c)
	er.ErrorPrint(err)
	sess.Values[key] = data
	err = sess.Save(c.Request(), c.Response())
	er.ErrorPrint(err)
	return err == nil
}
