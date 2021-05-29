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

// Package gotest provides function to check whether is running in go test mode.
package gotest

import (
	"flag"
	"github.com/hidevopsio/hiboot/pkg/utils/str"
	"os"
	"strings"
)

// IsRunning return true if the go test is running
func IsRunning() (ok bool) {

	args := os.Args

	//log.Println("args: ", args)
	//log.Println("args[0]: ", args[0])

	if str.InSlice("-test.v", args) ||
		strings.Contains(args[0], ".test") {
		ok = true
	}
	return
}

// ParseArgs parse args
func ParseArgs(args []string) {

	a := os.Args[1:]
	if args != nil {
		a = args
	}

	flag.CommandLine.Parse(a)
}
