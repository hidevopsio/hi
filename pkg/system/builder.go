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

// TODO: app config should be generic kit

package system

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils/io"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"github.com/hidevopsio/hiboot/pkg/utils/replacer"
	"github.com/hidevopsio/hiboot/pkg/utils/str"
	"github.com/hidevopsio/mapstructure"
	"github.com/hidevopsio/viper"
	"gopkg.in/yaml.v2"
)

// Builder is the config file (yaml, json) builder
type Builder interface {
	Init() error
	Build(profiles ...string) (p interface{}, err error)
	BuildWithProfile(profile string) (interface{}, error)
	Load(properties interface{}, opts ...func (*mapstructure.DecoderConfig)) (err error)
	Save(p interface{}) (err error)
	Replace(source string) (retVal interface{})
	GetProperty(name string) (retVal interface{})
	SetProperty(name string, val interface{}) Builder
	SetDefaultProperty(name string, val interface{}) Builder
	SetConfiguration(in interface{})
}

// Deprecated, use propertyBuilder instead
type builder struct {
	*viper.Viper
	path             string
	name             string
	fileType         string
	configuration    interface{}
	customProperties map[string]interface{}
	profiles         []string
}

func (b *builder) Load(properties interface{}, opts ...func (*mapstructure.DecoderConfig)) (err error) {
	return
}

// Deprecated
// use NewPropertyBuilder instead
// NewBuilder is the constructor of system.Builder
func NewBuilder(configuration interface{}, path, name, fileType string, customProperties map[string]interface{}) Builder {
	return &builder{
		Viper:            viper.New(),
		path:             path,
		name:             name,
		fileType:         fileType,
		configuration:    configuration,
		customProperties: customProperties,
	}
}

// New create new viper instance
func (b *builder) config(name string) {
	viper.Reset()
	b.AutomaticEnv()
	viperReplacer := strings.NewReplacer(".", "_")
	b.SetEnvKeyReplacer(viperReplacer)

	b.AddConfigPath(b.path)
	b.SetConfigName(name)
	b.SetConfigType(b.fileType)
}

func (b *builder) SetConfiguration(in interface{}) {
	b.configuration = in
}

// Init create file if it's not exist
func (b *builder) Init() error {
	return io.CreateFile(b.path, b.name+"."+b.fileType)
}

func (b *builder) isFileNotExist(path string) bool {
	for _, ext := range viper.SupportedExts {
		if !io.IsPathNotExist(path + ext) {
			return false
		}
	}
	return true
}

// Build config file
func (b *builder) Build(profiles ...string) (conf interface{}, err error) {
	// parse profiles
	if b.configuration != nil {
		for _, field := range reflector.DeepFields(reflect.TypeOf(b.configuration)) {
			p, ok := field.Tag.Lookup("mapstructure")
			if ok && !str.InSlice(p, profiles) {
				profiles = append([]string{p}, profiles...)
			}
		}
	}
	// save profiles
	b.profiles = append(b.profiles, profiles...)

	if str.InSlice("default", profiles) {
		b.read(b.name)
		b.load(b.name, "")
	}

	for _, profile := range profiles {
		name := b.name + "-" + profile
		// allow the empty of the profile
		configFile := filepath.Join(b.path, name)
		if profile != "" && !b.isFileNotExist(configFile+".") {
			b.read(name)
		}
		b.load(name, profile)
	}
	return b.configuration, err
}

// BuildWithProfile build config file
func (b *builder) BuildWithProfile(profile string) (interface{}, error) {
	name := b.name
	if profile != "" {
		name = name + "-" + profile
	}
	// allow the empty of the profile
	if profile == "" || b.isFileNotExist(filepath.Join(b.path, name)+".") {
		return reflector.NewReflectType(b.configuration), nil
	}
	b.read(name)
	return b.load(name, profile)
}

// Read single file
func (b *builder) read(fullName string) {

	// config
	b.config(fullName)

	// read config
	b.MergeInConfig()
}

// Read single file
func (b *builder) load(fullName, profile string) (interface{}, error) {

	// create new instance
	conf := b.configuration
	val := reflect.ValueOf(conf)
	//log.Debugf("value of configuration: %v, kind: %v", val, val.Kind())
	if val.Kind() == reflect.Struct {
		conf = reflector.NewReflectType(conf)
	}

	// set custom properties
	for key, value := range b.customProperties {
		keyPath := strings.Split(key, ".")
		if str.InSlice(keyPath[0], b.profiles) {
			b.SetProperty(key, value)
		}
	}

	// iterate all and replace reference values or env
	allKeys := b.AllKeys()
	for _, key := range allKeys {
		val := b.GetString(key)
		if strings.Contains(val, "${") {
			newVal := b.Replace(val)
			b.SetConfig(key, newVal)
			log.Debugf(">>> replaced key: %v, value: %v, newVal: %v", key, val, newVal)
		}
	}

	err := b.Unmarshal(conf)
	return conf, err
}

// Save configurations to file
func (b *builder) Save(p interface{}) error {

	b.config(b.name)

	y, err := yaml.Marshal(p)
	if err == nil {
		err = b.ReadConfig(bytes.NewBuffer(y))
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return err
		}
	}

	return b.WriteConfig()
}

// Replace replace reference and
func (b *builder) Replace(source string) (retVal interface{}) {
	result := source
	matches := replacer.GetMatches(source)
	if len(matches) != 0 {
		for _, m := range matches {
			varFullName, varName := m[0], m[1]
			// check if it contains default value
			var defaultValue string
			n := strings.Index(varName, ":")
			if n > 0 {
				defaultValue = varName[n+1:]
				varName = varName[:n]
				//log.Debugf("name: %v, default value: %v", varName, defaultValue)
			}
			prop := b.Get(varName)

			var newVal string
			if prop != nil {
				switch prop.(type) {
				case string:
					newVal = prop.(string)
					result = strings.Replace(result, varFullName, newVal, -1)
				default:
					retVal = prop
					return
				}
			}

			envValue := os.Getenv(varName)
			// check if  varName == strings.ToUpper(varName), the assume that varName is environment variable
			if envValue != "" || (varName == strings.ToUpper(varName) && defaultValue == "") {
				result = strings.Replace(result, varFullName, envValue, -1)
			}

			if envValue == "" && newVal == "" && defaultValue != "" {
				result = strings.Replace(result, varFullName, defaultValue, -1)
			}
			log.Debugf("replaced %v to %v", varName, result)
		}
	}
	retVal = result
	return
}

func (b *builder) GetProperty(name string) (retVal interface{}) {
	retVal = b.Get(name)
	return
}

func (b *builder) SetProperty(name string, val interface{}) Builder {
	b.Set(name, val)
	return b
}

func (b *builder) SetDefaultProperty(name string, val interface{}) Builder {
	// TODO: bug ...
	b.SetDefault(name, val)
	return b
}
