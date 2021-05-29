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

package main

import (
	"embed"

	"github.com/hidevopsio/hiboot/examples/cli/advanced/cmd"
	"github.com/hidevopsio/hiboot/examples/cli/advanced/config"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/cli"
	"github.com/hidevopsio/hiboot/pkg/starter/logging"
)

//go:embed config/foo
var embedFS embed.FS

func main() {
	// create new cli application and run it
	cli.NewApplication(cmd.NewRootCommand).
		SetProperty(app.Config, &embedFS).
		SetProperty(logging.Level, logging.LevelError).
		SetProperty(app.ProfilesInclude, config.Profile, logging.Profile).
		Run()
}
