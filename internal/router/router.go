package myrouter

import (
	uh "prmeet/internal/handler/user"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(app *echo.Echo) {

	app.GET("", uh.IndexPage)

	// data related routes
	users := app.Group("/users")

	// users.Use(my_utils.Auth)

	users.GET("", uh.IndexPage)
	users.GET("/login", uh.UserLogin)
	users.GET("/logout", uh.UserLogout)
	users.POST("/login", uh.UserDoLogin)
	users.GET("/signup", uh.UserSignup)
	users.POST("/signup", uh.UserDoSignup)

}
