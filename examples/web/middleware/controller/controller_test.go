package controller_test

import (
	_ "hidevops.io/hiboot/examples/web/middleware/controller"
	_ "hidevops.io/hiboot/examples/web/middleware/logging"
	"hidevops.io/hiboot/pkg/app/web"
	"net/http"
	"testing"
)

func TestMiddleware(t *testing.T) {
	testApp := web.NewTestApp().Run(t)

	t.Run("should get user", func(t *testing.T) {
		testApp.Get("/user/123456").
			Expect().Status(http.StatusOK).
			Body().Contains("123456")
	})

	t.Run("should delete user", func(t *testing.T) {
		testApp.Delete("/user/123456").
			Expect().Status(http.StatusOK)
	})
}
