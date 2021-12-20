package my_utils

import (
	"errors"
	"log"
	"net/http"
	me "prmeet/internal/utils/my_errors"
	"time"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"golang.org/x/crypto/bcrypt"
)

// type to use with templates data
type M map[string]interface{}

func AskPipe(subj string, request M) M {

	// Connect to NATS encoded
	enc := ConnectNatsEncoded()
	defer enc.Close()

	response := M{}
	err := enc.Request(subj, request, &response, 2*time.Second)
	me.ErrorPrint(err)

	return response
}

func PasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func PasswordCheckHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// fills in feed["user"] even if user is nil
func LoadAppData(c echo.Context) M {
	feed := M{}
	sess, err := session.Get("session", c)
	me.ErrorPrint(err)
	feed["user"] = sess.Values["user"]
	if sess.Values["err"] != nil {
		feed["err"] = sess.Values["err"]
		sess.Values["err"] = nil
	}
	if sess.Values["success"] != nil {
		feed["success"] = sess.Values["success"]
		sess.Values["success"] = nil
	}
	defer sess.Save(c.Request(), c.Response())

	return feed
}

// authentication middleware
func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		feed := LoadAppData(c)
		if feed["user"] == nil {
			feed["err"] = errors.New("Please login first.")
			return c.Render(http.StatusBadRequest, "login.go.html", feed)
		}
		return next(c)
	}
}

func SaveToSession(c echo.Context, key string, data interface{}) bool {
	sess, err := session.Get("session", c)
	me.ErrorPrint(err)
	sess.Values[key] = data
	err = sess.Save(c.Request(), c.Response())
	me.ErrorPrint(err)
	return err == nil
}

// connect encoded nats
func ConnectNatsEncoded() *nats.EncodedConn {

	nc := ConnectNats()
	enc, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Fatal(err)
	}
	return enc
}
