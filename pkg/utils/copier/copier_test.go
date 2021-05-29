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

// Notes: this source code is originally copied from https://github.com/jinzhu/copier

package copier_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/utils/copier"
	"reflect"
	"testing"
	"time"
)

type Foo struct {
	Name string
}

func (f *Foo) get(in string) string {
	return f.Name + in
}

type Bar struct {
	Name string
}

type Baz struct {
	Bar
}

type NestedFoo struct {
	Name string
	Foo  Foo
}

type NestedBar struct {
	Name string
	Foo  Foo
}

type MyString string

func TestCopier(t *testing.T) {

	testCases := []struct {
		name string
		from interface{}
		to   interface{}
		err  error
	}{
		{
			name: "should copy struct slice",
			from:   &[]Bar{},
			to: &[]Foo{{Name: "foo"}, {Name: "bar"}},
			err:  nil,
		},
		{
			name: "Should copy from Foo to Bar",
			from: &Foo{Name: "foo"},
			to:   &Bar{},
			err:  nil,
		},
		{
			name: "Should copy from Bar to Baz",
			from: &Bar{Name: "bar"},
			to:   &Baz{},
			err:  nil,
		},
		{
			name: "copy to value is unaddressable",
			from: &Foo{Name: "foo"},
			to:   nil,
			err:  errors.New("copy to value is unaddressable"),
		},
		{
			name: "copy source is invalid",
			from: nil,
			to:   &Bar{},
			err:  nil,
		},
		{
			name: "should copy directly",
			from: &Foo{Name: "foo"},
			to:   &Foo{},
			err:  nil,
		},
		{
			name: "should ignore unmatched type",
			from: []string{"a", "b"},
			to:   &Foo{Name: "foo"},
			err:  nil,
		},
		{
			name: "should copy struct slice",
			from: &[]Foo{{Name: "foo"}, {Name: "bar"}},
			to:   &[]Bar{},
			err:  nil,
		},
		{
			name: "should copy struct slice",
			from: &Foo{Name: "bar"},
			to:   &[]Bar{},
			err:  nil,
		},
		{
			name: "should copy struct slice",
			from: &NestedFoo{Name: "nested foo", Foo: Foo{Name: "foo"}},
			to:   &NestedFoo{},
			err:  nil,
		},
		{
			name: "copy source and destination are both invalid",
			from: nil,
			to:   nil,
			err:  errors.New("copy to value is unaddressable"),
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := copier.Copy(testCase.to, testCase.from)
			assert.Equal(t, testCase.err, err)
		})
	}
}


func TestCopierWithConfig(t *testing.T) {

	type src struct {
		Str string
		Int int
		Slice []string
		Foo *Foo
	}

	type dst struct {
		Str string
		Int int
		Slice []string
		Foo *Foo
	}

	type p1 struct {
		Foo *Foo
	}
	type p2 struct {
		Foo *Foo
	}

	type s1 struct {
		Foo []string
	}
	type s2 struct {
		Foo []string
	}

	type baz struct {
		Name string
	}
	type bar struct{
		Name string
		Arr []baz
	}
	type foo struct {
		Name string
		Bar bar
	}
	type fake struct {
		Name string
		Foo foo
	}

	testCases := []struct {
		name     string
		from     interface{}
		to       interface{}
		expected interface{}
		err      error
	}{
		{
			name: "should copy struct slice deeply",
			from: &fake{
				Name: "fake",
				Foo:  foo{
					Name: "foo",
					Bar:  bar{
						Name: "bar",
					},
				},
			},
			to:   &fake{
				Name: "fake",
				Foo:  foo{
					Name: "foo",
					Bar:  bar{
						Name: "bar",
						Arr: []baz{{Name: "foo"}, {Name: "bar"}, {Name: "baz"}},
					},
				},
			},
			expected: &fake{
				Name: "fake",
				Foo:  foo{
					Name: "foo",
					Bar:  bar{
						Name: "bar",
					},
				},
			},
			err:  nil,
		},
		{
			name: "should copy struct slice deeply",
			from: &fake{
				Name: "fake",
				Foo:  foo{
					Name: "foo",
					Bar:  bar{
						Name: "bar",
						Arr: []baz{{Name: "foo"}, {Name: "bar"}, {Name: "baz"}},
					},
				},
			},
			to:   &fake{},
			expected: &fake{
				Name: "fake",
				Foo:  foo{
					Name: "foo",
					Bar:  bar{
						Name: "bar",
						Arr: []baz{{Name: "foo"}, {Name: "bar"}, {Name: "baz"}},
					},
				},
			},
			err:  nil,
		},
		{
			name:     "should copy struct slice",
			from:     &[]Foo{{Name: "foo"}, {Name: "bar"}},
			to:       &[]Foo{{Name: "fb"}},
			expected: &[]Foo{{Name: "foo"}, {Name: "bar"}},
			err:      nil,
		},
		{
			name:     "should copy struct slice",
			from:     &[]Foo{{Name: "foo"}, {Name: "bar"}},
			to:       &[]Foo{},
			expected: &[]Foo{{Name: "foo"}, {Name: "bar"}},
			err:      nil,
		},
		{
			name:     "should copy struct slice",
			from:     nil,
			to:       &[]Foo{{Name: "foo"}, {Name: "bar"}},
			expected: &[]Foo{{Name: "foo"}, {Name: "bar"}},
			err:      nil,
		},
		{
			name:     "should copy struct slice",
			from:     &[]Bar{},
			to:       &[]Foo{{Name: "foo"}, {Name: "bar"}},
			expected: &[]Foo{{Name: "foo"}, {Name: "bar"}},
			err:      nil,
		},
		{
			name:     "Should copy from src to dst",
			from:     &s1{},
			to:       &s2{Foo: []string{"a", "b", "c"}},
			expected: &s2{Foo: []string{"a", "b", "c"}},
			err:      nil,
		},
		{
			name:     "Should copy from src to dst",
			from:     &s1{Foo: []string{}},
			to:       &s2{Foo: []string{"a", "b", "c"}},
			expected: &s2{Foo: []string{"a", "b", "c"}},
			err:      nil,
		},
		{
			name:     "Should copy from src to dst",
			from:     &p1{Foo: &Foo{Name: "foo"}},
			to:       &p2{Foo: &Foo{Name: "bar"}},
			expected: &p2{Foo: &Foo{Name: "foo"}},
			err:      nil,
		},
		{
			name:     "Should copy from src to dst",
			from:     &src{Str: "foo", Int: 1, Slice: []string{"a", "b", "c"}, Foo: &Foo{Name: "foo"}},
			to:       &dst{Str: "bar", Slice: []string{"one", "two", "three"}, Foo: &Foo{Name: "bar"}},
			expected: &dst{Str: "foo", Int: 1, Slice: []string{"a", "b", "c"}, Foo: &Foo{Name: "foo"}},
			err:      nil,
		},
		{
			name:     "Should copy from src to dst",
			from:     &src{Int: 1, },
			to:       &dst{Str: "bar", Slice: []string{"one", "two", "three"}, Foo: &Foo{Name: "bar"}},
			expected: &dst{Str: "bar", Int: 1, Slice: []string{"one", "two", "three"}, Foo: &Foo{Name: "bar"}},
			err:      nil,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := copier.Copy(testCase.to, testCase.from, copier.IgnoreEmptyValue)
			assert.EqualValues(t, testCase.expected, testCase.to)
			assert.Equal(t, testCase.err, err)
		})
	}
}


type User struct {
	Name     string
	Birthday *time.Time
	Nickname string
	Role     string
	Age      int32
	FakeAge  *int32
	Notes    []string
	flags    []byte
}

func (user User) DoubleAge() int32 {
	return 2 * user.Age
}

type Employee struct {
	Name      string
	Birthday  *time.Time
	Nickname  *string
	Age       int64
	FakeAge   int
	EmployeID int64
	DoubleAge int32
	SuperRule string
	Notes     []string
	flags     []byte
}

func (employee *Employee) Role(role string) {
	employee.SuperRule = "Super " + role
}

func checkEmployee(employee Employee, user User, t *testing.T, testCase string) {
	if employee.Name != user.Name {
		t.Errorf("%v: Name haven't been copied correctly.", testCase)
	}
	if employee.Nickname == nil || *employee.Nickname != user.Nickname {
		t.Errorf("%v: NickName haven't been copied correctly.", testCase)
	}
	if employee.Birthday == nil && user.Birthday != nil {
		t.Errorf("%v: Birthday haven't been copied correctly.", testCase)
	}
	if employee.Birthday != nil && user.Birthday == nil {
		t.Errorf("%v: Birthday haven't been copied correctly.", testCase)
	}
	if employee.Age != int64(user.Age) {
		t.Errorf("%v: Age haven't been copied correctly.", testCase)
	}
	if user.FakeAge != nil && employee.FakeAge != int(*user.FakeAge) {
		t.Errorf("%v: FakeAge haven't been copied correctly.", testCase)
	}
	if employee.DoubleAge != user.DoubleAge() {
		t.Errorf("%v: Copy from method doesn't work", testCase)
	}
	if employee.SuperRule != "Super "+user.Role {
		t.Errorf("%v: Copy to method doesn't work", testCase)
	}
	if !reflect.DeepEqual(employee.Notes, user.Notes) {
		t.Errorf("%v: Copy from slice doen't work", testCase)
	}
}

func TestCopyStruct(t *testing.T) {
	var fakeAge int32 = 12
	user := User{Name: "Jinzhu", Nickname: "jinzhu", Age: 18, FakeAge: &fakeAge, Role: "Admin", Notes: []string{"hello world", "welcome"}, flags: []byte{'x'}}
	employee := Employee{}

	if err := copier.Copy(employee, &user); err == nil {
		t.Errorf("Copy to unaddressable value should get error")
	}

	copier.Copy(&employee, &user)
	checkEmployee(employee, user, t, "Copy From Ptr To Ptr")

	copier.Copy(employee, &user)
	checkEmployee(employee, user, t, "Copy From Ptr To Struct")

	copier.Copy(employee, user)
	checkEmployee(employee, user, t, "Copy From Struct To Struct")

	employee2 := Employee{}
	copier.Copy(&employee2, user)
	checkEmployee(employee2, user, t, "Copy From Struct To Ptr")

	employee3 := Employee{}
	ptrToUser := &user
	copier.Copy(&employee3, &ptrToUser)
	checkEmployee(employee3, user, t, "Copy From Double Ptr To Ptr")

	employee4 := &Employee{}
	copier.Copy(&employee4, user)
	checkEmployee(*employee4, user, t, "Copy From Ptr To Double Ptr")
}

func TestCopyFromStructToSlice(t *testing.T) {
	user := User{Name: "Jinzhu", Age: 18, Role: "Admin", Notes: []string{"hello world"}}
	employees := []Employee{}

	if err := copier.Copy(employees, &user); err != nil && len(employees) != 0 {
		t.Errorf("Copy to unaddressable value should get error")
	}

	if copier.Copy(&employees, &user); len(employees) != 1 {
		t.Errorf("Should only have one elem when copy struct to slice")
	} else {
		checkEmployee(employees[0], user, t, "Copy From Struct To Slice Ptr")
	}

	employees2 := &[]Employee{}
	if copier.Copy(&employees2, user); len(*employees2) != 1 {
		t.Errorf("Should only have one elem when copy struct to slice")
	} else {
		checkEmployee((*employees2)[0], user, t, "Copy From Struct To Double Slice Ptr")
	}

	employees3 := []*Employee{}
	if copier.Copy(&employees3, user); len(employees3) != 1 {
		t.Errorf("Should only have one elem when copy struct to slice")
	} else {
		checkEmployee(*(employees3[0]), user, t, "Copy From Struct To Ptr Slice Ptr")
	}

	employees4 := &[]*Employee{}
	if copier.Copy(&employees4, user); len(*employees4) != 1 {
		t.Errorf("Should only have one elem when copy struct to slice")
	} else {
		checkEmployee(*((*employees4)[0]), user, t, "Copy From Struct To Double Ptr Slice Ptr")
	}
}

func TestCopyFromSliceToSlice(t *testing.T) {
	users := []User{{Name: "Jinzhu", Age: 18, Role: "Admin", Notes: []string{"hello world"}}, {Name: "Jinzhu2", Age: 22, Role: "Dev", Notes: []string{"hello world", "hello"}}}
	employees := []Employee{}

	if copier.Copy(&employees, users); len(employees) != 2 {
		t.Errorf("Should have two elems when copy slice to slice")
	} else {
		checkEmployee(employees[0], users[0], t, "Copy From Slice To Slice Ptr @ 1")
		checkEmployee(employees[1], users[1], t, "Copy From Slice To Slice Ptr @ 2")
	}

	employees2 := &[]Employee{}
	if copier.Copy(&employees2, &users); len(*employees2) != 2 {
		t.Errorf("Should have two elems when copy slice to slice")
	} else {
		checkEmployee((*employees2)[0], users[0], t, "Copy From Slice Ptr To Double Slice Ptr @ 1")
		checkEmployee((*employees2)[1], users[1], t, "Copy From Slice Ptr To Double Slice Ptr @ 2")
	}

	employees3 := []*Employee{}
	if copier.Copy(&employees3, users); len(employees3) != 2 {
		t.Errorf("Should have two elems when copy slice to slice")
	} else {
		checkEmployee(*(employees3[0]), users[0], t, "Copy From Slice To Ptr Slice Ptr @ 1")
		checkEmployee(*(employees3[1]), users[1], t, "Copy From Slice To Ptr Slice Ptr @ 2")
	}

	employees4 := &[]*Employee{}
	if copier.Copy(&employees4, users); len(*employees4) != 2 {
		t.Errorf("Should have two elems when copy slice to slice")
	} else {
		checkEmployee(*((*employees4)[0]), users[0], t, "Copy From Slice Ptr To Double Ptr Slice Ptr @ 1")
		checkEmployee(*((*employees4)[1]), users[1], t, "Copy From Slice Ptr To Double Ptr Slice Ptr @ 2")
	}
}

func TestEmbedded(t *testing.T) {
	type Base struct {
		BaseField1 int
		BaseField2 int
	}

	type Embedded struct {
		EmbeddedField1 int
		EmbeddedField2 int
		Base
	}

	base := Base{}
	embedded := Embedded{}
	embedded.BaseField1 = 1
	embedded.BaseField2 = 2
	embedded.EmbeddedField1 = 3
	embedded.EmbeddedField2 = 4

	copier.Copy(&base, &embedded)

	if base.BaseField1 != 1 {
		t.Error("Embedded fields not copied")
	}
}

type structSameName1 struct {
	A string
	B int64
	C time.Time
}

type structSameName2 struct {
	A string
	B time.Time
	C int64
}

func TestCopyFieldsWithSameNameButDifferentTypes(t *testing.T) {
	obj1 := structSameName1{A: "123", B: 2, C: time.Now()}
	obj2 := &structSameName2{}
	err := copier.Copy(obj2, &obj1)
	if err != nil {
		t.Error("Should not raise error")
	}

	if obj2.A != obj1.A {
		t.Errorf("Field A should be copied")
	}
}

type ScannerValue struct {
	V int
}

func (s *ScannerValue) Scan(src interface{}) error {
	return errors.New("I failed")
}

type ScannerStruct struct {
	V *ScannerValue
}

type ScannerStructTo struct {
	V *ScannerValue
}

func TestScanner(t *testing.T) {
	s := &ScannerStruct{
		V: &ScannerValue{
			V: 12,
		},
	}

	s2 := &ScannerStructTo{}

	err := copier.Copy(s2, s)
	if err != nil {
		t.Error("Should not raise error")
	}

	if s.V.V != s2.V.V {
		t.Errorf("Field V should be copied")
	}
}

func TestCopyMap(t *testing.T) {
	src := map[string]interface{}{
		"abc": map[string]interface{}{
			"str": "string value",
			"int": 123,
			"int-slice": []int{1,2,3},
			"str-slice": []string{"a","b","c"},
			"str-ifc": []interface{}{"a",100,"c"},
			"do-not-copy": nil,
		},
		"opq": map[string]interface{}{
			"str": "str value",
			"int": 666,
			"int-slice": []int{1,2,3},
			"str-slice": []string{"j","b","c"},
			"str-ifc": []interface{}{"a",333,"c"},
		},
		"def": "define",
		"nil": nil,
		"empty": "",
	}
	dst := map[string]interface{}{
		"abc": map[string]interface{}{
			"str2": "hhh",
			"int2": 456,
			"int-slice": []int{1,2,3},
			"str-slice": []string{"j","f","k"},
			"str-ifc": []interface{}{"a",100,"c"},
			"do-not-copy": []string{"x","y"},
		},
		"hi": "Hi all",
		"nil": "Nil",
	}

	copier.CopyMap(dst, src, copier.IgnoreEmptyValue)
	assert.Equal(t, 123, dst["abc"].(map[string]interface{})["int"])
	assert.Equal(t, "define", dst["def"])
}