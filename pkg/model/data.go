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

package model

import "time"

// BaseData specifies the base data fields for database models
type BaseData struct {
	// IsDeleted for soft delete
	IsDeleted bool `json:"is_deleted,omitempty"`

	// CreatedAt data created time
	CreatedAt time.Time `json:"created_at,omitempty"`

	// UpdatedAt data updated time
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
