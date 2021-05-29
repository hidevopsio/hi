package main

import (
	"embed"
	"fmt"
	"net/http"
	"time"

	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/app/web/context"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/model"
	"github.com/hidevopsio/hiboot/pkg/starter/actuator"
	"github.com/hidevopsio/hiboot/pkg/starter/logging"
	"github.com/hidevopsio/hiboot/pkg/starter/swagger"
)

type Asset struct {
	ID         int    `schema:"The asset ID" json:"id" example:"1234567890"`
	Name       string `schema:"The asset name" json:"name" example:"John Deng"`
	Amount     float64 `json:"amount" example:"987654321"`
	Type       string `schema:"The asset type" json:"type" example:"book"`
	ExpirationTime time.Time `json:"expiration_time" example:"Sun Sep 29 15:47:50 CST 2019"`
}

type Manager struct {
	ID int `schema:"The manager ID" json:"id" default:"1000000"`
	Name string `schema:"The manager name of the employee" json:"name" example:"John Deng"`
}

type Employee struct {
	Id        int     `schema:"The auto generated employee ID" json:"id" example:"123456"`
	FirstName string  `schema:"The employee first name" json:"first_name" example:"John"`
	LastName  string  `schema:"The employee last name" json:"last_name" example:"Deng"`
	Manger    Manager `schema:"The manager" json:"manger"`
	Assets    []Asset `schema:"The assets list of the employee" json:"assets"`
}

type ResponseInfo struct {
	at.Schema
	model.BaseResponseInfo
}


type UpdateEmployeeRequest struct {
	at.RequestBody
	at.Schema
	Employee
}

type CreateEmployeeRequest struct {
	at.RequestBody
	at.Schema
	Employee
}

type EmployeeResponse struct {
	at.ResponseBody `json:"-"`
	at.Schema

	model.BaseResponseInfo
	Data *Employee `json:"data,omitempty" schema:"The employee data"`
}

type ListEmployeeResponse struct {
	at.ResponseBody `json:"-"`
	at.Schema
	model.BaseResponse
	Data []*Employee `json:"data,omitempty" schema:"The employee data list"`
}

type employeeController struct {
	at.RestController
	at.RequestMapping `value:"/employee"`

	dummyData []*Employee
}

// newEmployeeController is the constructor for orgController
// you may inject dependency to constructor
func newEmployeeController() *employeeController {
	return &employeeController{
		dummyData: []*Employee{
			{
				Id:        123,
				FirstName: "John",
				LastName:  "Deng",
			},
			{
				Id:        456,
				FirstName: "Mike",
				LastName:  "Philip",
			},
		},
	}
}

// Before
func (c *employeeController) BeforeMethod(at struct{ at.BeforeMethod }, ctx context.Context) {
	log.Debugf("%v %v before method", ctx.GetCurrentRoute().Method(), ctx.GetCurrentRoute().Path())
	ctx.Next()
	return
}


type Foo struct {
	at.RequestBody
	at.Schema

	Name string `json:"name"`
	Child *Foo `json:"child"`
	Children []*Foo `json:"children"`
	GradChildren []Foo `json:"grad_children"`
}

type Bar struct {
	at.RequestBody
	at.Schema

	Name string `json:"name"`

}

// Foo
func (c *employeeController) Foo(at struct {
	at.PostMapping `value:"/foo"`
	at.Operation   `id:"Foo" description:"This is the foo test api"`
	at.Consumes    `values:"application/json"`
	at.Produces    `values:"application/json"`
	at.Tags 	   `values:"foo,bar,test"`
	Parameters     struct {
		at.Parameter `name:"foo" in:"body" description:"foo request body" `
		Foo
	}
	Responses struct {
		StatusOK struct {
			at.Response `code:"200" description:"returns foo"`
			Foo
		}
	}
}, request *Foo) (response model.Response, err error) {
	response = new(model.BaseResponse)

	// Just for the demo purpose
	response.SetData(&Foo{Name: request.Name, Child: &Foo{
		Name: "foo1",
	}})

	return
}

// Foo
func (c *employeeController) Bar(at struct {
	at.PostMapping `value:"/bar"`
	at.Operation   `id:"Bar" description:"This is the bar test api" deprecated:"true"`
	at.Consumes    `values:"application/json"`
	at.Produces    `values:"application/json"`
	Parameters     struct {
		at.Parameter `name:"foo" in:"body" description:"foo request body" `
		Bar
	}
	Responses struct {
		StatusOK struct {
			at.Response `code:"200" description:"returns foo"`
			Bar
		}
	}
}, request *Bar) (response model.Response, err error) {
	response = new(model.BaseResponse)

	// Just for the demo purpose
	response.SetData(&Bar{Name: request.Name})

	return
}

// GetEmployee
func (c *employeeController) CreateEmployee(at struct {
	at.PostMapping `value:"/"`
	at.Operation   `id:"Create Employee" description:"This is the employee creation api"`
	at.Consumes    `values:"application/json"`
	at.Produces    `values:"application/json"`
	at.ExternalDocs `url:"http://hiboot.hidevops.io" description:"HiBoot Official Site"`
	Parameters     struct {
		at.Parameter `name:"token" in:"header" type:"string" description:"JWT token (fake token - for demo only)" `
		Body struct {
			at.Parameter `name:"employee" in:"body" description:"Employee request body" `
			CreateEmployeeRequest
		}
	}
	Responses struct {
		StatusOK struct {
			at.Response `code:"200" description:"returns a employee with ID"`
			XRateLimit struct {
				at.Header `value:"X-Rate-Limit" type:"integer" format:"int32" description:"calls per hour allowed by the user"`
			}
			XExpiresAfter struct{
				at.Header `value:"X-Expires-After" type:"string" format:"date-time" description:"date in UTC when token expires"`
			}
			EmployeeResponse
		}
	}
}, request *CreateEmployeeRequest) (response *EmployeeResponse, err error) {
	response = new(EmployeeResponse)

	// Just for the demo purpose
	request.Employee.Id = 654321
	response.SetCode(http.StatusOK)
	response.Data = &request.Employee

	return
}


// GetEmployee
func (c *employeeController) UpdateEmployee(at struct {
	at.PutMapping `value:"/"`
	at.Operation   `id:"Update Employee" description:"This is the employee update api"`
	at.Consumes    `values:"application/json"`
	at.Produces    `values:"application/json"`
	Parameters     struct {
		at.Parameter `name:"employee" in:"body" description:"Employee request body" `
		UpdateEmployeeRequest
	}
	Responses struct {
		StatusOK struct {
			at.Response `code:"200" description:"returns a employee with ID"`
			EmployeeResponse
		}
	}
}, request *UpdateEmployeeRequest) (response model.Response, err error) {
	response = new(model.BaseResponse)

	// Just for the demo purpose
	request.Employee.Id = 67890
	response.SetData(request.Employee)

	return
}

// GetEmployee
func (c *employeeController) GetEmployee(at struct {
	at.GetMapping `value:"/{id}"`
	at.Operation  `id:"Get Employee" description:"This is get employees api"`
	at.Produces   `values:"application/json"`
	Parameters    struct {
		ID struct {
			at.Parameter `type:"integer" name:"id" in:"path" description:"Path variable employee ID" required:"true"`
		}
	}
	Responses struct {
		StatusOK struct {
			at.Response `code:"200" description:"returns a employee"`
			EmployeeResponse
		}
		StatusNotFound struct {
			at.Response `code:"404" description:"the employee you are looking for is not found"`
			ResponseInfo
		}
	}
}, id int) (response *EmployeeResponse, err error) {
	response = new(EmployeeResponse)
	for _, e := range c.dummyData {
		if id == e.Id {
			response.Code = http.StatusOK
			response.Message = "success"
			response.Data = e
			break
		}
	}

	if response.Data == nil {
		errMsg := fmt.Sprintf("employee %v is not found", id)
		response.Code = http.StatusNotFound
		response.Message = errMsg
	}

	return
}


// GetEmployeeName
func (c *employeeController) GetEmployeeName(at struct {
	at.GetMapping `value:"/{id}/name"`
	at.Operation  `id:"Get Employee Name" description:"This is the api that get employee name"`
	at.Produces   `values:"text/plain"`
	Parameters    struct {
		ID struct {
			at.Parameter `type:"integer" name:"id" in:"path" description:"Path variable employee ID" required:"true"`
		}
	}
	Responses struct {
		StatusOK struct {
			at.Response `code:"200" description:"returns the employee name"`
			at.Schema   `type:"string" description:"contains the actual employee name as plain text"`
		}
		StatusNotFound struct {
			at.Response `code:"404" description:"employee is not found"`
			at.Schema   `type:"string" description:"Report 'not found' error message"`
		}
	}
}, id int) (name string) {
	return "Donald Trump"
}


// ListEmployee
func (c *employeeController) ListEmployee(at struct {
	at.GetMapping `value:"/"`
	at.Operation  `id:"List Employee" description:"This is employees list api"`
	at.Produces   `values:"application/json"`
	Responses     struct {
		StatusOK struct {
			at.Response `code:"200" description:"returns a set of employees"`
			ListEmployeeResponse
		}
		StatusNotFound struct {
			at.Response `code:"404" description:"the employees you are looking for is not found"`
			ResponseInfo
		}
	}
}) (response *ListEmployeeResponse, err error) {
	response = new(ListEmployeeResponse)
	response.Code = http.StatusOK
	response.Message = "success"
	response.Data = c.dummyData
	return
}

// DeleteEmployee
// at.DeleteEmployee is an annotation to define request mapping for http method DELETE,
func (c *employeeController) DeleteEmployee(at struct {
	at.DeleteMapping `value:"/{id}"`
	at.Operation     `id:"Delete Employee" description:"This is delete employees api"`
	at.Produces      `values:"application/json"`
	Parameters       struct {
		at.Parameter `type:"integer" name:"id" in:"path" description:"Path variable employee ID" required:"true"`
	}
	Responses struct {
		StatusOK struct {
			at.Response `code:"200" description:"returns success message"`
			ResponseInfo
		}
		StatusNotFound struct {
			at.Response `code:"404" description:"the employee is not found"`
			ResponseInfo
		}
	}
}, id int) (response model.ResponseInfo, err error) {
	response = new(model.BaseResponseInfo)
	return
}

// AddEmployeeAsserts
func (c *employeeController) AddEmployeeAsserts(at struct {
	at.PostMapping `value:"/add-assets"`
	at.Operation     `id:"Add Employee's Assets" description:"This is the api that adding assets for employees"`
	at.Consumes      `values:"application/json"`
	at.Produces      `values:"application/json"`
	Parameters    struct {
		at.Parameter `in:"body" description:"Employee request body" `
		at.Schema
		Assets []*Asset
	}
	Responses struct {
		StatusOK struct {
			at.Response `code:"200" description:"returns a employee with ID"`
			at.Schema
			Assets []*Asset
		}
	}
}, ctx context.Context) (response model.Response, err error) {
	var assets []*Asset
	_ = ctx.ReadJSON(&assets)
	response = new(model.BaseResponse)
	response.SetData(assets)
	return
}

// After
func (c *employeeController) AfterMethod(at struct{ at.AfterMethod }, ctx context.Context) {
	log.Debugf("%v %v after method", ctx.GetCurrentRoute().Method(), ctx.GetCurrentRoute().Path())
	ctx.Next()
	return
}

func init() {
	app.Register(
		newEmployeeController,
		swagger.ApiInfoBuilder().
			// Schemes, Version, Host, BasePath use default settings from server properties instead
			// they will be replaced with program arguments as highest priority,
			// for example: my-app --app.version=v.1.0.0  --server.schemes=http, --server.host=example.com --server.context_path=/api/v1/my-app
			//Schemes("http").
			//Version("1.1.0").
			//Host("localhost:8080").
			//BasePath("/").
			ContactName("John Deng").
			ContactEmail("john.deng@outlook.com").
			ContactURL("https://hiboot.hidevops.io").
			Title("HiBoot Swagger Demo Application - Simple CRUD Demo Application - 演示代码").
			Description(`
## Description

  This is an application that demonstrate the usage of Swagger2.0 Annotations, see [example](https://github.com/hidevopsio/hiboot/tree/v2/examples/web/swagger/simple-server) for more details

  这是一个使用HiBoot的注解实现的，支持简单服务在线接口文档，遵循Swagger2.0规范，欲了解更多详情，请看这个[例子](https://github.com/hidevopsio/hiboot/tree/v2/examples/web/swagger/simple-server)

`),
	)
}

//go:embed config/*.yml
var embedFS embed.FS

// Hiboot main function
func main() {
	// create new web application and run it
	web.NewApplication().
		SetProperty(app.Config, &embedFS).
		SetProperty(app.ProfilesInclude,
			actuator.Profile,
			swagger.Profile,
			logging.Profile,
		).
		Run()
}
