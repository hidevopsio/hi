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

// Package jaeger provides the hiboot starter for injectable jaeger dependency
package jaeger

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go/config"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/app/web/context"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/log"
	"io"
)

const (
	// Profile is the profile of jwt, it should be as same as the package name
	Profile = "jaeger"
)

type configuration struct {
	at.AutoConfiguration

	Properties *properties
	Closer     io.Closer

	ServiceName string `value:"${app.name}"`
}

func init() {
	app.Register(newConfiguration, new(properties))
}

func newConfiguration(properties *properties) *configuration {
	return &configuration{Properties: properties}
}

//Tracer returns an instance of Jaeger Tracer that samples 100% of traces and logs all spans to stdout.
func (c *configuration) Tracer() (tracer Tracer) {
	var err error
	if c.Properties.Config.ServiceName == "" {
		c.Properties.Config.ServiceName = c.ServiceName
	}
	tracer, c.Closer, err = c.Properties.Config.NewTracer(config.Logger(&Logger{}))
	log.Debug(err)
	opentracing.SetGlobalTracer(tracer)
	return tracer
}

func (c *configuration) path(ctx context.Context) (path string) {
	currentRoute := ctx.GetCurrentRoute()
	path = currentRoute.Path() + " => " + currentRoute.MainHandlerName() + "()"
	return
}

//Span returns an instance of Jaeger root span.
func (c *configuration) Span(ctx context.Context, tracer Tracer) (span *Span) {
	span = new(Span)
	span.Span = tracer.StartSpan( c.path(ctx) )
	span.context = ctx
	return span
}

//ChildSpan returns an instance of Jaeger child span from parent Span.
//1. Extract the span context from the incoming request using tracer.Extract
//2. Start a new child span representing the work of the server
func (c *configuration) ChildSpan(ctx context.Context, tracer Tracer) (span *ChildSpan) {
	span = new(ChildSpan)

	spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(ctx.Request().Header))
	span.Span = tracer.StartSpan( c.path(ctx), ext.RPCServerOption(spanCtx))
	span.context = ctx
	return span
}
