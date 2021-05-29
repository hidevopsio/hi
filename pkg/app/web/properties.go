package web

import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/at"
)

const (
	// ViewEnabled is the property for enabling web view
	ViewEnabled = "web.view.enabled"
	// ContextPath is the property for setting web view context path
	ContextPath = "web.view.contextPath"
	// DefaultPage is the property for setting default page
	DefaultPage = "web.view.defaultPage"
	// ResourcePath is the property for setting resource path
	ResourcePath = "web.view.resourcePath"
	// Extension is the property for setting extension
	Extension = "web.view.extension"
)

type view struct {
	// ViewEnabled is the property for enabling web view
	Enabled bool `json:"enabled"`
	// ContextPath is the property for setting web view context path
	ContextPath string `json:"context_path" default:"/"`
	// DefaultPage is the property for setting default page
	DefaultPage string `json:"default_page" default:"index.html"`
	// ResourcePath is the property for setting resource path
	ResourcePath string `json:"resource_path" default:"./static"`
	// Extension is the property for setting extension
	Extension string `json:"extension" default:".html"`
}

type properties struct {
	at.ConfigurationProperties `value:"web"`
	at.AutoWired

	// View is the properties for setting web view
	View view `json:"view"`
}

func init() {
	app.Register(new(properties))
}
