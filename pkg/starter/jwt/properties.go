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

package jwt

import "github.com/hidevopsio/hiboot/pkg/at"

// Properties the jwt properties
type Properties struct {
	at.ConfigurationProperties `value:"jwt"`
	at.AutoWired

	PrivateKeyPath string `json:"private_key_path" default:"config/ssl/app.rsa"`
	PublicKeyPath  string `json:"public_key_path" default:"config/ssl/app.rsa.pub"`
}
