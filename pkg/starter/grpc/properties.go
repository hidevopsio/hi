package grpc

import (
	"github.com/hidevopsio/hiboot/pkg/at"
	"time"
)

type server struct {
	Enabled bool `json:"enabled" default:"false"`

	// The network must be "tcp", "tcp4", "tcp6", "unix" or "unixpacket".
	Network string `json:"network" default:"tcp"`

	// The address can use a host name, but this is not recommended,
	// because it will create a listener for at most one of the host's IP
	// addresses.
	// If the port in the address parameter is empty or "0", as in
	// "127.0.0.1:" or "[::1]:0", a port number is automatically chosen.
	// The Addr method of Listener can be used to discover the chosen
	// port.
	// address = host:port
	// e.g. :7575 means that the address is 127.0.0.1 and port is 7575
	Host string `json:"host"`
	// server port, default is 7575
	Port string `json:"port" default:"7575"`
}

type keepAlive struct {
	Enabled bool   `json:"enabled" default:"true"`
	Delay   uint64 `json:"delay" default:"10"`
	Timeout uint64 `json:"timeout" default:"120"`
}

// ClientProperties used for grpc client injection
type ClientProperties struct {
	Host      string    `json:"host"`
	Port      string    `json:"port" default:"7575"`
	PlainText bool      `json:"plain_text" default:"true"`
	KeepAlive keepAlive `json:"keep_alive"`
}

type properties struct {
	at.ConfigurationProperties `value:"grpc"`
	at.AutoWired

	TimeoutSecond time.Duration          `json:"timeout_second"`
	Server        server                 `json:"server"`
	Client        map[string]interface{} `json:"client"`
}
