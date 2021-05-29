package actuator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfiguration(t *testing.T) {
	c := newConfiguration()
	assert.NotEqual(t, nil, c)
}
