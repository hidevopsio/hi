package webutils

import (
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/inject/annotation"
)

func GetHttpMethod(atMethod *annotation.Annotations) (method string, path string) {
	// parse http method
	hma := annotation.Find(atMethod, at.HttpMethod{})
	if hma != nil {
		hm := hma.Field.Value.Interface()
		switch hm.(type) {
		case at.GetMapping:
			httpMethod := hm.(at.GetMapping)
			method, path = httpMethod.AtMethod, httpMethod.AtValue
		case at.PostMapping:
			httpMethod := hm.(at.PostMapping)
			method, path = httpMethod.AtMethod, httpMethod.AtValue
		case at.PutMapping:
			httpMethod := hm.(at.PutMapping)
			method, path = httpMethod.AtMethod, httpMethod.AtValue
		case at.DeleteMapping:
			httpMethod := hm.(at.DeleteMapping)
			method, path = httpMethod.AtMethod, httpMethod.AtValue
		case at.PatchMapping:
			httpMethod := hm.(at.PatchMapping)
			method, path = httpMethod.AtMethod, httpMethod.AtValue
		case at.OptionsMapping:
			httpMethod := hm.(at.OptionsMapping)
			method, path = httpMethod.AtMethod, httpMethod.AtValue
		case at.AnyMapping:
			httpMethod := hm.(at.AnyMapping)
			method, path = httpMethod.AtMethod, httpMethod.AtValue
		case at.TraceMapping:
			httpMethod := hm.(at.TraceMapping)
			method, path = httpMethod.AtMethod, httpMethod.AtValue
		}
	}

	return
}
