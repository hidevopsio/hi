package idgen

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/log"
)

func TestNext(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	t.Run("shoud parse mac address", func(t *testing.T) {
		ma := strings.Split("88:e9:fe:7a:86:a0", ":")
		tmp := ma[0] + ma[1]
		mid, err := strconv.ParseInt(tmp, 16, 64)
		assert.Equal(t, nil, err)
		assert.Equal(t, int64(35049), mid)
	})

	t.Run("should generate id in uint", func(t *testing.T) {
		id, err := Next()
		assert.Equal(t, nil, err)
		log.Info(id)
		assert.NotEqual(t, 0, id)
	})

	t.Run("should generate id in string", func(t *testing.T) {
		id, err := NextString()
		assert.Equal(t, nil, err)
		log.Info(id)
		assert.NotEqual(t, 0, id)
	})

}
