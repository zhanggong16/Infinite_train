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
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"
)

//FileHandler writes log to a file.
type FileHandler struct {
	fd *os.File
}

func NewFileHandler(fileName string, flag int) (*FileHandler, error) {
	dir := path.Dir(fileName)
	mkDirErr := os.Mkdir(dir, 0777)
	if mkDirErr != nil {
		return nil, mkDirErr
	}
	f, err := os.OpenFile(fileName, flag, 0)
	if err != nil {
		return nil, err
	}

	h := new(FileHandler)

	h.fd = f

	return h, nil
}

func (h *FileHandler) Write(b []byte) (n int, err error) {
	return h.fd.Write(b)
}

func (h *FileHandler) Close() error {
	return h.fd.Close()
}

//RotatingFileHandler writes log a file, if file size exceeds maxBytes,
//it will backup current file and open a new one.
//
//max backup file number is set by backupCount, it will delete oldest if backups too many.
type RotatingFileHandler struct {
	fd *os.File

	fileName    string
	maxBytes    int
	backupCount int
}

func NewRotatingFileHandler(fileName string, maxBytes int, backupCount int) (*RotatingFileHandler, error) {
	//dir := path.Dir(fileName)
	dir, _ := filepath.Split(fileName)
	//fmt.Println("filename: ", fileName)
	//fmt.Println("dir: ", dir)
	mkDirErr := os.MkdirAll(dir, 0777)
	if mkDirErr != nil {
		if os.IsExist(mkDirErr) {
			fmt.Printf("already exist.")
		} else {
			return nil, mkDirErr
		}
	}

	h := new(RotatingFileHandler)

	if maxBytes <= 0 {
		return nil, fmt.Errorf("invalid max bytes")
	}

	h.fileName = fileName
	h.maxBytes = maxBytes
	h.backupCount = backupCount

	var err error
	h.fd, err = os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return h, nil
}

func (h *RotatingFileHandler) Write(p []byte) (n int, err error) {
	h.doRollover()
	return h.fd.Write(p)
}

func (h *RotatingFileHandler) Close() error {
	if h.fd != nil {
		return h.fd.Close()
	}
	return nil
}

func (h *RotatingFileHandler) PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {

		return false, nil

	}
	return false, err
}

func (h *RotatingFileHandler) doRollover() {
	f, err := h.fd.Stat()
	if err != nil {
		fmt.Printf("Stat error:%s\n", err.Error())
		return
	}

	if h.maxBytes <= 0 {
		return
	} else if f.Size() < int64(h.maxBytes) {
		return
	}

	if h.backupCount > 0 {
		err = h.fd.Close()
		if err != nil {
			return
		}

		for i := h.backupCount - 1; i > 0; i-- {
			sfn := fmt.Sprintf("%s.%d", h.fileName, i)
			dfn := fmt.Sprintf("%s.%d", h.fileName, i+1)
			var isSourceExist = false
			isSourceExist, err = h.PathExists(sfn)
			if err != nil {
				fmt.Printf("PathExists error:%s\n", err.Error())
				return
			}
			if isSourceExist {
				err = os.Rename(sfn, dfn)
				if err != nil {
					fmt.Printf("Rename error:%s\n", err.Error())
					return
				}
			}
		}

		dfn := fmt.Sprintf("%s.1", h.fileName)
		err = os.Rename(h.fileName, dfn)
		if err != nil {
			fmt.Printf("Rename error:%s\n", err.Error())
			return
		}
		h.fd, err = os.OpenFile(h.fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Printf("OpenFile error:%s\n", err.Error())
		}
		return
	}
}

//TimeRotatingFileHandler writes log to a file,
//it will backup current and open a new one, with a period time you sepecified.
//
//refer: http://docs.python.org/2/library/logging.handlers.html.
//same like python TimedRotatingFileHandler.
type TimeRotatingFileHandler struct {
	fd *os.File

	baseName   string
	interval   int64
	suffix     string
	rolloverAt int64
}

const (
	WhenSecond = iota
	WhenMinute
	WhenHour
	WhenDay
)

func NewTimeRotatingFileHandler(baseName string, when int8, interval int) (*TimeRotatingFileHandler, error) {
	dir := path.Dir(baseName)
	mkDirErr := os.Mkdir(dir, 0777)
	if mkDirErr != nil {
		return nil, mkDirErr
	}

	h := new(TimeRotatingFileHandler)

	h.baseName = baseName

	switch when {
	case WhenSecond:
		h.interval = 1
		h.suffix = "2006-01-02_15-04-05"
	case WhenMinute:
		h.interval = 60
		h.suffix = "2006-01-02_15-04"
	case WhenHour:
		h.interval = 3600
		h.suffix = "2006-01-02_15"
	case WhenDay:
		h.interval = 3600 * 24
		h.suffix = "2006-01-02"
	default:
		return nil, fmt.Errorf("invalid when_rotate: %d", when)
	}

	h.interval = h.interval * int64(interval)

	var err error
	h.fd, err = os.OpenFile(h.baseName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	fInfo, err := h.fd.Stat()
	if err != nil {
		fmt.Printf("State error:%s", err.Error())
	}
	h.rolloverAt = fInfo.ModTime().Unix() + h.interval

	return h, nil
}

func (h *TimeRotatingFileHandler) doRollover() {
	//refer http://hg.python.org/cpython/file/2.7/Lib/logging/handlers.py
	now := time.Now()

	if h.rolloverAt <= now.Unix() {
		fName := h.baseName + now.Format(h.suffix)
		err := h.fd.Close()
		if err != nil {
			fmt.Printf("close file error:%s", err.Error())
		}
		e := os.Rename(h.baseName, fName)
		if e != nil {
			panic(e)
		}

		h.fd, err = os.OpenFile(h.baseName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Printf("OpenFile error:%s", err.Error())
		}
		h.rolloverAt = time.Now().Unix() + h.interval
	}
}

func (h *TimeRotatingFileHandler) Write(b []byte) (n int, err error) {
	h.doRollover()
	return h.fd.Write(b)
}

func (h *TimeRotatingFileHandler) Close() error {
	return h.fd.Close()
}
