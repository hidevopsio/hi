// Copyright 2018 John Deng (hi.devops.io@gmail.com).
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package system

import "github.com/hidevopsio/hiboot/pkg/at"

// Profiles is app profiles
// .include auto configuration starter should be included inside this slide
// .active active profile
type Profiles struct {
	// included profiles
	Include []string `json:"include,omitempty"`
	// active profile
	Active string `json:"active,omitempty" default:"default"`
}

type banner struct {
	// disable banner
	Disabled bool `json:"disabled" default:"false"`
	Custom string `json:"custom"`
}

type ContactInfo struct {
	Name  string `json:"name,omitempty"`
	URL   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

type License struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

// App is the properties of the application, it hold the base info of the application
type App struct {
	// at.ConfigurationProperties annotation
	at.ConfigurationProperties `value:"app" json:"-"`
	at.AutoWired

	// project name
	Title string `json:"title,omitempty" default:"HiBoot Demo Application"`
	// project name
	Project string `json:"project,omitempty" default:"hidevopsio"`
	// app name
	Name string `json:"name,omitempty" default:"hiboot-app"`
	// app description
	Description string `json:"description,omitempty" default:"${app.name} is a Hiboot Application"`
	// profiles
	Profiles Profiles `json:"profiles"`
	// banner
	Banner banner `json:"banner"`
	// Version
	Version string `json:"version,omitempty" default:"${APP_VERSION:v1}"`
	// TermsOfService
	TermsOfService string       `json:"termsOfService,omitempty"`
	Contact        *ContactInfo `json:"contact,omitempty"`
	License        *License     `json:"license,omitempty"`
}

// Server is the properties of http server
type Server struct {
	// annotation
	at.ConfigurationProperties `value:"server" json:"-"`
	at.AutoWired

	Schemes                    []string `json:"schemes,omitempty" default:"http"`
	Host                       string   `json:"host,omitempty" default:"localhost"`
	Port                       string   `json:"port,omitempty" default:"8080"`
	ContextPath                string   `json:"context_path,omitempty" default:"/"`
}

// Logging is the properties of logging
type Logging struct {
	// annotation
	at.ConfigurationProperties `value:"logging" json:"-"`
	at.AutoWired

	Level string `json:"level,omitempty" default:"info"`
}

