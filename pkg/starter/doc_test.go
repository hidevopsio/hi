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

package starter_test

import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"os"
)

//This example shows the guide to make customized auto configuration
//for more details, see https://github.com/hidevopsio/hiboot-data/blob/master/starter/bolt/autoconfigure.go
func Example() {
}

// properties
type properties struct {
	Database string      `json:"database" default:"hiboot.db"`
	Mode     os.FileMode `json:"mode" default:"0600"`
	Timeout  int64       `json:"timeout" default:"2"`
}

// declare boltConfiguration
type boltConfiguration struct {
	app.Configuration
	// the properties member name must be Bolt if the mapstructure is bolt,
	// so that the reference can be parsed
	BoltProperties properties `mapstructure:"bolt"`
}

// BoltRepository
type BoltRepository struct {
}

func init() {
	// register newBoltConfiguration as AutoConfiguration
	app.Register(newBoltConfiguration)
}

// boltConfiguration constructor
func newBoltConfiguration() *boltConfiguration {
	return &boltConfiguration{}
}

func (c *boltConfiguration) BoltRepository() *BoltRepository {

	repo := &BoltRepository{}

	// config repo according to c.BoltProperties

	return repo
}
