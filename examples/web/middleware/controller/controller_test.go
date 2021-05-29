package controller_test

import (
	_ "github.com/hidevopsio/hiboot/examples/web/middleware/controller"
	_ "github.com/hidevopsio/hiboot/examples/web/middleware/logging"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"net/http"
	"testing"
)

func TestMiddleware(t *testing.T) {
	testApp := web.NewTestApp().Run(t)

	t.Run("should get user by id", func(t *testing.T) {
		testApp.Get("/user/123456").
			Expect().Status(http.StatusOK).
			Body().Contains("123456")
	})

	t.Run("should get user by name", func(t *testing.T) {
		testApp.Get("/user/name/john.deng").
			Expect().Status(http.StatusOK).
			Body().Contains("john.deng")
	})

	t.Run("should get user", func(t *testing.T) {
		testApp.Get("/user/query").
			WithQuery("id", 123456).
			Expect().Status(http.StatusOK).
			Body().Contains("123456")
	})

	t.Run("should get users", func(t *testing.T) {
		testApp.Get("/user").
			WithQuery("page", 1).
			WithQuery("per_page", 10).
			Expect().Status(http.StatusOK)
	})


	t.Run("should get users without page, per_page params", func(t *testing.T) {
		testApp.Get("/user").
			Expect().Status(http.StatusBadRequest)
	})

	t.Run("should delete user", func(t *testing.T) {
		testApp.Delete("/user/123456").
			Expect().Status(http.StatusOK)
	})
}
