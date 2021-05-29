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

package io

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/stretchr/testify/assert"
)

var testPath = ""

func init() {
	testPath = filepath.Join(os.TempDir(), "a", "b", "c.txt")
	log.SetLevel(log.DebugLevel)
}

func TestChangeWorkDir(t *testing.T) {
	wd1 := GetWorkDir()
	ChangeWorkDir("..")
	wd2 := GetWorkDir()
	ChangeWorkDir(wd1)

	assert.NotEqual(t, wd1, wd2)
}

func TestGetWorkingDir(t *testing.T) {
	wd := GetWorkDir()
	expected, err := os.Getwd()
	assert.Equal(t, nil, err)
	assert.Equal(t, expected, wd)
}

func TestGetRelativePath(t *testing.T) {
	p := GetRelativePath(1)

	assert.Equal(t, "io", DirName(p))
}

func TestIsPathNotExist(t *testing.T) {
	assert.Equal(t, true, IsPathNotExist("/TestNotExistPath"))
}

func TestCreateFile(t *testing.T) {
	err := CreateFile(os.TempDir(), "test.txt")
	assert.Equal(t, nil, err)

	bytesToWrite := "hello"
	b, err := WriterFile(os.TempDir(), "test.txt", []byte(bytesToWrite))
	assert.Equal(t, nil, err)
	assert.Equal(t, len(bytesToWrite), b)
}

func TestWriteFile(t *testing.T) {
	in := "hello, world"
	buf := bytes.NewBufferString(in)
	path := filepath.Join(os.TempDir(), "foo")
	err := os.RemoveAll(path) // remove it first
	assert.Equal(t, nil, err)
	log.Println("path: ", path)
	n, err := WriterFile(path, "test.txt", buf.Bytes())
	assert.Equal(t, nil, err)
	assert.Equal(t, len(in), n)

	if runtime.GOOS != "windows" {
		_, err = WriterFile("/should-not-have-access-permission", "test.txt", buf.Bytes())
		log.Debugf("err: %v", err)
	}
}

func TestListFiles(t *testing.T) {
	var files []string

	root := GetWorkDir()
	err := filepath.Walk(root, Visit(&files))
	assert.Equal(t, nil, err)

	for _, file := range files {
		log.Debug(file)
	}

	err = filepath.Walk("dir-does-not-exist", Visit(&files))
	assert.NotEqual(t, nil, err)
}

func TestBasename(t *testing.T) {
	b := Basename(testPath)
	assert.Equal(t, filepath.Join(os.TempDir(), "a", "b", "c"), b)

	b = Basename(filepath.Join(".a", "b", "c"))
	assert.Equal(t, filepath.Join(".a", "b", "c"), b)

	b = Basename(filepath.Join(".a"))
	assert.Equal(t, filepath.Join(".a"), b)
}

func TestFilename(t *testing.T) {
	b := Filename(testPath)
	assert.Equal(t, "c.txt", b)

	b = Filename("test.txt")
	assert.Equal(t, "test.txt", b)

	b = Filename(filepath.Join(os.TempDir(), "test.txt"))
	assert.Equal(t, "test.txt", b)
}

func TestBaseDir(t *testing.T) {
	bd := BaseDir(testPath)
	assert.NotEqual(t, testPath, bd)
}

func TestDirName(t *testing.T) {
	d := DirName("/a/b")
	assert.Equal(t, "b", d)

	d = DirName("/a")
	assert.Equal(t, "a", d)

	d = DirName("a")
	assert.Equal(t, "a", d)
}

func TestEnsureWorkDir(t *testing.T) {
	wd := GetWorkDir()

	t.Run("should return false on EnsureWorkDir when dir does not exist", func(t *testing.T) {
		res := EnsureWorkDir(1, "dir-does-not-exist")
		assert.Equal(t, false, res)
	})

	t.Run("should ensure work dir", func(t *testing.T) {
		res := EnsureWorkDir(1, filepath.Join("config", "application.yml"))
		assert.Equal(t, true, res)
		assert.NotEqual(t, wd, GetWorkDir())
	})

	t.Run("should ensure work dir on go build and run", func(t *testing.T) {
		os.Args[0] = "/private/var/folders/zw/ynx11jnx6j5gwbhkqqdqf0wm0000gn/T/___TestEnsureWorkDir_should_ensure_work_dir_on_go_build_and_run_in_hidevops_io_hiboot_pkg_utils_io"
		os.Args = append(os.Args, "go_build_")
		res := EnsureWorkDir(1, filepath.Join("config", "application.yml"))
		assert.Equal(t, true, res)
		assert.NotEqual(t, wd, GetWorkDir())
	})
}

func TestCallerInfo(t *testing.T) {
	file, line, fn := CallerInfo(1)
	assert.Contains(t, file, "io_test.go")
	assert.NotEqual(t, 0, line)
	assert.Contains(t, fn, "TestCallerInfo")
}
