package httpclient

import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/at"
)

const (
	// Profile is a http client profile.
	Profile = "httpclient"
)

type configuration struct {
	at.AutoConfiguration
}

func init() {
	app.Register(newConfiguration)
}

func newConfiguration() *configuration {
	return &configuration{}
}

// client returns an instance of Client
func (c *configuration) Client() Client {
	return NewClient()
}
