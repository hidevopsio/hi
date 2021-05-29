package at

type Logical string

const (
	AND Logical = "and"
	OR Logical = "or"
)

// RequiresLogical  is the annotation that annotate the method for requires logical
type RequiresLogical  struct {
	Annotation

	BaseAnnotation

	// AtLogical is the logical operator, default value is 'and'
	AtLogical Logical `json:"-" at:"logical" logical:"and"` // default value is and
}

// RequiresAuthentication is the annotation that annotate the method for authorization
type RequiresAuthentication struct {
	Annotation

	BaseAnnotation
}


// RequiresRoles is the annotation that annotate the method for requires roles
type RequiresRoles struct {
	Annotation

	RequiresLogical
}


// RequiresPermissions  is the annotation that annotate the method for requires permissions
type RequiresPermissions  struct {
	Annotation

	RequiresLogical

	// AtValues hold the permission values as an array,  e.g. `values:"user:read,team:read"`
	AtValues []string `at:"values" json:"-"`

	// AtType is for data permission, user can specify his/her own type and then implement it in middleware, e.g. `type:"pagination"`
	AtType string `json:"-" at:"type"`

	// AtIn is the input field name of query parameters, e.g. `in:"page,per_page"`; page,per_page is the default values that indicate
	AtIn []string `json:"-" at:"in"`

	// AtOut is the output field name of query parameters, e.g. `out:"expr"` <where in (1,2,3)>; expr is the default value, it can be any query parameters field name
	AtOut []string `json:"-" at:"out"`
}

// RequiresUser  is the annotation that annotate the method for requires users
type RequiresUser  struct {
	Annotation

	BaseAnnotation
}


// RequiresGuest  is the annotation that annotate the method for requires guest
type RequiresGuest  struct {
	Annotation

	BaseAnnotation
}

