package main

import (
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"net/http"
	"sync"
	"testing"
	"time"
)

var mu sync.Mutex
func TestRunMain(t *testing.T) {
	mu.Lock()
	go main()
	mu.Unlock()
}

func TestController(t *testing.T) {
	time.Sleep(time.Second)
	web.NewTestApp().
		SetProperty(web.ViewEnabled, true).
		Run(t).
		Get("/").
		Expect().Status(http.StatusOK)
}
