package swagger_test

import (
	"fmt"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/app/web/context"
	"github.com/hidevopsio/hiboot/pkg/app/web/server"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/model"
	"github.com/hidevopsio/hiboot/pkg/starter/swagger"
	"net/http"
	"testing"
	"time"
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
	at.Schema   `json:"-"`
	Id          int     `schema:"The auto generated employee ID" json:"id"`
	FirstName   string  `schema:"The employee first name" json:"first_name" example:"John"`
	LastName    string  `schema:"The employee last name" json:"last_name"`
	Email       string  `schema:"The email of the employee"`
	Address     string  `schema:"The address of the employee"`
	PhoneNumber string  `json:"phone_number"`
	Manger      Manager `schema:"The manager" json:"manger"`
	Assets      []Asset `schema:"The assets list of the employee" json:"assets"`
}

type ErrorResponse struct {
	at.Schema `json:"-"`
	model.BaseResponseInfo
}

type UpdateEmployeeRequest struct {
	at.RequestBody
	Employee
}

type CreateEmployeeRequest struct {
	at.RequestBody
	Employee
}

type EmployeeResponse struct {
	at.ResponseBody `json:"-"`
	at.Schema       `json:"-"`

	model.BaseResponseInfo
	Data *Employee `json:"data,omitempty" schema:"The employee data"`
}

type ListEmployeeResponse struct {
	at.ResponseBody `json:"-"`
	at.Schema       `json:"-"`
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
	log.Debug("before method")
	ctx.Next()
	return
}
type Foo struct {
	at.Schema

	Name string `json:"name"`
	Child *Foo `json:"child"`
	Children []*Foo `json:"children"`
	GradChildren []Foo `json:"grad_children"`
}

// Foo
func (c *employeeController) Foo(at struct {
	at.PostMapping `value:"/foo"`
	at.Operation   `id:"Foo" description:"This is the foo test api"`
	at.Consumes    `values:"application/json"`
	at.Produces    `values:"application/json"`
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
	response.SetData(&Foo{Name: "foo", Child: &Foo{
		Name: "foo1",
	}})

	return
}

// GetEmployee
func (c *employeeController) CreateEmployee(at struct {
	at.PostMapping `value:"/"`
	at.Operation   `id:"Create Employee" description:"This is the employee creation api"`
	at.Consumes    `values:"application/json"`
	at.Produces    `values:"application/json"`
	at.Tags 	   `values:"create,employee,new"`
	at.ExternalDocs `url:"http://hiboot.hidevops.io" description:"HiBoot Official Site"`
	Parameters     struct {
		at.Parameter `name:"employee" in:"body" description:"Employee request body" `
		CreateEmployeeRequest
	}
	Responses struct {
		StatusOK struct {
			at.Response `code:"200" description:"returns a employee with ID"`
			Headers struct{
				XRateLimit struct{
					at.Header `value:"X-Rate-Limit" type:"integer" format:"int32" description:"calls per hour allowed by the user"`
				}
				XExpiresAfter struct {
					at.Header `value:"X-Expires-After" type:"string" format:"date-time" description:"date in UTC when token expires"`
				}
			}
			EmployeeResponse
		}
	}
}, request *CreateEmployeeRequest) (response model.Response, err error) {
	response = new(model.BaseResponse)

	// Just for the demo purpose
	request.Employee.Id = 654321
	response.SetData(request.Employee)

	return
}

// GetEmployee
func (c *employeeController) UpdateEmployee(at struct {
	at.PutMapping `value:"/"`
	at.Operation  `id:"Update Employee" description:"This is the employee update api"`
	at.Consumes   `values:"application/json"`
	at.Produces   `values:"application/json"`
	Parameters    struct {
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
			ErrorResponse
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
			ErrorResponse
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
			ErrorResponse
		}
		StatusNotFound struct {
			at.Response `code:"404" description:"the employee is not found"`
			ErrorResponse
		}
	}
}, id int) (response model.ResponseInfo, err error) {
	response = new(model.BaseResponseInfo)
	return
}

// AddEmployeeAsserts
func (c *employeeController) AddEmployeeAsserts(at struct {
	at.PostMapping `value:"/add-assets"`
	at.Operation   `id:"Add Employee's Assets" description:"This is the api that adding assets for employees"`
	at.Produces    `values:"application/json"`
	Parameters     struct {
		at.Parameter `name:"assets" in:"body" description:"Employee request body" `
		at.Schema    `value:"array" description:"The assets parameter"`
		assets       []*Asset
	}
	Responses struct {
		StatusOK struct {
			at.Response `code:"200" description:"returns a employee with ID"`
			at.Schema   `value:"array" description:"The assets response"`
			Assets      []*Asset
		}
	}
}) (response model.ResponseInfo, err error) {
	response = new(model.BaseResponseInfo)
	return
}

// After
func (c *employeeController) AfterMethod(at struct{ at.AfterMethod }, ctx context.Context) {
	log.Debug("before method")
	ctx.Next()
	return
}

func TestController(t *testing.T) {
	app.Register(
		newEmployeeController,
		swagger.ApiInfoBuilder().
			Title("HiBoot Swagger Demo Application - Simple CRUD Demo Application - 演示代码").
			Description("Simple Server is an application that demonstrate the usage of Swagger Annotations"),
	)
	web.NewTestApp(t).
		SetProperty(server.Schemes, "http,https").
		SetProperty(server.Host, "localhost:8080").
		SetProperty(server.ContextPath, "/v2").
		SetProperty(app.Version, "v2").
		SetProperty(app.ProfilesInclude, web.Profile, swagger.Profile).
		Run(t)

	app.Register(
		swagger.ApiInfoBuilder().
			Version("1.1.0").
			Title("HiBoot Swagger Demo Application - Simple CRUD Demo Application - 演示代码").
			Description("Simple Server is an application that demonstrate the usage of Swagger Annotations").
			Schemes("http").
			Host("localhost:8080").
			BasePath("/").Contact(swagger.Contact{
			Name:  "foo",
			URL:   "http://bar.com",
			Email: "foo@bar.com",
		}).License(swagger.License{
			Name: "foo-lic",
			URL:  "http://bar.com",
		}).TermsOfServiceUrl(`
Copyright 2018 John Deng (hi.devops.io@gmail.com).

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
`),
	)

	time.Sleep(time.Second)
	testApp := web.NewTestApp(t).SetProperty(app.ProfilesInclude, web.Profile, swagger.Profile).Run(t)
	employee := Employee{
		Id:        12345,
		FirstName: "foo",
		LastName:  "bar",
		Manger: Manager{
			ID:   23345,
			Name: "baz",
		},
		Assets: []Asset{
			{
				ID:   1234,
				Name: "abc",
			},
			{
				ID:   5678,
				Name: "def",
			},
		},
	}

	t.Run("should get employee ", func(t *testing.T) {
		testApp.Get("/employee/123").
			Expect().Status(http.StatusOK)
	})


	t.Run("should get employee ", func(t *testing.T) {
		testApp.Get("/employee/123/name").
			Expect().Status(http.StatusOK)
	})



	t.Run("should delete employee ", func(t *testing.T) {
		testApp.Delete("/employee/333").
			Expect().Status(http.StatusOK)
	})

	t.Run("should report 404 when employee does not exist", func(t *testing.T) {
		testApp.Get("/employee/100").
			Expect().Status(http.StatusNotFound)
	})

	t.Run("should list employee", func(t *testing.T) {
		testApp.Get("/employee").
			Expect().Status(http.StatusOK)
	})

	t.Run("should update employee", func(t *testing.T) {
		testApp.Put("/employee").
			WithJSON(&UpdateEmployeeRequest{
				Employee: employee,
			}).Expect().Status(http.StatusOK)
	})
	t.Run("should create employee", func(t *testing.T) {
		testApp.Post("/employee").
			WithJSON(&CreateEmployeeRequest{
				Employee: employee,
			}).Expect().Status(http.StatusOK)
	})

	t.Run("should report 500 error if create employee without request body", func(t *testing.T) {
		testApp.Post("/employee").
			Expect().Status(http.StatusInternalServerError)
	})

	t.Run("should get employees", func(t *testing.T) {
		testApp.Get("/employee").
			Expect().Status(http.StatusOK)
	})

	t.Run("should get swagger-ui", func(t *testing.T) {
		testApp.Get("/swagger-ui").
			Expect().Status(http.StatusOK)
	})

	t.Run("should get swagger.json", func(t *testing.T) {
		testApp.Get("/swagger.json").
			Expect().Status(http.StatusOK)
	})

}
