package main

import (
	"net/http"
	"sync"
	"testing"

	"hidevops.io/hiboot/pkg/app/web"
)

func TestController(t *testing.T) {
	testApp := web.NewTestApp(t).Run(t)

	t.Run("should get index.html ", func(t *testing.T) {
		testApp.Get("/public/ui").
			Expect().Status(http.StatusOK)
	})

	t.Run("should get hello.txt ", func(t *testing.T) {
		testApp.Get("/public/ui/hello.txt").
			Expect().Status(http.StatusOK)
	})
}

var mu sync.Mutex
func TestRunMain(t *testing.T) {
	mu.Lock()
	go main()
	mu.Unlock()
}