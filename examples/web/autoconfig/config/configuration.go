package config

import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web/context"
	"github.com/hidevopsio/hiboot/pkg/at"
	"net/http"
)

const Profile string = "config"

type configuration struct {
	at.AutoConfiguration

	properties *properties
}

func newConfiguration(properties *properties) *configuration  {
	return &configuration{properties: properties}
}

func init() {
	app.Register(newConfiguration)
}

type Foo struct {
	Name string `json:"name" value:"foo"`
}

type Bar struct {
	at.ContextAware

	Name string `json:"name" value:"bar"`
}

func (c *configuration) Foo() *Foo {
	return &Foo{}
}

func (c *configuration) Bar(ctx context.Context) *Bar {
	if ctx.GetHeader("Authorization") == "fake" {
		ctx.StatusCode(http.StatusUnauthorized)
		return nil
	}
	return &Bar{Name: c.properties.Name}
}