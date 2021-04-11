package at

// Middleware is the annotation that annotate the controller or method use middleware
type Middleware struct {
	Annotation

	BaseAnnotation
}

// BeforeHandler is the annotation that annotate the controller or method use middleware
type BeforeHandler struct {
	Annotation

	BaseAnnotation
}


// EmbedHandler is the annotation that annotate the controller or method use middleware
// This handler is embedded in to the method
type EmbedHandler struct {
	Annotation

	BaseAnnotation
}

// AfterHandler is the annotation that annotate the controller or method use middleware
type AfterHandler struct {
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