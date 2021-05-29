package jaeger

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	webctx "github.com/hidevopsio/hiboot/pkg/app/web/context"
	"github.com/hidevopsio/hiboot/pkg/at"
	"net/http"
)

//Span  is the wrap of opentracing.Span
type Span struct {
	at.ContextAware
	opentracing.Span
	context webctx.Context
}

//ChildSpan
type ChildSpan Span


func (s *Span) Inject(ctx context.Context, method string, url string, req *http.Request) opentracing.Span {
	c := opentracing.ContextWithSpan(ctx, s)
	newSpan, _ := opentracing.StartSpanFromContext(c, req.RequestURI)

	ext.SpanKindRPCClient.Set(newSpan)
	ext.HTTPUrl.Set(newSpan, url)
	ext.HTTPMethod.Set(newSpan, method)
	newSpan.Tracer().Inject(
		newSpan.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)

	return newSpan
}


func (s *ChildSpan) Inject(ctx context.Context, method string, url string, req *http.Request) opentracing.Span {
	c := opentracing.ContextWithSpan(ctx, s)
	newSpan, _ := opentracing.StartSpanFromContext(c, req.RequestURI)

	ext.SpanKindRPCClient.Set(newSpan)
	ext.HTTPUrl.Set(newSpan, url)
	ext.HTTPMethod.Set(newSpan, method)
	newSpan.Tracer().Inject(
		newSpan.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)

	return newSpan
}
