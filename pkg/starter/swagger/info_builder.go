// Copyright 2018~now John Deng (hi.devops.io@gmail.com).
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

package swagger

import (
	"github.com/go-openapi/spec"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/system"
)

type ApiInfoBuilderInterface interface {
	Title(value string) ApiInfoBuilderInterface
	Description(value string) ApiInfoBuilderInterface
	TermsOfServiceUrl(value string) ApiInfoBuilderInterface
	Version(value string) ApiInfoBuilderInterface
	ContactName(value string) ApiInfoBuilderInterface
	ContactEmail(value string) ApiInfoBuilderInterface
	ContactURL(value string) ApiInfoBuilderInterface
	Contact(value Contact) ApiInfoBuilderInterface
	LicenseName(value string) ApiInfoBuilderInterface
	LicenseURL(value string) ApiInfoBuilderInterface
	License(value License) ApiInfoBuilderInterface
	Schemes(values ...string) ApiInfoBuilderInterface
	Host(values string) ApiInfoBuilderInterface
	BasePath(values string) ApiInfoBuilderInterface
	Redoc(value string) ApiInfoBuilderInterface
}

type Contact struct {
	Name  string `json:"name,omitempty"`
	URL   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

type License struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

type apiInfoBuilder struct {
	at.ConfigurationProperties `value:"swagger"`
	at.AutoWired
	spec.Swagger

	SystemServer *system.Server
	AppVersion   string `value:"${app.version}"`
	RedocURL string `json:"redoc_url" value:"${redoc.url:https://cdn.jsdelivr.net/npm/redoc@next/bundles/redoc.standalone.js}"`
}

func ApiInfoBuilder() ApiInfoBuilderInterface {
	return &apiInfoBuilder{
		Swagger: spec.Swagger{
			SwaggerProps: spec.SwaggerProps{
				Swagger: "2.0",
				Info:    &spec.Info{},
				Paths: &spec.Paths{
					VendorExtensible: spec.VendorExtensible{},
					Paths:            make(map[string]spec.PathItem),
				},
			},
		},
	}
}

func (b *apiInfoBuilder) Title(value string) ApiInfoBuilderInterface {
	b.Swagger.SwaggerProps.Info.Title = value
	return b
}

func (b *apiInfoBuilder) Description(value string) ApiInfoBuilderInterface {
	b.Swagger.SwaggerProps.Info.Description = value
	return b
}

func (b *apiInfoBuilder) TermsOfServiceUrl(value string) ApiInfoBuilderInterface {
	b.Swagger.SwaggerProps.Info.TermsOfService = value
	return b
}

func (b *apiInfoBuilder) Version(value string) ApiInfoBuilderInterface {
	b.Swagger.SwaggerProps.Info.Version = value
	return b
}

func (b *apiInfoBuilder) Host(values string) ApiInfoBuilderInterface {
	b.Swagger.SwaggerProps.Host = values
	return b
}

func (b *apiInfoBuilder) BasePath(values string) ApiInfoBuilderInterface {
	b.Swagger.SwaggerProps.BasePath = values
	return b
}

func (b *apiInfoBuilder) Schemes(values ...string) ApiInfoBuilderInterface {
	b.Swagger.SwaggerProps.Schemes = values
	return b
}


func (b *apiInfoBuilder) ContactName(value string) ApiInfoBuilderInterface {
	b.ensureContact()
	b.Swagger.SwaggerProps.Info.Contact.Name = value
	return b
}

func (b *apiInfoBuilder) ensureContact() {
	if b.Swagger.SwaggerProps.Info.Contact == nil {
		b.Swagger.SwaggerProps.Info.Contact = &spec.ContactInfo{}
	}
}

func (b *apiInfoBuilder) ContactEmail(value string) ApiInfoBuilderInterface {
	b.ensureContact()
	b.Swagger.SwaggerProps.Info.Contact.Email = value
	return b
}

func (b *apiInfoBuilder) ContactURL(value string) ApiInfoBuilderInterface {
	b.ensureContact()
	b.Swagger.SwaggerProps.Info.Contact.URL = value
	return b
}

func (b *apiInfoBuilder) Contact(value Contact) ApiInfoBuilderInterface {
	b.ContactName(value.Name)
	b.ContactEmail(value.Email)
	b.ContactURL(value.URL)
	return b
}

func (b *apiInfoBuilder) ensureLicense() {
	if b.Swagger.SwaggerProps.Info.License == nil {
		b.Swagger.SwaggerProps.Info.License = &spec.License{}
	}
}

func (b *apiInfoBuilder) LicenseName(value string) ApiInfoBuilderInterface {
	b.ensureLicense()
	b.Swagger.SwaggerProps.Info.License.Name = value
	return b
}

func (b *apiInfoBuilder) LicenseURL(value string) ApiInfoBuilderInterface {
	b.ensureLicense()
	b.Swagger.SwaggerProps.Info.License.URL = value
	return b
}

func (b *apiInfoBuilder) License(value License) ApiInfoBuilderInterface {
	b.LicenseName(value.Name)
	b.LicenseURL(value.URL)
	return b
}

func (b *apiInfoBuilder) Redoc(value string) ApiInfoBuilderInterface {
	b.RedocURL = value
	return b
}
