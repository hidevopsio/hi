package web

import "github.com/hidevopsio/hiboot/pkg/inject/annotation"

// HttpMethodSubscriber
type HttpMethodSubscriber interface {
	Subscribe(atController *annotation.Annotations, atMethod *annotation.Annotations)
}