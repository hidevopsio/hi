package main

import (
	"net/http"
	"sync"
	"testing"

	"github.com/hidevopsio/hiboot/pkg/app/web"
)
var mu sync.Mutex

func TestController(t *testing.T) {
	mu.Lock()
	testApp := web.NewTestApp(t).SetProperty("server.port", "8081").Run(t)

	t.Run("should get index.html ", func(t *testing.T) {
		testApp.Get("/public/ui").
			Expect().Status(http.StatusOK)
	})

	t.Run("should get hello.txt ", func(t *testing.T) {
		testApp.Get("/public/ui/hello.txt").
			Expect().Status(http.StatusOK)
	})
	mu.Unlock()
}


func TestRunMain(t *testing.T) {
	mu.Lock()
	go main()
	mu.Unlock()
}