package main

import (
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/hidevopsio/hiboot/pkg/app/web"
)

var mu sync.Mutex
func TestRunMain(t *testing.T) {
	mu.Lock()
	go main()
	mu.Unlock()
}

func TestController(t *testing.T) {
	mu.Lock()
	testApp := web.NewTestApp(t).Run(t)

	t.Run("should get employee ", func(t *testing.T) {
		testApp.Get("/employee/123").
			Expect().Status(http.StatusOK)
	})

	t.Run("should get employee ", func(t *testing.T) {
		testApp.Get("/employee/999/name").
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
				Employee: Employee{
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
				},
			}).Expect().Status(http.StatusOK)
	})

	t.Run("should create employee by post /employee", func(t *testing.T) {
		testApp.Post("/employee").
			WithJSON(&CreateEmployeeRequest{
				Employee: Employee{
					Id:        12345,
					FirstName: "foo",
					LastName:  "bar",
					Manger: Manager{
						ID:   23345,
						Name: "baz",
					},
					Assets: []Asset{
						{
							ID:             1234,
							Name:           "abc",
							Amount:         1.234,
							ExpirationTime: time.Now(),
						},
						{
							ID:             5678,
							Name:           "def",
							Amount:         3.14,
							ExpirationTime: time.Now(),
						},
					},
				},
			}).Expect().Status(http.StatusOK).Body().Contains("abc").Contains("def")
	})

	t.Run("should post /foo", func(t *testing.T) {
		testApp.Post("/employee/foo").
			WithJSON(&Foo{
				Name: "foo"}).Expect().Status(http.StatusOK).Body().Contains("foo")
	})

	t.Run("should post /bar", func(t *testing.T) {
		testApp.Post("/employee/bar").
			WithJSON(&Bar{
				Name: "bar"}).Expect().Status(http.StatusOK).Body().Contains("bar")
	})

	t.Run("should add asset", func(t *testing.T) {
		testApp.Post("/employee/add-assets").
			WithJSON([]*Asset{
				{
					ID:             1234,
					Name:           "abc",
					Amount:         1.234,
					Type:			"foo",
					ExpirationTime: time.Now(),
				},
				{
					ID:             5678,
					Name:           "def",
					Type:			"bar",
					Amount:         3.14,
					ExpirationTime: time.Now(),
				},
			}).Expect().Status(http.StatusOK).Body().Contains("abc")
	})

	t.Run("should report 500 error if create employee without request body", func(t *testing.T) {
		testApp.Post("/employee").
			Expect().Status(http.StatusInternalServerError)
	})

	t.Run("should get employees", func(t *testing.T) {
		testApp.Get("/employee").
			Expect().Status(http.StatusOK)
	})
	mu.Unlock()
}
