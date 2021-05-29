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

type profiles struct {
	Include []string `json:"include"`
	Active  string   `json:"active"`
}

type properties struct {
	at.ConfigurationProperties `value:"app"`
	at.AutoWired

	Name string `json:"name"`
	Profiles profiles `json:"profiles"`
}

type fakeConfiguration struct {
	Properties properties `mapstructure:"app"`
}

func init() {
	io.EnsureWorkDir(1, "config/application.yml")
	log.SetLevel(log.DebugLevel)
	customProps["app.project"] = "system-test"
}

var customProps = make(map[string]interface{})


func TestBuilderBuild(t *testing.T) {
	err := os.Setenv("APP_NAME", "hiboot-app")

	testProject := "hidevopsio"
	customProps["app.project"] = testProject
	b := NewBuilder(NewConfiguration(),
		filepath.Join(io.GetWorkDir(), "config"),
		"application",
		"yaml",
		customProps)

	// not used
	_ = b.Load(nil)

	cp, err := b.Build("default", "local")
	t.Run("should build configuration properly", func(t *testing.T) {
		c := cp.(*Configuration)
		assert.Equal(t, nil, err)
		assert.Equal(t, "hiboot-app", c.App.Name)
		assert.Equal(t, testProject, c.App.Project)
	})

	t.Run("should build mock and local configuration properly and merge them", func(t *testing.T) {
		cp, err = b.Build("mock", "local")
		assert.Equal(t, nil, err)
		assert.Equal(t, "hiboot-app-mocking", b.GetProperty("mock.nickname"))
		assert.Equal(t, "hiboot-app-user-local", b.GetProperty("mock.username"))
		assert.Equal(t, "hiboot-mock-local", b.GetProperty("mock.name"))
	})

	t.Run("should build configuration properly", func(t *testing.T) {
		assert.Equal(t, nil, err)
		c := cp.(*Configuration)
		assert.Equal(t, "hiboot-app", c.App.Name)
		log.Debugf("app: %v", b.GetProperty("app"))
		log.Debugf("app.name: %v", b.GetProperty("app.name"))
		log.Debugf("server: %v", b.GetProperty("server"))
		log.Debugf("server port: %v", b.GetProperty("server.port"))
	})

	t.Run("should build fake configuration", func(t *testing.T) {
		/*
				# config/application.yml
			    # fake.name is a reference of app.name
				app:
				  project: hidevopsio
				  name: hiboot
				  profiles:
					include:
					- foo

				logging:
				  level: debug

				# added for test only
				fake:
				  name: ${app.name}
		*/

		err = os.Setenv("APP_URL", "https://examples.org/api")
		assert.Equal(t, nil, err)
		fc := new(fakeConfiguration)
		b.SetConfiguration(fc)
		br, err := b.Build("fake")
		assert.Equal(t, nil, err)
		assert.NotEqual(t, nil, br)
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

func TestBuilderBuildWithError(t *testing.T) {

	b := NewBuilder(nil, "", "", "", nil)

	_, err := b.Build()
	assert.Equal(t, nil, err)

}

func TestBuilderBuildWithProfile(t *testing.T) {
	customProps["app.project"] = "local-test"
	b := NewBuilder(&Configuration{},
		filepath.Join(io.GetWorkDir(), "config"),
		"application",
		"yaml",
		customProps)

	cp, err := b.BuildWithProfile("local")
	assert.Equal(t, nil, err)

	c := cp.(*Configuration)
	assert.Equal(t, "8081", c.Server.Port)
	log.Print(c)

	_, err = b.BuildWithProfile("")
	assert.Equal(t, nil, err)

}

func TestFileDoesNotExist(t *testing.T) {

	b := NewBuilder(&Configuration{},
		filepath.Join(io.GetWorkDir(), "config"),
		"application",
		"yaml",
		customProps)
	t.Run("use default profile if custom profile does not exist", func(t *testing.T) {
		_, err := b.Build("does-not-exist")
		assert.Equal(t, nil, err)
	})
}

func TestWrongFileFormat(t *testing.T) {

	b := NewBuilder(&Configuration{},
		filepath.Join(os.TempDir(), "config"),
		"test",
		"yaml",
		customProps)
	io.CreateFile(os.TempDir(), "test-abc.yml")
	io.WriterFile(os.TempDir(), "test-abc.yml", []byte(": 1234"))
	t.Run("should report error: did not find expected key", func(t *testing.T) {
		_, err := b.Build("abc")
		assert.Equal(t, nil, err)
	})
	io.WriterFile(os.TempDir(), "test-abc.yml", []byte("abc:"))
	t.Run("use default profile if custom profile does not exist", func(t *testing.T) {
		_, err := b.Build("default")
		assert.Equal(t, nil, err)
	})
}

func TestDefaultProfileOnly(t *testing.T) {
	type emptyConfig struct {
	}
	b := NewBuilder(emptyConfig{},
		filepath.Join(io.GetWorkDir(), "config"),
		"application",
		"yaml",
		customProps)

	t.Run("use default profile if custom profile does not exist", func(t *testing.T) {
		_, err := b.Build("default")
		assert.Equal(t, nil, err)
	})
}

func TestWithoutReplacer(t *testing.T) {

	path := filepath.Join(io.GetWorkDir(), "config")
	testProfile := "xxx"
	appConfig := "application"
	testFile := appConfig + "-" + testProfile + ".yml"
	b := NewBuilder(&Configuration{},
		path,
		"application",
		"yaml",
		customProps)
	io.CreateFile(path, testFile)
	_, err := b.Build("xxx")
	os.Remove(filepath.Join(path, testFile))
	assert.Equal(t, nil, err)

}

func TestBuilderInit(t *testing.T) {
	b := NewBuilder(&Configuration{},
		filepath.Join(os.TempDir(), "test-init"),
		"foo",
		"yaml",
		customProps)

	err := b.Init()
	assert.Equal(t, nil, err)
}

func TestBuilderSave(t *testing.T) {
	b := NewBuilder(nil,
		filepath.Join(os.TempDir(), "test-save"),
		"foo",
		"yaml",
		customProps)

	err := b.Init()
	b.SetConfiguration(&Configuration{})
	assert.Equal(t, nil, err)

	c := &Configuration{
		App: &App{
			Name:    "foo",
			Project: "bar",
		},
		Server: &Server{
			Port: "8080",
		},
	}

	t.Run("should save struct to file", func(t *testing.T) {
		err = b.Save(c)
		assert.Equal(t, nil, err)
	})

	t.Run("should save struct to file", func(t *testing.T) {
		err = b.Save("wrong-format")
		assert.Contains(t, err.Error(), "wrong")
	})
}
