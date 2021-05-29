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

import "github.com/hidevopsio/hiboot/pkg/at"

var (
	// RequestTypeBody means it is RequestBody
	RequestTypeBody = "RequestBody"
	// RequestTypeParams means it is RequestParams
	RequestTypeParams = "RequestParams"
	// RequestTypeForm means it is RequestForm
	RequestTypeForm = "RequestForm"
	// Context means it is Context
	Context = "Context"
)

// RequestBody the annotation RequestBody
type RequestBody struct{
	at.RequestBody
}

// RequestForm the annotation RequestForm
type RequestForm struct{
	at.RequestForm
}

// RequestParams the annotation RequestParams
type RequestParams struct{
	at.RequestParams
}

// ListOptions specifies the optional parameters to various List methods that
// support pagination.
type ListOptions struct {
	// For paginated result sets, page of results to retrieve.
	Page int `url:"page,omitempty" json:"page,omitempty" validate:"min=1"`

	// For paginated result sets, the number of results to include per page.
	PerPage int `url:"per_page,omitempty" json:"per_page,omitempty" validate:"min=1"`
}