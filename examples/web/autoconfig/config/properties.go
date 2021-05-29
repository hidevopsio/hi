package config

import "github.com/hidevopsio/hiboot/pkg/at"

type properties struct {
	at.ConfigurationProperties `value:"config"`

	Enabled bool `json:"enabled" `

	Name string `json:"name" default:"foo"`

	DefaultName string `json:"default_name" default:"foo"`
}