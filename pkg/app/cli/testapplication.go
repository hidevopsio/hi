package cli

import (
	"bytes"
	"github.com/hidevopsio/hiboot/pkg/log"
	"testing"
)

// TestApplication the interface of cli test application
type TestApplication interface {
	Initialize() error
	SetProperty(name string, value ...interface{}) TestApplication
	Run(args ...string) (output string, err error)
}

type testApplication struct {
	application
}

// NewTestApplication is the test application constructor
func NewTestApplication(t *testing.T, cmd ...interface{}) TestApplication {
	a := new(testApplication)
	err := a.initialize(cmd...)
	log.Debug(err)
	err = a.build()
	log.Debug(err)
	return a
}

// SetProperty set application property
func (a *testApplication) SetProperty(name string, value ...interface{}) TestApplication {
	a.BaseApplication.SetProperty(name, value...)
	return a
}

func (a *testApplication) Run(args ...string) (output string, err error) {
	buf := new(bytes.Buffer)
	if a.root != nil {
		a.root.SetOutput(buf)
		a.root.SetArgs(args)

		_, err = a.root.ExecuteC()

		output = buf.String()
	}
	return
}
