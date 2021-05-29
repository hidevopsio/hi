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

package web

import (
	"net/http"
	"testing"

	"github.com/iris-contrib/httpexpect"
	"github.com/kataras/iris/httptest"
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/log"
)

// TestApplication the test web application interface for unit test only
type TestApplication interface {
	Initialize() error
	SetProperty(name string, value ...interface{}) TestApplication
	Run(t *testing.T) TestApplication
	Request(method, path string, pathargs ...interface{}) *httpexpect.Request
	Post(path string, pathargs ...interface{}) *httpexpect.Request
	Get(path string, pathargs ...interface{}) *httpexpect.Request
	Put(path string, pathargs ...interface{}) *httpexpect.Request
	Delete(path string, pathargs ...interface{}) *httpexpect.Request
	Patch(path string, pathargs ...interface{}) *httpexpect.Request
	Options(path string, pathargs ...interface{}) *httpexpect.Request
}

// TestApplication the test web application for unit test only
type testApplication struct {
	application
	expect *httpexpect.Expect
}

// RunTestApplication returns the new test application
func RunTestApplication(t *testing.T, controllers ...interface{}) TestApplication {
	log.SetLevel(log.DebugLevel)
	a := new(testApplication)
	err := a.initialize(controllers...)
	assert.Equal(t, nil, err)
	a.Run(t)
	return a
}

// NewTestApplication is the alias of RunTestApplication
// Deprecated, you should use RunTestApplication instead
var NewTestApplication = RunTestApplication

// NewTestApp returns the new test application
func NewTestApp(controllers ...interface{}) TestApplication {
	log.SetLevel(log.DebugLevel)
	a := new(testApplication)
	a.initialize(controllers...)
	return a
}

// SetProperty set application property
func (a *testApplication) SetProperty(name string, value ...interface{}) TestApplication {
	a.BaseApplication.SetProperty(name, value...)
	return a
}

// RunTestServer run the test server
func (a *testApplication) Run(t *testing.T) TestApplication {
	err := a.build()
	if err != nil {
		log.Error(err)
	}
	a.expect = httptest.New(t, a.webApp.Application)
	// Reset log level as httptest.New disabled log as default
	log.SetLevel(a.SystemConfig().Logging.Level)
	return a
}

// Request request for unit test
func (a *testApplication) Request(method, path string, pathargs ...interface{}) *httpexpect.Request {
	return a.expect.Request(method, path, pathargs...)
}

// Post wrap of Request with POST method
func (a *testApplication) Post(path string, pathargs ...interface{}) *httpexpect.Request {
	return a.expect.Request(http.MethodPost, path, pathargs...)
}

// Put wrap of Request with Put method
func (a *testApplication) Put(path string, pathargs ...interface{}) *httpexpect.Request {
	return a.expect.Request(http.MethodPut, path, pathargs...)
}

// Patch wrap of Request with Patch method
func (a *testApplication) Patch(path string, pathargs ...interface{}) *httpexpect.Request {
	return a.expect.Request(http.MethodPatch, path, pathargs...)
}

// Get wrap of Request with Get method
func (a *testApplication) Get(path string, pathargs ...interface{}) *httpexpect.Request {
	return a.expect.Request(http.MethodGet, path, pathargs...)
}

// Delete wrap of Request with Delete method
func (a *testApplication) Delete(path string, pathargs ...interface{}) *httpexpect.Request {
	return a.expect.Request(http.MethodDelete, path, pathargs...)
}

// Delete wrap of Request with Delete method
func (a *testApplication) Options(path string, pathargs ...interface{}) *httpexpect.Request {
	return a.expect.Request(http.MethodOptions, path, pathargs...)
}
