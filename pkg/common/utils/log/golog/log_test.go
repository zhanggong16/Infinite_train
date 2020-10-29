// Copyright 2016 The kingshard Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package golog

import (
	//"os"
	"fmt"
	"os"
	"testing"
)

const (
	ModMySQL = "mysql"
)

func TestLog(t *testing.T) {
	SetLevel(LevelTrace)
	Error("0", ModMySQL, "select * from qing_user")
	Errorf("0", "select * from qing_user")
	Errorf("0", "sql is %s, count=%d", "select * from qing_user", 1)
	Info("0", ModMySQL, "select * from qing_user")
	Infof("0", "sql is %s", "select * from qing_user")
	Trace("0", ModMySQL, "select * from qing_user")
	Tracef("0", "sql is %s", "select * from qing_user")
	Debug("0", ModMySQL, "select * from qing_user")
	Debugf("0", "sql is %s", "select * from qing_user")
	Warn("0", ModMySQL, "select * from qing_user")
	Warnf("0", "sql is %s", "select * from qing_user")
	Fatal("0", ModMySQL, "a", "b", "c", "d", "e")
	Fatal("0", ModMySQL, "%3d", 123)
	Fatalf("0", ModMySQL+" %%3d=%3d", 123)
	fmt.Print("test finished")
	for _, GlobalSysLogger := range GlobalSysLoggers {
		GlobalSysLogger.Close()
	}
}

func TestEscape(t *testing.T) {
	r := escape("abc= %|", false)
	if r != "abc= %%" {
		t.Fatal("invalid result ", r)
	}

	r = escape("abc= %|", true)
	if r != "abc %%" {
		t.Fatal("invalid result ", r)
	}

	if r := escape("%3d", false); r != "%%3d" {
		t.Fatal("invalid result ", r)
	}
}

func TestRotatingFileLog(t *testing.T) {
	path := "/tmp/test_log"
	os.RemoveAll(path)

	os.Mkdir(path, 0777)
	fileName := path + "/test"

	h, err := NewRotatingFileHandler(fileName, 10*1, 2)
	if err != nil {
		t.Fatal(err)
	}
	GlobalSysLoggers = GlobalSysLoggers[:0]
	GlobalSysLoggers = append(GlobalSysLoggers, New(h, Lfile|Ltime|Llevel, LogDefaultKey))
	for _, GlobalSysLoggers := range GlobalSysLoggers {
		GlobalSysLoggers.SetLevel(LevelTrace)
	}

	Debug("log", "hello,world", "OK", 0, "fileName", fileName, "fileName2", fileName, "fileName3", fileName)
	Debug("log", "hello,world1", "OK", 0, "fileName", fileName, "fileName2", fileName, "fileName3", fileName)
	Debug("log", "hello,world2", "OK", 0, "fileName", fileName, "fileName2", fileName, "fileName3", fileName)
	for _, GlobalSysLoggers := range GlobalSysLoggers {
		GlobalSysLoggers.Close()
	}

	//os.RemoveAll(path)
}
