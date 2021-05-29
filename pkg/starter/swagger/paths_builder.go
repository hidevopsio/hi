package swagger

import (
	"fmt"
	"github.com/go-openapi/spec"
	"github.com/hidevopsio/hiboot/pkg/app/web/webutils"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/inject/annotation"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"github.com/hidevopsio/hiboot/pkg/utils/str"
	"github.com/hidevopsio/hiboot/pkg/utils/structtag"
	"path"
	"reflect"
	"strings"
	"time"
)

const refPrefix = "#/definitions/"

type apiPathsBuilder struct {
	apiInfoBuilder *apiInfoBuilder
	primitiveTypes map[string]string
}

func newApiPathsBuilder(builder *apiInfoBuilder) *apiPathsBuilder {
	if builder.SystemServer != nil {
		if builder.SystemServer.Host != "" {
			if builder.SystemServer.Port != "" {
				builder.SwaggerProps.Host = builder.SystemServer.Host + ":" + builder.SystemServer.Port
			} else {
				builder.SwaggerProps.Host = builder.SystemServer.Host
			}
		}
		if builder.SystemServer.ContextPath != "" {
			builder.SwaggerProps.BasePath = builder.SystemServer.ContextPath
		}
		if len(builder.SystemServer.Schemes) > 0 {
			builder.SwaggerProps.Schemes = builder.SystemServer.Schemes
		}
	}
	if builder.AppVersion != "" {
		builder.Info.Version = builder.AppVersion
	}
	// TODO: save visit for later use
	visit := fmt.Sprintf("%s://%s/swagger-ui", builder.SwaggerProps.Schemes[0], path.Join(builder.SwaggerProps.Host, builder.SwaggerProps.BasePath))
	log.Infof("visit %v to open api doc", visit)

	return &apiPathsBuilder{
		apiInfoBuilder: builder,
		primitiveTypes: map[string]string{
			// array, boolean, integer, number, object, string
			"string": "string",
			"int": "integer",
			"int8": "integer",
			"int16": "integer",
			"int32": "integer",
			"int64": "integer",
			"uint": "integer",
			"uint8": "integer",
			"uint16": "integer",
			"uint32": "integer",
			"uint64": "integer",
			"float32": "number",
			"float64": "number",
			"struct": "object",
			"slice": "array",
			"bool": "boolean",
			"Time": "string",
		},
	}
}

func deepFields(reflectType reflect.Type) []reflect.StructField {
	var fields []reflect.StructField

	if reflectType = reflector.IndirectType(reflectType); reflectType.Kind() == reflect.Struct {
		for i := 0; i < reflectType.NumField(); i++ {
			v := reflectType.Field(i)

			if annotation.IsAnnotation(v.Type) {
				continue
			}

			if v.Anonymous {
				vk :=  reflector.IndirectType(v.Type).Kind()
				if vk == reflect.Struct || vk == reflect.Interface {
					fields = append(fields, deepFields(v.Type)...)
				}
			} else {
				fields = append(fields, v)
			}
		}
	}

	return fields
}

func (b *apiPathsBuilder) buildSchemaArray(definition *spec.Schema, typ reflect.Type, recursive bool)  {
	definition.Type = spec.StringOrArray{"array"}
	// array items
	arrSchema := spec.Schema{}
	arrType := reflector.IndirectType(typ.Elem())
	b.buildSchemaObject(&arrSchema, arrType, recursive)
	definition.Items = &spec.SchemaOrArray{Schema: &arrSchema}
}

func (b *apiPathsBuilder) buildSchemaProperty(definition *spec.Schema, typ reflect.Type, recursive bool)  {
	kind := typ.Kind()

	if kind == reflect.Slice {
		b.buildSchemaArray(definition, typ, recursive)
	} else if kind == reflect.Struct {
		definition.Properties = make(map[string]spec.Schema)
		definition.Type = spec.StringOrArray{"object"}

		for _, f := range deepFields(typ) {
			var jsonName string
			var fieldName string
			var descName string
			var descTag *structtag.Tag

			tags, _ := structtag.Parse(string(f.Tag))

			jsonTag, err := tags.Get("json")
			if err == nil {
				jsonName = jsonTag.Name
			}

			// first, check if field tag schema is presented
			schemaTag, err := tags.Get("schema")
			if err == nil {
				descTag = schemaTag
			} else {
				// assign field tag json to desc
				descTag = jsonTag
			}

			// assign schema properties
			typName := f.Type.Name()
			ps := spec.Schema{}
			fieldKind := f.Type.Kind()
			iTyp := reflector.IndirectType(f.Type)
			switch fieldKind {
			case reflect.Slice:
				if recursive && typ == iTyp {
					ps.Type = spec.StringOrArray{"array"}
					arrSchema := spec.Schema{}
					// ref to itself for recursive object
					arrSchema.Ref = spec.MustCreateRef(refPrefix + iTyp.Name())
					ps.Items = &spec.SchemaOrArray{Schema: &arrSchema}
				} else {
					b.buildSchemaArray(&ps, f.Type, typ == iTyp)
				}

			case reflect.Struct:
				b.buildSchemaObject(&ps, f.Type, typ == iTyp)

			case reflect.Ptr:
				if recursive && typ == iTyp {
					// ref to itself for recursive object
					ps.Ref = spec.MustCreateRef(refPrefix + iTyp.Name())
				} else {
					if iTyp.Kind() == reflect.Struct {
						b.buildSchemaObject(&ps, iTyp, typ == iTyp)
					}
				}

			default:
				// convert primitive types
				swgTypName := b.primitiveTypes[typName]
				ps.Type = spec.StringOrArray{swgTypName}
			}

			// assign schema
			if jsonName != "" {
				fieldName = jsonName
			} else {
				fieldName = str.ToLowerCamel(f.Name)
			}
			if descTag != nil && ps.Ref.Ref.String() == "" {

				if schemaTag == nil {
					descName = str.ToKebab(f.Name)
					descName = strings.Replace(descName, "-", " ", -1)
					descName = str.UpperFirst(descName)
				} else {
					descName = schemaTag.Name
				}

				ps.Title = f.Name
				ps.Description = descName
				ps.Format = typName
				// example
				example := b.parseExample(f.Tag)
				if example != "" {
					ps.Example = example
				}
			}

			if descTag != nil {
				definition.Properties[fieldName] = ps
			}
		}
	}
}

func (b *apiPathsBuilder) parseExample(tagVal reflect.StructTag) string {
	example := tagVal.Get("example")
	if example == "" {
		example = tagVal.Get("default")
	}
	return example
}

func (b *apiPathsBuilder) buildSchemaObject(ps *spec.Schema, typ reflect.Type, recursive bool) (ok bool) {
	if typ == reflect.TypeOf(time.Time{}) {
		swgTypName := b.primitiveTypes[typ.Name()]
		ps.Type = spec.StringOrArray{swgTypName}
	} else {
		// try to find the definition ref first, if it does not exist, then build the schema property, otherwise just assign ref to schema
		refName := typ.Name()
		_, ok = b.apiInfoBuilder.Definitions[refName]
		if !ok {
			newSchema := spec.Schema{}
			newSchema.Type = spec.StringOrArray{"object"}
			b.buildSchemaProperty(&newSchema, typ, recursive)
			b.apiInfoBuilder.Definitions[refName] = newSchema
		}
		ps.Ref = spec.MustCreateRef(refPrefix + refName)
	}
	return
}

func (b *apiPathsBuilder) buildSchema(ann *annotation.Annotation, field *reflect.StructField) (schema *spec.Schema) {
	if field == nil {
		field = &ann.Field.StructField
	}

	atSchema := annotation.GetAnnotation(ann.Parent.Interface, at.Schema{})

	s := atSchema.Field.Value.Interface().(at.Schema)
	schemaType := s.AtType
	primitiveTypes := b.primitiveTypes[schemaType]

	schema = &spec.Schema{}
	if primitiveTypes == "" {
		err := annotation.Inject(atSchema)
		if err == nil {
			// parse body schema and assign to definitions
			schema.Ref = spec.MustCreateRef(refPrefix + field.Name)

			if b.apiInfoBuilder.Definitions == nil {
				def := make(spec.Definitions)
				b.apiInfoBuilder.Definitions = def
			}

			definition, ok := b.apiInfoBuilder.Definitions[field.Name]
			if !ok {
				definition = spec.Schema{}
				b.buildSchemaProperty(&definition, field.Type, false)
				b.apiInfoBuilder.Definitions[field.Name] = definition
			}
		}
	} else {
		schema.Type = spec.StringOrArray{s.AtType}
		schema.Description = s.AtDescription
	}

	return
}

func (b *apiPathsBuilder) buildParameter(operation *spec.Operation, annotations *annotation.Annotations, a *annotation.Annotation) {
	ao := a.Field.Value.Interface()
	atParameter := ao.(at.Parameter)
	// copy values
	parameter := spec.Parameter{}
	parameter.Name = atParameter.AtName
	parameter.Type = atParameter.AtType
	parameter.In = atParameter.AtIn
	parameter.Description = atParameter.AtDescription

	if atParameter.AtIn == "body" || atParameter.AtIn == "array" {

		atSchema := annotation.Find(annotations, at.Schema{})

		if atSchema != nil {

			field := b.findArrayField(atSchema)

			parameter.Schema = b.buildSchema(atSchema, field)
		}
	}

	operation.Parameters = append(operation.Parameters, parameter)
	return
}

func (b *apiPathsBuilder) findArrayField(schema *annotation.Annotation) (field *reflect.StructField) {
	parentType := schema.Parent.Type
	numField := parentType.NumField()
	for i := 0; i < numField; i++ {
		f := parentType.Field(i)
		nextIndex := f.Index[0] + 1
		if f.Type == reflect.TypeOf(at.Schema{}) && nextIndex < numField {
			nextField := parentType.Field(f.Index[0] + 1)
			if nextField.Type.Kind() == reflect.Slice {
				field = &nextField
				break
			}
		}
	}
	return field
}

func (b *apiPathsBuilder) buildResponse(operation *spec.Operation, annotations *annotation.Annotations, a *annotation.Annotation) {
	ao := a.Field.Value.Interface()
	atResponse := ao.(at.Response)
	if operation.Responses == nil {
		operation.Responses = new(spec.Responses)
		operation.Responses.StatusCodeResponses = make(map[int]spec.Response)
	}
	atSchema := annotation.Find(annotations, at.Schema{})

	response := spec.Response{}
	response.Description = atResponse.AtDescription
	if atSchema != nil {
		field := b.findArrayField(atSchema)

		response.Schema = b.buildSchema(atSchema, field)
	}

	// build headers
	atHeaders := annotation.FilterIn(annotations, at.Header{})
	if len(atHeaders) > 0 {
		response.Headers = b.buildHeaders(atHeaders)
	}

	operation.Responses.StatusCodeResponses[atResponse.AtCode] = response
	return
}

func (b *apiPathsBuilder) buildOperation(operation *spec.Operation, annotations *annotation.Annotations)  {
	for _, a := range annotations.Items {
		ao := a.Field.Value.Interface()
		switch ao.(type) {
		case at.Tags:
			ann := ao.(at.Tags)
			operation.Tags = append(operation.Tags, ann.AtValues...)
		case at.Consumes:
			ann := ao.(at.Consumes)
			operation.Consumes = append(operation.Consumes, ann.AtValues...)
		case at.Produces:
			ann := ao.(at.Produces)
			operation.Produces = append(operation.Produces, ann.AtValues...)
		case at.Parameter:
			b.buildParameter(operation, annotations, a)
		case at.Response:
			b.buildResponse(operation, annotations, a)
		case at.ExternalDocs:
			ann := ao.(at.ExternalDocs)
			extDoc := new(spec.ExternalDocumentation)
			extDoc.Description = ann.AtDescription
			extDoc.URL = ann.AtURL
			operation.ExternalDocs = extDoc
		}
	}

	for _, child := range annotations.Children {
		b.buildOperation(operation, child)
	}
}

func (b *apiPathsBuilder) Build(atController *annotation.Annotations, atMethod *annotation.Annotations) {

	if !annotation.ContainsChild(atMethod, at.Operation{}) {
		//log.Debugf("does not found any swagger annotations in %v", atController.Items[0].Parent.Type)
		return
	}

	method, pth := webutils.GetHttpMethod(atMethod)
	if method != "" {
		atRequestMapping := annotation.GetAnnotation(atController, at.RequestMapping{})
		if atRequestMapping != nil {
			ann := atRequestMapping.Field.Value.Interface().(at.RequestMapping)
			pth = path.Join(ann.AtValue, pth)
		}
		//log.Debugf("%v:%v", method, path)

		pathItem := b.apiInfoBuilder.Paths.Paths[pth]

		ann :=  annotation.GetAnnotation(atMethod, at.Operation{})

		atOperationInterface := ann.Field.Value.Interface()
		atOperation := atOperationInterface.(at.Operation)

		// copy values
		operation := &spec.Operation{}
		operation.ID = atOperation.AtID
		operation.Description = atOperation.AtDescription
		operation.Summary = atOperation.AtSummary
		operation.Deprecated = atOperation.AtDeprecated

		method = strings.Title(strings.ToLower(method))
		err := reflector.SetFieldValue(&pathItem, method, operation)
		if err == nil {
			b.buildOperation(operation, atMethod)

			// add new path item
			b.apiInfoBuilder.Paths.Paths[pth] = pathItem
			//log.Debug(b.openAPIDefinition.Paths.Paths[path])
		}
	}
}

func (b *apiPathsBuilder) buildHeaders(annotations []*annotation.Annotation) (headers map[string]spec.Header) {
	headers = make(map[string]spec.Header)

	for _, ann := range annotations {
		header := spec.Header{}
		atHeader := ann.Field.Value.Interface().(at.Header)
		header.Type = atHeader.AtType
		header.Description = atHeader.AtDescription
		header.Format = atHeader.AtFormat
		headers[atHeader.AtValue] = header
	}
	return
}
