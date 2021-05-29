package at

// Component is the annotation that the dependency is injected in app init.
type Component struct {
	Annotation

	BaseAnnotation
}

// AutoWired is the annotation that auto inject instance to object
type AutoWired struct {
	Annotation

	BaseAnnotation
}
