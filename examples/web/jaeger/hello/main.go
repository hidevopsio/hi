// Copyright 2018 John Deng (hi.devops.io@gmail.com).
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package jaeger provides the quick start jaeger application example
// main package
package main

// import web starter from hiboot
import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/starter/httpclient"
	"github.com/hidevopsio/hiboot/pkg/starter/jaeger"
	"github.com/hidevopsio/hiboot/pkg/starter/logging"
	"io/ioutil"
	"net/http"
)

// Controller Rest Controller with path /
// RESTful Controller, derived from at.RestController. The context mapping of this controller is '/' by default
type Controller struct {
	// at.RestController or at.RestController must be embedded here
	at.RestController

	Formatter string `value:"${provider.formatter}"`
	Publisher string `value:"${provider.publisher}"`

	client httpclient.Client
}

// Get GET /greeting/{greeting}/name/{name}
func (c *Controller) GetByGreetingName(greeting, name string, span *jaeger.Span) string {
	defer span.Finish()
	span.SetTag("hello-to", name)
	span.SetBaggageItem("greeting", greeting)

	helloStr, err := c.formatString(span, name)
	if err != nil {
		return ""
	}
	c.printHello(span, helloStr)

	// response
	return helloStr
}

func (c *Controller) formatString(span *jaeger.Span, helloTo string) (string, error) {

	finalUrl := c.Formatter + "/formatter/" + helloTo

	var newSpan opentracing.Span

	resp, err := c.client.Get(finalUrl, nil, func(req *http.Request) {
		// call formatter service
		newSpan = span.Inject(context.Background(), "GET", finalUrl, req)
	})
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	helloStr := string(body)

	newSpan.LogFields(
		log.String("event", "string-format"),
		log.String("value", helloStr),
	)

	newSpan.Finish()
	return helloStr, nil
}

func (c *Controller) printHello(span *jaeger.Span, helloStr string) {

	finalUrl := c.Publisher + "/publisher/" + helloStr

	var newSpan opentracing.Span

	c.client.Get(finalUrl, nil, func(req *http.Request) {
		// call formatter service
		newSpan = span.Inject(context.Background(), "GET", finalUrl, req)
	})

	newSpan.Finish()

}

func newController(client httpclient.Client) *Controller {
	return &Controller{
		client: client,
	}
}

func init() {
	app.Register(newController)
}

// main function
func main() {
	// create new web application and run it
	web.NewApplication().
		SetProperty(app.ProfilesInclude,
			logging.Profile,
			web.Profile,
			jaeger.Profile,
			httpclient.Profile).
		Run()
}
