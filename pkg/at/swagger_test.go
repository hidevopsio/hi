package at_test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/at"
	"testing"
)

func TestMarshal(t *testing.T) {
	type Foo struct {
		at.Schema
		Name string `json:"name"`
	}

	foo := &Foo{
		Name:   "foo",
	}

	b, err := json.Marshal(foo)

	assert.Equal(t, nil, err)
	assert.Equal(t, `{"name":"foo"}`, string(b))

	backFoo := &Foo{}

	err = json.Unmarshal(b, backFoo)
	assert.Equal(t, nil, err)
	assert.Equal(t, "foo", backFoo.Name)

}