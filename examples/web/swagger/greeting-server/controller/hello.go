package controller

import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/at"
)

type controller struct {
	at.RestController
	at.RequestMapping `value:"/"`
}

func init() {
	app.Register(newHelloController)
}

func newHelloController() *controller {
	return &controller{}
}

type HelloQueryParam struct {
	at.RequestParams
	at.Schema
	Name        string `schema:"defaults to World if not given"`
}

// Hello
func (c *controller) Hello(at struct {
	at.GetMapping `value:"/hello"`
	at.Operation  `id:"hello" description:"This is the Greeting api for demo"`
	at.Produces   `values:"text/plain"`
	Parameters     struct {
		Name struct{
			at.Parameter `type:"string" name:"name" in:"query" description:"defaults to World if not given" `
		}
	}
	Responses struct {
		StatusOK struct {
			at.Response `code:"200" description:"returns a greeting"`
			at.Schema   `type:"string" description:"contains the actual greeting as plain text"`
		}
		StatusNotFound struct {
			at.Response `code:"404" description:"greeter is not available"`
			at.Schema   `type:"string" description:"Report 'not found' error message"`
		}
	}
	// Response MyResponse
}, param *HelloQueryParam) (response string) {
	if param.Name == "" {
		param.Name = "world"
	}
	response = "Hello, " + param.Name
	return
}


// Hey
func (c *controller) Hey(at struct {
	at.GetMapping `value:"/hey"`
	at.Operation  `id:"hey" description:"This is the another Greeting api for demo"`
	at.Produces   `values:"text/plain"`
	Parameters     struct {
		Name struct{
			at.Parameter `type:"string" name:"name" in:"query" description:"defaults to HiBoot if not given" `
		}
	}
	Responses struct {
		StatusOK struct {
			at.Response `code:"200" description:"returns a greeting"`
			at.Schema   `type:"string" description:"contains the actual greeting as plain text"`
		}
		StatusNotFound struct {
			at.Response `code:"404" description:"greeter is not available"`
			at.Schema   `type:"string" description:"Report 'not found' error message"`
		}
		StatusUnauthorized struct {
			at.Response `code:"401" description:"greeter is not allowed"`
			at.Schema   `type:"string" description:"Report 'not allowed' error message"`
		}
	}
	// Response MyResponse
}, param *HelloQueryParam) (response string) {
	if param.Name == "" {
		param.Name = "HiBoot"
	}
	response = "Hey, " + param.Name
	return
}
