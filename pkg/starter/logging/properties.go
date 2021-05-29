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

package logging

import "github.com/hidevopsio/hiboot/pkg/at"

// Properties is the logging properties
type properties struct {
	at.ConfigurationProperties `value:"logging"`
	at.AutoWired

	Level       string   `json:"level" default:"info"`
	Status      bool     `json:"status" default:"true"`
	IP          bool     `json:"ip" default:"true"`
	Method      bool     `json:"method" default:"true"`
	Path        bool     `json:"path" default:"true"`
	Query       bool     `json:"query" default:"false"`
	Columns     bool     `json:"columns" default:"false"`
	ContextKeys []string `json:"context_keys" default:"logger_message"`
	HeaderKeys  []string `json:"header_keys" default:"User-Agent"`
}
