package controller

import (
	"github.com/kataras/iris/mvc"
	"citrix.com/xaxdcloud/common-web-backend/service_profile/model"
	"citrix.com/xaxdcloud/common-web-backend/service_profile/service"
)

// UsersController is our /users API controller.
// GET				/users  | get all
// GET				/users/{id:long} | get by id
// PUT				/users/{id:long} | update by id
// DELETE			/users/{id:long} | delete by id
// Requires basic authentication.
type UsersController struct {
	mvc.C

	Service service.UserService
}

// This could be possible but we should not call handlers inside the `BeginRequest`.
// Because `BeginRequest` was introduced to set common, shared variables between all method handlers
// before their execution.
// We will add this middleware from our `app.Controller` call.
//
// var authMiddleware = basicauth.New(basicauth.Config{
// 	Users: map[string]string{
// 		"admin": "password",
// 	},
// })
//
// func (c *UsersController) BeginRequest(ctx iris.Context) {
//	c.C.BeginRequest(ctx)
//
// 	if !ctx.Proceed(authMiddleware) {
// 		ctx.StopExecution()
// 	}
// }

// Get returns list of the users.
// Demo:
// curl -i -u admin:password http://localhost:8080/users
//
// The correct way if you have sensitive data:
// func (c *UsersController) Get() (results []viewmodels.User) {
// 	data := c.Service.GetAll()
//
// 	for _, user := range data {
// 		results = append(results, viewmodels.User{user})
// 	}
// 	return
// }
// otherwise just return the model.
func (c *UsersController) Get() (results []model.User) {
	return c.Service.GetAll()
}

// GetBy returns a user.
// Demo:
// curl -i -u admin:password http://localhost:8080/users/1
func (c *UsersController) GetBy(id int64) (user model.User, found bool) {
	u, found := c.Service.GetByID(id)
	if !found {
		// this message will be binded to the
		// main.go -> app.OnAnyErrorCode -> NotFound -> shared/error.html -> .Message text.
		c.Ctx.Values().Set("message", "User couldn't be found!")
	}
	return u, found // it will throw/emit 404 if found == false.
}

// PutBy updates a user.
// Demo:
// curl -i -X PUT -u admin:password -F "username=kataras"
// -F "password=rawPasswordIsNotSafeIfOrNotHTTPs_You_Should_Use_A_client_side_lib_for_hash_as_well"
// http://localhost:8080/users/1
func (c *UsersController) PutBy(id int64) (model.User, error) {
	// username := c.Ctx.FormValue("username")
	// password := c.Ctx.FormValue("password")
	u := model.User{}
	if err := c.Ctx.ReadForm(&u); err != nil {
		return u, err
	}

	return c.Service.Update(id, u)
}

// DeleteBy deletes a user.
// Demo:
// curl -i -X DELETE -u admin:password http://localhost:8080/users/1
func (c *UsersController) DeleteBy(id int64) interface{} {
	wasDel := c.Service.DeleteByID(id)
	if wasDel {
		// return the deleted user's ID
		return map[string]interface{}{"deleted": id}
	}
	// right here we can see that a method function
	// can return any of those two types(map or int),
	// we don't have to specify the return type to a specific type.
	return 400 // same as `iris.StatusBadRequest`.
}
