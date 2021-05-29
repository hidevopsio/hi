package factory

import (
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/inject/annotation"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/system/types"
	"github.com/hidevopsio/hiboot/pkg/utils/io"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"github.com/hidevopsio/hiboot/pkg/utils/str"
	"reflect"
	"strings"
)

// MetaData is the injectable object meta data
type MetaData struct {
	Kind         string
	Name         string
	ShortName    string
	TypeName     string
	PkgName      string
	ObjectOwner  interface{}
	MetaObject   interface{}
	Type         reflect.Type
	DepNames     []string
	DepMetaData  []*MetaData
	ContextAware bool
	Instance     interface{}
}

func appendDep(deps, dep string) (retVal string) {
	if deps == "" {
		retVal = dep
	} else {
		retVal = deps + "," + dep
	}
	return
}

func findDep(objTyp, inTyp reflect.Type) (name string) {
	indInTyp := reflector.IndirectType(inTyp)
	for _, field := range reflector.DeepFields(objTyp) {
		indFieldTyp := reflector.IndirectType(field.Type)
		//log.Debugf("%v <> %v", indFieldTyp, indInTyp)
		if indFieldTyp == indInTyp {
			name = str.ToLowerCamel(field.Name)
			depPkgName := io.DirName(indFieldTyp.PkgPath())
			if depPkgName != "" {
				name = depPkgName + "." + name
			}
			break
		}
	}
	if name == "" {
		name = reflector.GetLowerCamelFullNameByType(inTyp)
	}
	return
}

func parseDependencies(object interface{}, kind string, typ reflect.Type) (deps []string) {
	var depNames string
	switch kind {
	case types.Func:
		fn := reflect.ValueOf(object)
		numIn := fn.Type().NumIn()
		for i := 0; i < numIn; i++ {
			inTyp := fn.Type().In(i)
			depNames = appendDep(depNames, findDep(typ, inTyp))
		}
	case types.Method:
		method := object.(reflect.Method)
		numIn := method.Type.NumIn()
		for i := 1; i < numIn; i++ {
			inTyp := method.Type.In(i)
			if annotation.IsAnnotation(inTyp) {
				log.Debugf("%v is annotation", inTyp.Name())
			} else {
				depNames = appendDep(depNames, findDep(typ, inTyp))
			}
		}
	default:
		// find user specific inject tag
		for _, field := range reflector.DeepFields(typ) {
			tag, ok := field.Tag.Lookup("inject")
			if ok {
				name := tag
				if name == "" {
					name = str.ToLowerCamel(field.Type.Name())
				}
				depNames = appendDep(depNames, name)
			}
		}

		// find user specific depends tag
		var depTag string
		var hasDepTag bool
		f := reflector.GetEmbeddedField(object, "", reflect.Struct)
		depTag, hasDepTag = f.Tag.Lookup("depends")
		if hasDepTag {
			depNames = appendDep(depNames, depTag)
		}
	}

	if depNames != "" {
		deps = strings.Split(depNames, ",")
	}
	return
}

func getFullName(object interface{}, n string) (name string) {
	name = n
	if object != nil {
		pkgName, typeName := reflector.GetPkgAndName(object)
		if pkgName != "" {
			if n == "" {
				name = pkgName + "." + str.ToLowerCamel(typeName)
			} else if !strings.Contains(n, ".") {
				name = pkgName + "." + name
			}
		}
	}
	return
}

// ParseParams parse parameters
func ParseParams(params ...interface{}) (name string, object interface{}) {
	numParams := len(params)
	if numParams != 0 && params[0] != nil {
		kind := reflect.TypeOf(params[0]).Kind()
		if numParams == 1 {
			if kind == reflect.String {
				name = str.LowerFirst(params[0].(string))
				return
			}
			object = params[0]
		} else {
			if kind == reflect.String {
				name = str.LowerFirst(params[0].(string))
				object = params[1]
			} else {
				t := params[0]
				switch t.(type) {
				case reflect.Type:
					name = reflector.GetLowerCamelFullNameByType(t.(reflect.Type))
				default:
					name = reflector.GetLowerCamelFullName(t)
				}
				object = params[1]
				return
			}
		}

		md := CastMetaData(object)
		if md != nil {
			name = md.Name
			name = GetObjectQualifierName(md.MetaObject, name)
		} else {
			name = getFullName(object, name)
			name = GetObjectQualifierName(object, name)
		}
	}

	return
}

// GetObjectQualifierName get the qualifier's name of object
func GetObjectQualifierName(object interface{}, name string) string {
	// overwrite with qualifier's name
	qf := annotation.GetAnnotation(object, at.Qualifier{})
	if qf != nil {
		name = qf.Field.StructField.Tag.Get("value")
		//log.Debugf("Qualifier's name: %v", name)
	}
	return name
}

// NewMetaData create new meta data
func NewMetaData(params ...interface{}) (metaData *MetaData) {
	var name string
	var shortName string
	var metaObject interface{}
	var owner interface{}
	var deps []string

	numParams := len(params)
	if numParams != 0 && params[0] != nil {
		if len(params) == 2 {
			if reflect.TypeOf(params[0]).Kind() == reflect.String {
				name = params[0].(string)
				metaObject = params[1]
			} else {
				owner = params[0]
				metaObject = params[1]
			}
		} else {
			metaObject = params[0]
		}

		switch metaObject.(type) {
		case *MetaData:
			md := metaObject.(*MetaData)
			deps = append(deps, md.DepNames...)
			metaObject = md.MetaObject
			name = md.Name
		}

		pkgName, typeName := reflector.GetPkgAndName(metaObject)
		typ := reflect.TypeOf(metaObject)
		kind := typ.Kind()
		kindName := kind.String()

		if pkgName != "" {
			shortName = str.ToLowerCamel(typeName)
			if name == "" {
				name = pkgName + "." + shortName
			} else if !strings.Contains(name, ".") {
				name = pkgName + "." + name
			}
		}
		if kind == reflect.Struct && typ.Name() == types.Method {
			kindName = types.Method
		}
		var instance interface{}
		if kindName == types.Method || kindName == types.Func {
			t, ok := reflector.GetObjectType(metaObject)
			if ok {
				typ = t
			}
		} else {
			instance = metaObject
		}

		deps = append(deps, parseDependencies(metaObject, kindName, typ)...)

		// check if it is contextAware
		contextAware := annotation.Contains(owner, at.ContextAware{}) || annotation.Contains(metaObject, at.ContextAware{})

		name = GetObjectQualifierName(metaObject, name)

		metaData = &MetaData{
			Kind:         kindName,
			PkgName:      pkgName,
			TypeName:     typeName,
			Name:         name,
			ShortName:    shortName,
			ObjectOwner:  owner,
			MetaObject:   metaObject,
			Type:         typ,
			DepNames:     deps,
			ContextAware: contextAware,
			Instance:     instance,
		}
	}

	return metaData
}

// CloneMetaData is the func for cloning meta data
func CloneMetaData(src *MetaData) (dst *MetaData) {
	dst = &MetaData{
		Kind:         src.Kind,
		Name:         src.Name,
		ShortName:    src.ShortName,
		TypeName:     src.TypeName,
		PkgName:      src.PkgName,
		ObjectOwner:  src.ObjectOwner,
		MetaObject:   src.MetaObject,
		Type:         src.Type,
		DepNames:     src.DepNames,
		DepMetaData:  src.DepMetaData,
		ContextAware: src.ContextAware,
		Instance:     src.Instance,
	}
	return dst
}
