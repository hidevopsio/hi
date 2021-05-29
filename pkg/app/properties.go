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

package app

import "github.com/hidevopsio/hiboot/pkg/system"

const (
	// Config is the embed config
	Config = system.Config

	// ConfigDir is the embed config dir
	ConfigDir = system.ConfigDir

	// BannerDisabled is the property that allow use to enable / disable banner display on terminal
	BannerDisabled = "app.banner.disabled"

	// ProfilesInclude is the property that allow user include profiles at runtime
	ProfilesInclude = "app.profiles.include"

	// Version is the property key of app.version
	Version = "app.version"
)