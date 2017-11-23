// file: main.go

package main

import (
	"time"

	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"

	"citrix.com/xaxdcloud/common-web-backend/service_profile/datasource"
	"citrix.com/xaxdcloud/common-web-backend/service_profile/repository"
	"citrix.com/xaxdcloud/common-web-backend/service_profile/service"
	"citrix.com/xaxdcloud/common-web-backend/service_profile/web/controller"
	"citrix.com/xaxdcloud/common-web-backend/service_profile/web/middleware"
)

func main() {
	app := iris.New()
	// You got full debug messages, useful when using MVC and you want to make
	// sure that your code is aligned with the Iris' MVC Architecture.
	app.Logger().SetLevel("debug")

	// Load the template files.
	tmpl := iris.HTML("./web/views", ".html").
		Layout("shared/layout.html").
		Reload(true)
	app.RegisterView(tmpl)

	app.StaticWeb("/public", "./web/public")

	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("Message", ctx.Values().
			GetStringDefault("message", "The page you're looking for doesn't exist"))
		ctx.View("shared/error.html")
	})

	// Create our repositories and services.
	db, err := datasource.LoadUsers(datasource.Memory)
	if err != nil {
		app.Logger().Fatalf("error while loading the users: %v", err)
		return
	}
	repo := repository.NewUserRepository(db)
	userService := service.NewUserService(repo)

	// Register our controllers.
	app.Controller("/users", new(controller.UsersController),
		// Add the basic authentication(admin:password) middleware
		// for the /users based requests.
		middleware.BasicAuth,
		// Bind the "userService" to the UserController's Service (interface) field.
		userService,
	)

	sessManager := sessions.New(sessions.Config{
		Cookie:  "sessioncookiename",
		Expires: 24 * time.Hour,
	})
	app.Controller("/user", new(controller.UserController), userService, sessManager)

	// Start the web server at localhost:8080
	// http://localhost:8080/hello
	// http://localhost:8080/hello/iris
	// http://localhost:8080/users/1
	app.Run(
		iris.Addr("localhost:7000"),
		iris.WithoutVersionChecker,
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations, // enables faster json serialization and more
	)
}
