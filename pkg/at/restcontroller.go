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

// Package at provides annotations for web RestController
package at

// RestController is the annotation that declare current controller is the RESTful Controller
type RestController struct {
	Annotation

	BaseAnnotation
}

// JwtRestController is the annotation that declare current controller is the RESTful Controller with JWT support
type JwtRestController struct {
	Annotation

	RestController
	UseJwt
}

type HttpMethodSubscriber struct{
	Annotation

	BaseAnnotation
}

// ContextPath is the annotation that set the context path of a controller
type ContextPath struct {
	Annotation

	BaseAnnotation
}

// RequestMapping is the annotation that set the RequestMapping of a controller
type RequestMapping struct {
	Annotation

	HttpMethod
}

// Method is the annotation that set the RequestMethod of a controller
type Method struct {
	Annotation

	BaseAnnotation
}

// HttpMethod is the annotation that the http method of a controller
type HttpMethod struct {
	Annotation

	BaseAnnotation

	AtNoContextPath bool `json:"-" at:"no_context_path"`
}

// BeforeMethod is the annotation that set the called before the http method of a controller
type BeforeMethod struct {
	Annotation

	HttpMethod
}

// AfterMethod is the annotation that set the called after the http method of a controller
type AfterMethod struct {
	Annotation

	HttpMethod
}

// GetMapping is the annotation that set the GetMapping of a controller
type GetMapping struct {
	Annotation

	RequestMapping
	AtMethod string `method:"GET" at:"method" json:"-"`
}

// PostMapping is the annotation that set the PostMapping of a controller
type PostMapping struct {
	Annotation

	RequestMapping
	AtMethod string  `method:"POST" at:"method" json:"-"`
}

// PutMapping is the annotation that set the PutMapping of a controller
type PutMapping struct {
	Annotation

	RequestMapping
	AtMethod string  `method:"PUT" at:"method" json:"-"`
}

// PatchMapping is the annotation that set the PatchMapping of a controller
type PatchMapping struct {
	Annotation

	RequestMapping
	AtMethod string  `method:"PATCH" at:"method" json:"-"`
}

// DeleteMapping is the annotation that set the DeleteMapping of a controller
type DeleteMapping struct {
	Annotation

	RequestMapping
	AtMethod string  `method:"DELETE" at:"method" json:"-"`
}

// AnyMapping is the annotation that set the AnyMapping of a controller
type AnyMapping struct {
	Annotation

	RequestMapping
	AtMethod string  `value:"ANY" at:"method" json:"-"`
}

// OptionsMapping is the annotation that set the OptionsMapping of a controller
type OptionsMapping struct {
	Annotation

	RequestMapping
	AtMethod string  `method:"OPTIONS" at:"method" json:"-"`
}

// TraceMapping is the annotation that set the TraceMapping of a controller
type TraceMapping struct {
	Annotation

	RequestMapping
	AtMethod string  `method:"TRACE" at:"method" json:"-"`
}

// StaticResource is the annotation that set the StaticResource of a controller
// value: static resource dir
type FileServer struct {
	Annotation

	RequestMapping
}