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

package utils

import (
	"runtime"
	"os"
	"path/filepath"
	"github.com/hidevopsio/hiboot/pkg/log"
	"strings"
)

func GetWorkingDir(file string) string {
	wd, _ := os.Getwd()
	if file == "" {
		return wd
	}
	return wd
}


func GetRelativePath(skip int) string {
	_, path, _, _ := runtime.Caller(skip)

	return filepath.Dir(path)
}

func IsPathNotExist(path string) bool {
	_, err := os.Stat(path)
	isNotExist := os.IsNotExist(err)
	return isNotExist
}

func Visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		*files = append(*files, path)
		return nil
	}
}

func Basename(s string) string {
	n := strings.LastIndexByte(s, '.')
	if n > 0 {
		return s[:n]
	}
	return s
}

func Filename(s string) string {
	n := strings.LastIndexByte(s, '/')
	if n > 0 {
		return s[n + 1:]
	}
	return s
}