package controller

import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web/context"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/model"
	"errors"
	"net/http"
)

type UserController struct {
	at.RestController

	at.RequestMapping `value:"/user" `
}

func init() {
	app.Register(newUserController)
}

func newUserController() *UserController {
	return &UserController{}
}

type User struct {
	ID int `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserRequests struct {
	at.RequestParams
	at.Schema

	// For paginated result sets, page of results to retrieve.
	Page int `url:"page,omitempty" json:"page,omitempty" validate:"min=1"`

	// For paginated result sets, the number of results to include per page.
	PerPage int `url:"per_page,omitempty" json:"per_page,omitempty" validate:"min=1"`

	Total int `json:"total"`

	Expr string `json:"expr"`
}

type UserResponse struct {
	at.ResponseBody `json:"-"`
	model.BaseResponse
	Data *User `json:"data"`
}

type ListUserResponse struct {
	at.ResponseBody `json:"-"`
	model.BaseResponse
	Data *UserRequests `json:"data"`
}

// GetUser
func (c *UserController) GetUser(_ struct{
	at.GetMapping `value:"/{id}"`
	at.Operation   `id:"Update Employee" description:"Get User by ID"`
	// /user/{id} -> `values:"user:read" type:"path" in:"id"`
	at.RequiresPermissions `values:"user:read" type:"path" in:"id"`
}, id int, ctx context.Context) (response *UserResponse, err error) {
	log.Debug("GetUserByID requested")
	response = new(UserResponse)
	response.SetCode(http.StatusOK)
	response.SetMessage("Success")
	response.Data = &User{ID: id, Username: "john.deng", Password: "magic-password"}
	err = errors.New("get user error")
	return
}

// GetUser
func (c *UserController) GetUserByName(_ struct{
	at.GetMapping `value:"/name/{name}"`
	// /user/{id} -> `values:"user:read" type:"path" in:"id"`
	at.RequiresPermissions `values:"user:read" type:"path" in:"name"`
}, name string) (response *UserResponse, err error) {
	response = new(UserResponse)
	response.SetCode(http.StatusOK)
	response.SetMessage("Success")
	response.Data = &User{ID: 123456, Username: name, Password: "magic-password"}
	return
}

// GetUser
func (c *UserController) GetUserQuery(_ struct{
	at.GetMapping `value:"/query"`
	at.Operation   `id:"Update Employee" description:"Query User"`
	// /user?id=12345 -> `values:"user:read" type:"query" in:"id"`
	at.RequiresPermissions `values:"user:read" type:"query" in:"id"`
}, ctx context.Context) (response *UserResponse, err error) {
	response = new(UserResponse)
	response.SetCode(http.StatusOK)
	response.SetMessage("Success")
	id, _ := ctx.URLParamInt("id")
	response.Data = &User{ID: id, Username: "john.deng", Password: "magic-password"}
	return
}


// GetUser
func (c *UserController) GetUsers(_ struct{
	at.GetMapping `value:"/"`
	at.Operation   `id:"Update Employee" description:"Get User List"`
	at.RequiresPermissions `values:"user:list" type:"query:pagination" in:"page,per_page,id" out:"expr,total"`
}, request *UserRequests, ctx context.Context) (response *ListUserResponse, err error) {
	log.Debugf("expr: %v", request.Expr)
	log.Debugf("total: %v", request.Total)
	response = new(ListUserResponse)
	response.SetCode(http.StatusOK)
	response.SetMessage("Success")
	response.Data = request
	return
}


// GetUser
func (c *UserController) DeleteUser(_ struct{
	at.DeleteMapping `value:"/{id}"`
	at.Operation   `id:"Update Employee" description:"Delete User by ID"`
}, id int) (response *UserResponse, err error) {
	response = new(UserResponse)
	response.SetCode(http.StatusOK)
	response.SetMessage("Success")
	return
}




