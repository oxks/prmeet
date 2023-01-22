package router

import (
	"prmeet/internal/handler/user"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(app *echo.Echo) {

	app.GET("", user.IndexPage)

	// data related routes
	users := app.Group("/users")

	// users.Use(my_utils.Auth)

	users.GET("", user.IndexPage)
	users.GET("/login", user.UserLogin)
	users.GET("/logout", user.UserLogout)
	users.POST("/login", user.UserDoLogin)
	users.GET("/signup", user.UserSignup)
	users.POST("/signup", user.UserDoSignup)

}
