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

// Package locale provides the hiboot starter for injectable locale (i18n) dependency
package locale

import (
	"github.com/kataras/iris/middleware/i18n"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web/context"
	"github.com/hidevopsio/hiboot/pkg/utils/io"
	"os"
	"path/filepath"
	"strings"
)

const (
	// Profile is the profile of locale, it should be as same as the package name
	Profile = "locale"
)

type configuration struct {
	app.Configuration
	Properties         *properties
	applicationContext app.ApplicationContext
}

func newConfiguration(applicationContext app.ApplicationContext) *configuration {
	return &configuration{
		applicationContext: applicationContext,
	}
}

func init() {
	app.Register(newConfiguration)
}

func (c *configuration) Handler() (handler context.Handler) {
	// TODO: localePath should be configurable in application.yml
	// locale:
	//   en-US: ./config/i18n/en-US.ini
	//   cn-ZH: ./config/i18n/cn-ZH.ini
	// TODO: or
	// locale:
	//   path: ./config/i18n/

	localePath, _ := filepath.Abs(c.Properties.LocalePath)
	// parse language files
	languages := make(map[string]string)
	err := filepath.Walk(localePath, func(lngPath string, info os.FileInfo, err error) error {
		if err == nil {
			//*files = append(*files, path)
			lng := strings.Replace(lngPath, localePath, "", 1)
			lng = io.BaseDir(lng)
			lng = strings.Replace(lng, string(filepath.Separator), "", -1)
			if lng != "" && lng != "." && lngPath != localePath+lng {
				//languages[lng] = path
				if languages[lng] == "" {
					languages[lng] = lngPath
				} else {
					languages[lng] = languages[lng] + ", " + lngPath
				}
				//log.Debugf("%v, %v", lng, languages[lng])
			}
		}
		return err
	})
	if err == nil && len(languages) != 0 {
		handler = context.NewHandler(i18n.New(i18n.Config{
			Default:      c.Properties.Default,
			URLParameter: c.Properties.URLParameter,
			Languages:    languages,
		}))

		c.applicationContext.Use(handler)
	}

	return handler
}
