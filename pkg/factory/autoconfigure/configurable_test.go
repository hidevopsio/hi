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

package autoconfigure_test

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/app/web/context"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/factory"
	"github.com/hidevopsio/hiboot/pkg/factory/autoconfigure"
	"github.com/hidevopsio/hiboot/pkg/factory/instantiate"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils/cmap"
	"github.com/hidevopsio/hiboot/pkg/utils/io"
	"github.com/stretchr/testify/assert"
)

type FakeProperties struct {
	at.ConfigurationProperties `value:"fake"`
	at.AutoWired

	Name     string `default:"foo"`
	Nickname string `default:"foobar"`
	Username string `default:"fb"`
	Org      string `default:"hidevopsio"`
	Profile  string `default:"${APP_PROFILES_ACTIVE:dev}"`
}

type Connection struct {
	context context.Context
}

type ContextAwareConfiguration struct {
	at.AutoConfiguration
	at.ContextAware
}

func newContextAwareConfiguration() *ContextAwareConfiguration {
	return &ContextAwareConfiguration{}
}

func (c *ContextAwareConfiguration) Connection(context context.Context) *Connection {
	return &Connection{context: context}
}

type unknownConfiguration struct {
	FakeProperties FakeProperties
}

func newUnknownConfiguration(fakeProperties FakeProperties) *unknownConfiguration {
	return &unknownConfiguration{FakeProperties: fakeProperties}
}

type FooInterface interface {
}

type unsupportedConfiguration struct {
	FooInterface
	FakeProperties *FakeProperties
}

func newUnsupportedConfiguration(fakeProperties *FakeProperties) *unsupportedConfiguration {
	return &unsupportedConfiguration{FakeProperties: fakeProperties}
}

type fooProperties struct {
	at.ConfigurationProperties `value:"foo"`
	at.AutoWired

	Name     string `default:"${fake.name}"`
	Nickname string `default:"foobar"`
	Username string `default:"fb"`
}

type foobarProperties struct {
	at.ConfigurationProperties `value:"foobar"`
	at.AutoWired

	Name     string `default:"${fake.name}"`
	Nickname string `default:"foobar"`
	Username string `default:"fb"`
}


type barProperties struct {
	at.ConfigurationProperties `value:"bar"`
	at.AutoWired

	Name     string `default:"${fake.name}"`
	Nickname string `default:"foobar"`
	Username string `default:"fb"`
}

type configuration struct {
	at.AutoConfiguration
}

func newConfiguration() *configuration {
	return &configuration{}
}

type emptyConfiguration struct {
}

func newEmptyConfiguration() *emptyConfiguration {
	return &emptyConfiguration{}
}

type fooConfiguration struct {
	at.AutoConfiguration
	FooProperties *fooProperties
}

func newFooConfiguration(fooProperties *fooProperties) *fooConfiguration {
	return &fooConfiguration{FooProperties: fooProperties}
}

type BarConfiguration struct {
	app.Configuration
	barProperties *barProperties
}

type FooBarConfiguration struct {
	app.Configuration
	FakeProperties *foobarProperties
	foobar         *FooBar
}

func NewFooBarConfiguration(fakeProperties *foobarProperties, foobar *FooBar) *FooBarConfiguration {
	return &FooBarConfiguration{FakeProperties: fakeProperties, foobar: foobar}
}

type mercury struct {
	jupiter *jupiter
}

type mars struct {
	mercury *mercury
}

type jupiter struct {
}

type marsConfiguration struct {
	app.Configuration
}

func newMarsConfiguration() *marsConfiguration {
	return &marsConfiguration{}
}

func (c *marsConfiguration) Mars(m *mercury) *mars {
	return &mars{mercury: m}
}

type mercuryConfiguration struct {
	app.Configuration
}

func newMercuryConfiguration() *mercuryConfiguration {
	return &mercuryConfiguration{}
}

func (c *mercuryConfiguration) Mercury(j *jupiter) *mercury {
	return &mercury{jupiter: j}
}

type jupiterConfiguration struct {
	app.Configuration
}

func newJupiterConfiguration() *jupiterConfiguration {
	return &jupiterConfiguration{}
}

func (c *jupiterConfiguration) Jupiter() *jupiter {
	return new(jupiter)
}

func newFooBarConfiguration(foobar *FooBar) *FooBarConfiguration {
	return &FooBarConfiguration{
		foobar: foobar,
	}
}

type HelloWorld string
type Hello string

func (c *fooConfiguration) HelloWorld(foo Hello) HelloWorld {
	return HelloWorld(foo + " world")
}

func (c *fooConfiguration) Hello() Hello {
	return Hello("Hello")
}

type Foo struct {
	Name string
	Bar  *Bar
}

type Bar struct {
	Name string
}

type FooBar struct {
	Name string
}

type barConfiguration struct {
	at.AutoConfiguration
	
	barProperties *barProperties 
}

func newBarConfiguration(barProperties *barProperties) *barConfiguration {
	return &barConfiguration{barProperties: barProperties}
}

func (c *barConfiguration) BarBar() *Bar {
	return new(Bar)
}

func init() {
	log.SetLevel(log.DebugLevel)
	io.EnsureWorkDir(1, "config/application.yml")
}

type fakeConfiguration struct {
	app.Configuration
	FakeProperties *FakeProperties 
}

func newFakeConfiguration(fakeProperties *FakeProperties) *fakeConfiguration {
	return &fakeConfiguration{FakeProperties: fakeProperties}
}

type foobarConfiguration struct {
	app.Configuration
	FoobarProperties *foobarProperties
}

func newFoobarConfiguration(properties *foobarProperties) *foobarConfiguration {
	return &foobarConfiguration{FoobarProperties: properties}
}

func (c *foobarConfiguration) Foo(bar *Bar) *Foo {
	f := new(Foo)
	f.Name = c.FoobarProperties.Name
	f.Bar = bar
	return f
}

func (c *foobarConfiguration) FooBar(foo *Foo) *FooBar {
	f := new(FooBar)
	f.Name = foo.Name

	return f
}

func (c *foobarConfiguration) Bar() *Bar {
	b := new(Bar)
	b.Name = "bar"

	return b
}

type EarthConfiguration struct {
	app.Configuration
	instantiateFactory factory.InstantiateFactory
}

func newEarthConfiguration(instantiateFactory factory.InstantiateFactory) *EarthConfiguration {
	c := &EarthConfiguration{
		instantiateFactory: instantiateFactory,
	}
	c.RuntimeDeps.Set(c.RuntimeTree, []string{"autoconfigure_test.leaf", "autoconfigure_test.branch"})
	return c
}

type Land struct {
	Mountain *Mountain
}

type RuntimeTree struct {
	leaf   *Leaf
	branch *Branch
}

type Tree interface {
	Family() string
}

type oakTree struct {
	Tree
}

func (t *oakTree) Family() string {
	return "Beech"
}

type mapleTree struct {
}

func (t *mapleTree) Family() string {
	return "Soapberry"
}

type Branch struct {
}

type Leaf struct {
}

type Mountain struct {
	oakTree Tree
}

func (c *EarthConfiguration) Land(mountain *Mountain) *Land {
	return &Land{Mountain: mountain}
}

func (c *EarthConfiguration) Mountain(tree Tree) *Mountain {
	return &Mountain{oakTree: tree}
}

func (c *EarthConfiguration) Branch() *Branch {
	return &Branch{}
}

func (c *EarthConfiguration) Leaf() *Leaf {
	return &Leaf{}
}

func (c *EarthConfiguration) RuntimeTree() *RuntimeTree {
	leaf := c.instantiateFactory.GetInstance("autoconfigure_test.leaf").(*Leaf)
	branch := c.instantiateFactory.GetInstance("autoconfigure_test.branch").(*Branch)
	return &RuntimeTree{leaf: leaf, branch: branch}
}

func (c *EarthConfiguration) OakTree() Tree {
	return &oakTree{}
}

func (c *EarthConfiguration) MapleTree() Tree {
	return &mapleTree{}
}

type helloService struct{ foo *Foo }

func newHelloService(foo *Foo) *helloService {
	return &helloService{foo: foo}
}

func setFactory(t *testing.T, configDir string, customProperties cmap.ConcurrentMap) factory.ConfigurableFactory {
	io.ChangeWorkDir(os.TempDir())

	configPath := filepath.Join(os.TempDir(), "config")
	defaultConfigFile := "application.yml"
	os.Remove(filepath.Join(configPath, defaultConfigFile))
	fakeContent := `
app:
  project: hidevopsio
  name: hiboot-test
  version: 0.0.1
  profiles:
    include:
    - fake
    - foo
    - autoconfigure_test
    - bar
    - mars
    - jupiter
    - mercury
    - foobar
    - earth
`
	n, err := io.WriterFile(configPath, defaultConfigFile, []byte(fakeContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, n, len(fakeContent))

	fooConfigFile := "application-foo.yml"
	os.Remove(filepath.Join(configPath, fooConfigFile))
	fooContent :=
		"foo:\n" +
			"  name: foo\n" +
			"  nickname: ${app.name} ${foo.name}\n" +
			"  username: ${unknown.name:bar}\n"
	_, err = io.WriterFile(configPath, fooConfigFile, []byte(fooContent))
	assert.Equal(t, nil, err)

	configContainers := cmap.New()
	customProperties.Set("logging.level", "debug")
	f := autoconfigure.NewConfigurableFactory(
		instantiate.NewInstantiateFactory(cmap.New(), []*factory.MetaData{}, customProperties),
		configContainers,
	)

	return f
}

func TestConfigurableFactory(t *testing.T) {
	defaultProperties := cmap.New()
	f := setFactory(t, "mercury", defaultProperties)

	var err error

	// backup profile
	profile := "default" //os.Getenv(autoconfigure.EnvAppProfilesActive)
	t.Run("should build app config", func(t *testing.T) {
		os.Setenv(autoconfigure.EnvAppProfilesActive, "")
		sc, err := f.BuildProperties()
		assert.Equal(t, nil, err)
		assert.Equal(t, "default", sc.App.Profiles.Active)
		// restore profile
		os.Setenv(autoconfigure.EnvAppProfilesActive, profile)
	})

	err = os.Setenv(autoconfigure.EnvAppProfilesActive, profile)
	assert.Equal(t, nil, err)

	t.Run("should build app config", func(t *testing.T) {
		sc, err := f.BuildProperties()
		assert.Equal(t, nil, err)
		assert.Equal(t, profile, sc.App.Profiles.Active)
		// restore profile
		os.Setenv(autoconfigure.EnvAppProfilesActive, profile)
	})

	t.Run("should build app config", func(t *testing.T) {
		_, err = f.BuildProperties()
		assert.Equal(t, nil, err)
	})
	
	f.Build([]*factory.MetaData{
		factory.NewMetaData(newEmptyConfiguration),
		factory.NewMetaData(newFakeConfiguration),
		factory.NewMetaData(newFooConfiguration),
		factory.NewMetaData(newConfiguration),
		factory.NewMetaData(newBarConfiguration),
		factory.NewMetaData(newMarsConfiguration),
		factory.NewMetaData(newJupiterConfiguration),
		factory.NewMetaData(newMercuryConfiguration),
		factory.NewMetaData(newUnknownConfiguration),
		factory.NewMetaData(newUnsupportedConfiguration),
		factory.NewMetaData(newFoobarConfiguration),
		factory.NewMetaData(newEarthConfiguration),
		factory.NewMetaData(newContextAwareConfiguration),
	})

	f.AppendComponent(newHelloService)
	ctx := web.NewContext(nil)
	f.AppendComponent("context.context", ctx)

	err = f.BuildComponents()

	t.Run("should get SystemConfiguration", func(t *testing.T) {
		sysCfg := f.SystemConfiguration()
		assert.NotEqual(t, nil, sysCfg)
	})

	t.Run("should get fake configuration", func(t *testing.T) {
		fake := f.Configuration("fake")
		assert.NotEqual(t, nil, fake)
	})

	t.Run("should not get non-existence configuration", func(t *testing.T) {
		nonExist := f.Configuration("non-existence")
		assert.Equal(t, nil, nonExist)
	})

	t.Run("should add instance to factory at runtime", func(t *testing.T) {
		fakeInstance := &struct{ Name string }{Name: "fake"}
		f.SetInstance("autoconfigure_test.fakeInstance", fakeInstance)
		gotFakeInstance := f.GetInstance("autoconfigure_test.fakeInstance")
		assert.Equal(t, fakeInstance, gotFakeInstance)
	})

	t.Run("should get foo configuration", func(t *testing.T) {
		helloWorld := f.GetInstance("autoconfigure_test.helloWorld")
		assert.NotEqual(t, nil, helloWorld)
		assert.Equal(t, HelloWorld("Hello world"), helloWorld)
	})

	t.Run("should get runtime created instances", func(t *testing.T) {
		runtimeTree := f.GetInstance(RuntimeTree{})
		leaf := f.GetInstance(Leaf{})
		branch := f.GetInstance(Branch{})
		assert.NotEqual(t, nil, runtimeTree)
		rt := runtimeTree.(*RuntimeTree)
		assert.Equal(t, branch, rt.branch)
		assert.Equal(t, leaf, rt.leaf)
	})

	t.Run("should report error on nil configuration", func(t *testing.T) {

	})
}

func TestReplacer(t *testing.T) {
	customProperties := cmap.New()
	customProperties.Set("app.profiles.filter", true)
	customProperties.Set("app.profiles.include", []string{"foo", "fake"})
	f := setFactory(t, "jupiter", customProperties)
	var err error

	_, err = f.BuildProperties()
	t.Run("should build app config", func(t *testing.T) {
		assert.Equal(t, nil, err)
	})


	type outConfiguration struct {
		at.AutoConfiguration
		Properties *fooProperties `inject:""`
	}

	f.Build([]*factory.MetaData{
		factory.NewMetaData(new(outConfiguration)),
		factory.NewMetaData(newFakeConfiguration),
		factory.NewMetaData(newFooConfiguration),
	})

	fp := f.GetInstance(fooProperties{})
	assert.NotEqual(t, nil, fp)
	fooProp := fp.(*fooProperties)
	
	t.Run("should get foo configuration", func(t *testing.T) {
		assert.Equal(t, "hiboot-test foo", fooProp.Nickname)
		assert.Equal(t, "bar", fooProp.Username)
		assert.Equal(t, "foo", fooProp.Name)
	})
}

var doneSch = make(chan bool)

type myService struct {
	at.EnableScheduling

	count int
}

func newMyService() *myService {
	return &myService{count: 10}
}

//_ struct{at.Scheduler `limit:"10"`}
func (s *myService) Run(_ struct{at.Scheduled `every:"200" unit:"milliseconds"`}) (done bool) {
	log.Info("Running Scheduler Task")

	if s.count <= 0 {
		done = true
		doneSch <- true
	}
	s.count--

	return
}

type controller struct {
	at.RestController
}

func (c *controller) Get() string {
	return "Hello scheduler"
}


func TestScheduler(t *testing.T) {
	app.Register(newMyService)
	testApp := web.NewTestApp(t, new(controller)).Run(t)


	t.Run("scheduler", func(t *testing.T) {
		testApp.Get("/").Expect().Status(http.StatusOK)
	})

	log.Infof("scheduler is done: %v", <- doneSch)
}