package at

// Middleware is the annotation that annotate the controller or method use middleware
type Middleware struct {
	Annotation

	BaseAnnotation
}

// MiddlewareHandler is the annotation that annotate the controller or method use middleware
type MiddlewareHandler struct {
	Annotation

	BaseAnnotation
}

// MiddlewarePostHandler is the annotation that annotate the controller or method use middleware
type MiddlewarePostHandler struct {
	Annotation

	BaseAnnotation
}

// UseMiddleware is the annotation that that annotate the controller or method use middleware based on condition
type UseMiddleware struct {
	Annotation

	Conditional
}

// UseJwt
type UseJwt struct {
	Annotation

	UseMiddleware
}