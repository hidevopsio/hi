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

package actuator

import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/starter/logging"
	"net/http"
	"testing"
)

type fakeHealthCheckService struct {
	at.HealthCheckService
}

func (s *fakeHealthCheckService) Name() string {
	return "fake"
}

func (s *fakeHealthCheckService) Status() bool {
	return true
}

func newFakeHealthCheckService() HealthService {
	return &fakeHealthCheckService{}
}

func init() {
	app.Register(newFakeHealthCheckService)
}

func TestHealthController(t *testing.T) {
	log.SetLevel(logging.LevelDebug)
	web.RunTestApplication(t).
		Get("/health").
		Expect().Status(http.StatusOK)
}
