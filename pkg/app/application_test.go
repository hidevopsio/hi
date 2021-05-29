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

package app_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/log"
	"os"
	"sync"
	"testing"
)

var mux = &sync.Mutex{}

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestApp(t *testing.T) {
	mux.Lock()
	type fakeProperties struct {
		at.ConfigurationProperties `value:"fake"`
		at.AutoWired

		Name string `default:"fake"`
	}
	type fakeConfiguration struct {
		app.Configuration
		Properties fakeProperties `mapstructure:"fake"`
	}

	t.Run("should add configuration", func(t *testing.T) {
		app.Register(new(fakeConfiguration))
	})

	//t.Run("should report duplication error", func(t *testing.T) {
	//	err := app.AutoConfiguration(new(fakeConfiguration))
	//	assert.Equal(t, app.ConfigurationNameIsTakenError, err)
	//})

	//t.Run("should not add invalid configuration", func(t *testing.T) {
	//	type fooConfiguration struct {
	//		Properties fakeProperties `mapstructure:"fake"`
	//	}
	//	err := app.AutoConfiguration(fooConfiguration{})
	//	assert.Equal(t, app.ErrInvalidObjectType, err)
	//})

	type configuration struct {
		app.Configuration
		Properties *fakeProperties `inject:""`
	}
	t.Run("should add configuration with pkg name", func(t *testing.T) {
		app.Register(new(configuration))
	})

	//t.Run("should add named configuration", func(t *testing.T) {
	//	err := app.AutoConfiguration("baz", new(configuration))
	//	assert.Equal(t, nil, err)
	//})

	t.Run("should not add invalid configuration", func(t *testing.T) {
		app.Register(nil)
	})

	t.Run("should add configuration with pkg name", func(t *testing.T) {
		type bazConfiguration struct {
			app.Configuration
			Properties fakeProperties `mapstructure:"fake"`
		}
		app.Register(new(bazConfiguration))
	})

	t.Run("should not add invalid component", func(t *testing.T) {
		app.Register(nil)
	})

	t.Run("should add new component", func(t *testing.T) {
		type fakeService interface{}
		type fakeServiceImpl struct{ fakeService }
		app.Register(new(fakeServiceImpl))
	})

	t.Run("should add new named component", func(t *testing.T) {
		type fakeService interface{}
		type fakeServiceImpl struct{ at.Qualifier `value:"myService"`; fakeService }
		 app.Register(new(fakeServiceImpl))
	})

	t.Run("should add more than one new component at the same time", func(t *testing.T) {
		type fakeService interface{}
		type fakeFooService struct{ fakeService }
		type fakeBarService struct{ fakeService }
		app.Register(new(fakeFooService), new(fakeBarService))
	})
	mux.Unlock()
}

func TestBaseApplication(t *testing.T) {

	mux.Lock()
	os.Args = append(os.Args, "--app.profiles.active=local", "--test.property")

	ba := new(app.BaseApplication)

	err := ba.Initialize()
	assert.Equal(t, nil, err)

	ba.Build()

	sc := ba.SystemConfig()
	assert.NotEqual(t, nil, sc)

	// TODO: check concurrency issue during test
	assert.NotEqual(t, nil, ba)
	err = ba.BuildConfigurations()
	//assert.Equal(t, nil, err)

	t.Run("should find instance by name", func(t *testing.T) {
		ba.GetInstance("foo")
	})

	cf := ba.ConfigurableFactory()
	assert.NotEqual(t, nil, cf)

	ba.SetAddCommandLineProperties(false)
	ba.AfterInitialization()

	ba.SetAddCommandLineProperties(true)
	ba.AfterInitialization()

	ba.RegisterController(nil)

	t.Run("should set PropertyBannerDisabled", func(t *testing.T) {
		ba.SetProperty(app.BannerDisabled, false)
		prop, ok := ba.GetProperty(app.BannerDisabled)
		assert.Equal(t, true, ok)
		assert.Equal(t, false, prop)
	})
	t.Run("should set PropertyBannerDisabled to true", func(t *testing.T) {
		ba.SetProperty(app.BannerDisabled, true)
		prop, ok := ba.GetProperty(app.BannerDisabled)
		assert.Equal(t, true, ok)
		assert.Equal(t, true, prop)
	})

	t.Run("should set profiles", func(t *testing.T) {
		ba.SetProperty(app.BannerDisabled, false).
			SetProperty(app.ProfilesInclude, "foo,bar")
		prop, ok := ba.GetProperty(app.ProfilesInclude)
		assert.Equal(t, true, ok)
		assert.Equal(t, []string{"foo", "bar"}, prop)
	})

	t.Run("should set profiles", func(t *testing.T) {
		ba.SetProperty(app.BannerDisabled, false).
			SetProperty(app.ProfilesInclude, "baz", "buz")
		prop, ok := ba.GetProperty(app.ProfilesInclude)
		assert.Equal(t, true, ok)
		assert.Equal(t, []interface{}{"baz", "buz"}, prop)
	})

	ba.PrintStartupMessages()

	ba.Use()

	ba.Run()

	ba.GetInstance("foo")
	mux.Unlock()
}
