# Web Application with JWT

This is the web application with JWT.


#### main.go
```go
package main

import (
	"github.com/hidevopsio/hiboot/pkg/starter/web"
	_ "github.com/hidevopsio/hiboot/examples/web/jwt/controllers"
)

func main()  {
	// create new web application and run it
	web.NewApplication().Run()
}
```

#### run unit test
```bash
go test ./...
```

#### run the example code
```bash
go run main.go
```

```bash
curl -H 'Accept-Language: cn-ZH' -H """Authorization: Bearer $(curl -d '{"username":"test","password":"123"}' -H "Content-Type: application/json" -X POST http://localhost:8080/foo/login 2>/dev/null | jq -r '.data') """ http://localhost:8080/bars/sayHello

# here is the output
{"code":200,"message":"成功","data":{"Greeting":"hello bar"}}
```