module github.com/hidevopsio/hiboot

go 1.16

require (
	github.com/codahale/hdrhistogram v0.0.0-20161010025455-3a0bb77429bd // indirect
	github.com/deckarep/golang-set v1.7.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fatih/camelcase v1.0.0
	github.com/go-co-op/gocron v1.5.0
	github.com/go-openapi/runtime v0.19.4
	github.com/go-openapi/spec v0.19.2
	github.com/go-openapi/swag v0.19.5 // indirect
	github.com/go-playground/validator/v10 v10.5.0
	github.com/gojektech/valkyrie v0.0.0-20190210220504-8f62c1e7ba45
	github.com/golang/mock v1.4.4
	github.com/golang/protobuf v1.4.2
	github.com/gorilla/handlers v1.4.2
	github.com/gorilla/mux v1.7.0
	github.com/gorilla/websocket v1.4.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2
	github.com/hidevopsio/mapstructure v1.1.3-0.20190908102033-f8832fd9e307
	github.com/hidevopsio/viper v1.2.2-0.20210220025633-ccb4b202d169
	github.com/iris-contrib/formBinder v0.0.0-20190104093907-fbd5963f41e1 // indirect
	github.com/iris-contrib/go.uuid v2.0.0+incompatible
	github.com/iris-contrib/httpexpect v0.0.0-20180314041918-ebe99fcebbce
	github.com/iris-contrib/middleware v0.0.0-20171114084220-1060fbb0ce08
	github.com/kataras/golog v0.0.0-20180321173939-03be10146386
	github.com/kataras/iris v11.0.3+incompatible
	github.com/kataras/pio v0.0.0-20180511174041-a9733b5b6b83
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/moul/http2curl v1.0.0 // indirect
	github.com/onsi/ginkgo v1.10.1 // indirect
	github.com/onsi/gomega v1.7.0 // indirect
	github.com/opentracing/opentracing-go v1.1.0
	github.com/pkg/errors v0.8.1
	github.com/rakyll/statik v0.1.6
	github.com/rs/cors v1.7.0 // indirect
	github.com/sony/sonyflake v0.0.0-20160530021500-fa881fb1052b
	github.com/spf13/afero v1.1.2
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.3
	github.com/stretchr/testify v1.7.0
	github.com/uber-go/atomic v1.3.2 // indirect
	github.com/uber/jaeger-client-go v2.15.0+incompatible
	github.com/uber/jaeger-lib v1.5.0+incompatible // indirect
	github.com/valyala/bytebufferpool v1.0.0
	golang.org/x/net v0.0.0-20200421231249-e086a090c8fd
	google.golang.org/grpc v1.29.1
	gopkg.in/yaml.v2 v2.2.4
)

replace (
	github.com/go-co-op/gocron => github.com/hidevopsio/gocron v1.6.1-0.20210527020624-ab5573f2755a
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)
