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

import (
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils/io"
	"os"
	"path/filepath"
	"testing"
)

func init() {
	io.EnsureWorkDir(1, "config/application.yml")
	log.SetLevel(log.DebugLevel)
	customProperties["app.project"] = "system-test"
}

var customProperties = make(map[string]interface{})

type StarProperties struct {
	at.ConfigurationProperties `value:"star"`
	at.AutoWired

	Name string `json:"name"`
	System string `json:"system"`
}

func TestPropertyBuilderBuild(t *testing.T) {
	os.Args = append(os.Args, "--logging.level=debug", "--foo.bar",  "--foobar=foo,bar", "--foobar=foo,bar")

	err := os.Setenv("APP_NAME", "hiboot-app")

	testProject := "system-configuration-test"
	customProperties["app.project"] = testProject
	b := NewPropertyBuilder(filepath.Join(io.GetWorkDir(), "config"), customProperties)

	// dummy
	_ = b.Init()
	_, _ = b.BuildWithProfile("")
	_ = b.Save(nil)
	b.SetConfiguration(nil)

	profile := "local" //os.Getenv("APP_PROFILES_ACTIVE")
	_, err = b.Build()

	_, err = b.Build(profile)
	assert.Equal(t, nil, err)

	appProp := &App{}
	t.Run("should build configuration properly", func(t *testing.T) {
		projectName := "hidevopsio"
		b.SetProperty("app.project", projectName)
		appName := b.GetProperty("app.name")
		err = b.Load(appProp)
		assert.Equal(t, nil, err)
		assert.Equal(t, appName, appProp.Name)
		assert.Equal(t, projectName, appProp.Project)
	})

	t.Run("should load properties", func(t *testing.T) {
		star := &StarProperties{}
		err = b.Load(star)
		assert.Equal(t, nil, err)
		assert.Equal(t, "Mars", star.Name)
		assert.Equal(t, "solar", star.System)
	})

	t.Run("should build mock and local configuration properly and merge them", func(t *testing.T) {
		assert.Equal(t, nil, err)
		assert.Equal(t, "hiboot-app-mocking", b.GetProperty("mock.nickname"))
		assert.Equal(t, "hiboot-app-user-local", b.GetProperty("mock.username"))
		assert.Equal(t, "hiboot-mock-local", b.GetProperty("mock.name"))
	})

	t.Run("should build configuration properly", func(t *testing.T) {
		assert.Equal(t, nil, err)
		assert.Equal(t, "hiboot-app", appProp.Name)
		log.Debugf("app: %v", b.GetProperty("app"))
		log.Debugf("app.name: %v", b.GetProperty("app.name"))
		log.Debugf("server: %v", b.GetProperty("server"))
		log.Debugf("server port: %v", b.GetProperty("server.port"))
	})

	t.Run("should build fake configuration", func(t *testing.T) {
		err = os.Setenv("APP_URL", "https://examples.org/api")
		assert.Equal(t, nil, err)
		assert.Equal(t, nil, err)
		assert.Equal(t, "hiboot-app", b.GetProperty("app.name"))

		log.Info(b.GetProperty("fake.name"))
		log.Info(b.GetProperty("app.name"))
		log.Info(b.GetProperty("app.url"))
		log.Info(b.GetProperty("app.profiles.active"))

	})

	t.Run("should set default property", func(t *testing.T) {
		b.SetDefaultProperty("plant.fruit", "apple").
			SetDefaultProperty("animal.mammal.carnivores", "tiger")

		assert.Equal(t, "apple", b.GetProperty("plant.fruit"))
		assert.Equal(t, "tiger", b.GetProperty("animal.mammal.carnivores"))
	})

	t.Run("should overwrite default property", func(t *testing.T) {
		b.SetProperty("plant.fruit", "banana")
		assert.Equal(t, "banana", b.GetProperty("plant.fruit"))
	})

	//t.Run("should overwrite with struct property", func(t *testing.T) {
	//	b.SetProperty("server", Server{
	//		Schemes:     []string{"https"},
	//		Host:        "test.hidevops.io",
	//		Port:        "9090",
	//		ContextPath: "/api/foo/bar",
	//	})
	//	assert.Equal(t, "/api/foo/bar", b.GetProperty("server.context_path"))
	//})

	//t.Run("should set struct property directly", func(t *testing.T) {
	//	b.SetProperty("server2", Server{
	//		Schemes:     []string{"https"},
	//		Host:        "test.hidevops.io",
	//		Port:        "9090",
	//		ContextPath: "/api/foo/bar",
	//	})
	//	assert.Equal(t, "/api/foo/bar", b.GetProperty("server2.context_path"))
	//})

	t.Run("should replace property", func(t *testing.T) {
		b.SetProperty("app.name", "foo").
			SetProperty("app.profiles.include", []string{"foo", "bar"})

		res := b.Replace("this is ${app.name} property")
		assert.Equal(t, "this is foo property", res)

		res = b.Replace("${app.profiles.include}")
		assert.Equal(t, []string{"foo", "bar"}, res)
	})

	t.Run("should replace with default property", func(t *testing.T) {
		res := b.Replace("this is ${default.value:default} property")
		assert.Equal(t, "this is default property", res)
	})

	t.Run("should replace with environment variable", func(t *testing.T) {
		res := b.Replace("this is ${HOME}")
		home := os.Getenv("HOME")
		assert.Equal(t, "this is "+home, res)
	})
}
