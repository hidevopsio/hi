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

package replacer

import (
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/log"
	"os"
	"reflect"
	"regexp"
	"testing"
)

type Bar struct {
	Name    string
	Profile string
	SubBar  SubBar
	SubMap  map[string]interface{}
}

type Foo struct {
	Name    string
	Project string
	Bar     Bar
	List    []string
	SubBars []SubBar
}

type SubBar struct {
	Name string
	Age  int
}

type FooBar struct {
	TheSubBar SubBar `mapstructure:"foo"`
}

func TestParseReferences(t *testing.T) {

	testCases := []struct {
		name     string
		src      interface{}
		vars     []string
		expected string
	}{
		{
			name: "test string",
			src:  &SubBar{Name: "bar"},
			vars: []string{
				"name",
			},
			expected: "bar",
		},
		{
			name: "test int",
			src:  &SubBar{Age: 18},
			vars: []string{
				"age",
			},
			expected: "18",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ref := ParseReferences(testCase.src, testCase.vars)
			assert.Equal(t, testCase.expected, ref)
		})
	}

}

func TestParseVariableName(t *testing.T) {

	re := regexp.MustCompile(`\$\{(.*?)\}`)
	matches := ParseVariables("foo ${name} bar", re)

	assert.Equal(t, "name", matches[0][1])
}

func TestGetFieldValue(t *testing.T) {
	foo := &Foo{Name: "foo"}

	field, err := GetFieldValue(foo, "Name")
	assert.Equal(t, nil, err)
	assert.Equal(t, reflect.String, field.Kind())
	assert.Equal(t, "foo", field.String())
}

func TestReplaceStringVariables(t *testing.T) {
	f := &Foo{
		Name: "foo",
		Bar: Bar{
			Name: "Hello ${name}",
		},
	}

	s := ReplaceStringVariables(f.Bar.Name, f)
	assert.Equal(t, "Hello foo", s)
}

func TestReplaceSlice(t *testing.T) {
	testData := []string{"foo", "bar", "baz"}
	f := &struct{ Options []string }{
		Options: testData,
	}
	s := ReplaceStringVariables("${options}", f)
	assert.NotEqual(t, nil, s)
	assert.Equal(t, testData, s)
}

func TestReplaceStringVariablesWithDefaultValue(t *testing.T) {
	f := &Foo{
		Name: "foo",
		Bar: Bar{
			Name: "Hello ${foo.name:foo:bar}",
		},
	}

	s := ReplaceStringVariables(f.Bar.Name, f)
	assert.Equal(t, "Hello foo:bar", s)
}

func TestReplaceMap(t *testing.T) {
	b := &Bar{
		Name: "bar",
		SubMap: map[string]interface{}{
			"name": "${name}",
			"nestedMap": map[string]interface{}{
				"name": "nested ${name}",
				"age":  18,
			},
		},
	}

	err := ReplaceMap(b.SubMap, b)
	assert.Equal(t, nil, err)
	assert.Equal(t, "bar", b.SubMap["name"])
	assert.Equal(t, "nested bar", b.SubMap["nestedMap"].(map[string]interface{})["name"])

	err = ReplaceMap(nil, nil)
	assert.Equal(t, ErrNilPointer, err)
}

func TestReplaceVariable(t *testing.T) {
	os.Setenv("FOO", "foo")
	os.Setenv("BAR", "bar")
	f := &Foo{
		Name:    "foo",
		Project: "it's ${FOO} project",
		Bar: Bar{
			Name:    "my name is ${BAR}",
			Profile: "${name}-bar",
			SubBar: SubBar{
				Name: "${bar.name}",
			},
			SubMap: map[string]interface{}{
				"barName": "${bar.name}",
				"name":    "${name}",
				"nestedMap": map[string]interface{}{
					"name": "${name}",
					"age":  18,
				},
			},
		},
		List: []string{
			"${bar.name} of ${Foo}",
		},
		SubBars: []SubBar{
			{
				Name: "${bar.name}",
			},
		},
	}
	log.Println(f)
	err := Replace(f, f)
	log.Println(f)
	assert.Equal(t, nil, err)
	assert.Equal(t, "it's foo project", f.Project)
	assert.Equal(t, "foo-bar", f.Bar.Profile)
	assert.Equal(t, "my name is bar", f.Bar.Name)
	assert.Equal(t, f.Bar.Name, f.Bar.SubBar.Name)
	assert.Equal(t, f.Name, f.Bar.SubMap["name"])
	assert.Equal(t, f.Name, f.Bar.SubMap["nestedMap"].(map[string]interface{})["name"])
}

func TestParseVariables(t *testing.T) {
	os.Setenv("FOO", "foo")
	os.Setenv("BAR", "bar")
	os.Setenv("foo.bar", "fb")
	source := "the-${FOO}-${BAR}-${foo.bar}-env-${url:http://localhost:8080}-${foo.bar:${nested.prop1}-${nested.prop2}}"

	re := regexp.MustCompile(`\$\{(.*?)\}`)

	matches := ParseVariables(source, re)

	log.Print(matches)
	assert.Equal(t, "${FOO}", matches[0][0])
	assert.Equal(t, "FOO", matches[0][1])
	assert.Equal(t, "${BAR}", matches[1][0])
	assert.Equal(t, "BAR", matches[1][1])

	for i, v := range matches {
		log.Debugf("%v, %v -> %v", i, v[0], v[1])
	}

}

func TestReplaceReferences(t *testing.T) {
	os.Setenv("FOO", "foo")
	os.Setenv("BAR", "bar")
	f := &Foo{
		Name:    "foo",
		Project: "it's ${FOO} project",
		Bar: Bar{
			Name:    "my name is ${BAR}",
			Profile: "${name}-bar",
			SubBar: SubBar{
				Name: "${bar.name}",
			},
		},
	}

	t.Run("should find top level reference, then replace it", func(t *testing.T) {
		res := ParseReferences(f, []string{"name"})
		assert.Equal(t, "foo", res)
	})

	t.Run("should find sub reference, then replace it", func(t *testing.T) {
		res := ParseReferences(f, []string{"bar", "name"})
		assert.Equal(t, "my name is ${BAR}", res)
	})

	t.Run("should return empty string when object is empty", func(t *testing.T) {
		res := ParseReferences(&Foo{}, []string{""})
		assert.Equal(t, "", res)
	})

	t.Run("should return empty string when varName is nil", func(t *testing.T) {
		res := ParseReferences(f, nil)
		assert.Equal(t, "", res)
	})
}

func TestGetReferenceValue(t *testing.T) {
	t.Run("should get reference value", func(t *testing.T) {
		fb := &FooBar{TheSubBar: SubBar{Name: "bar", Age: 18}}
		val, err := GetReferenceValue(fb, "foo")
		assert.Equal(t, nil, err)
		log.Debugf("%v : %v", fb.TheSubBar, val.Interface())
		assert.Equal(t, fb.TheSubBar, val.Interface())
	})

	t.Run("should failed to get reference value", func(t *testing.T) {
		_, err := GetReferenceValue((*FooBar)(nil), "foo")
		assert.Equal(t, ErrInvalidObject, err)
	})
}

func TestGetMatches(t *testing.T) {
	mcs := GetMatches("should find ${app.name} and ${app.role} here")
	assert.Equal(t, "app.name", mcs[0][1])
	assert.Equal(t, "app.role", mcs[1][1])
}
