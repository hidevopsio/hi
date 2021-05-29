package swagger


import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/at"
)

const (
    // Profile is the configuration name "swagger"
	Profile = "swagger"
)

type configuration struct {
	at.AutoConfiguration
}

func newConfiguration() *configuration {
	return &configuration{}
}

func init() {
	app.Register(newConfiguration)
}

func (c *configuration) Controller(builder *apiInfoBuilder) *controller {
	return newController(builder)
}

func (c *configuration) HttpMethodSubscriber(pathsBuilder *apiPathsBuilder) *httpMethodSubscriber {
	return newHttpMethodSubscriber(pathsBuilder)
}

func (c *configuration) ApiPathsBuilder(infoBuilder *apiInfoBuilder) *apiPathsBuilder {
	return newApiPathsBuilder(infoBuilder)
}



