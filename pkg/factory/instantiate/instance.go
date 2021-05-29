package instantiate

import (
	"fmt"
	"github.com/hidevopsio/hiboot/pkg/factory"
	"github.com/hidevopsio/hiboot/pkg/utils/cmap"
)

type instance struct {
	instMap cmap.ConcurrentMap
}

func newInstance(instMap cmap.ConcurrentMap) factory.Instance {
	if instMap == nil {
		instMap = cmap.New()
	}
	return &instance{
		instMap: instMap,
	}
}

// Get get instance
func (i *instance) Get(params ...interface{}) (retVal interface{}) {
	name, obj := factory.ParseParams(params...)

	// get from instance map if external instance map does not have it
	if md, ok := i.instMap.Get(name); ok {
		metaData := factory.CastMetaData(md)
		if metaData != nil {
			switch obj.(type) {
			case factory.MetaData:
				retVal = metaData
			default:
				retVal = metaData.Instance
			}
		}
	}

	return
}

// Set save instance
func (i *instance) Set(params ...interface{}) (err error) {
	name, inst := factory.ParseParams(params...)

	metaData := factory.CastMetaData(inst)
	if metaData == nil {
		metaData = factory.NewMetaData(inst)
	}

	old, ok := i.instMap.Get(name)
	if ok {
		oldMd := factory.CastMetaData(old)
		if oldMd.Instance != nil {
			err = fmt.Errorf("instance %v is already taken", name)
			//log.Warn(err)
			return
		}
	}

	i.instMap.Set(name, metaData)
	return
}

// Items return map items
func (i *instance) Items() map[string]interface{} {
	return i.instMap.Items()
}
