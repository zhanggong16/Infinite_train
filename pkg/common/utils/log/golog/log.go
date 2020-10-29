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
	"bytes"
	"fmt"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

//log level, from low to high, more high means more serious
const (
	LevelTrace = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

const (
	Ltime  = 1 << iota //time format "2006/01/02 15:04:05"
	Lfile              //file.go:123
	Llevel             //[Trace|Debug|Info...]
)

var LevelName [6]string = [6]string{"TRACE", "DEBUG", "INFO", "WARN", "ERROR", "FATAL"}

const (
	LogSqlOn       = "on"
	LogSqlOff      = "off"
	TimeFormat     = "2006/01/02 15:04:05.999999999"
	maxBufPoolSize = 16
)

const (
	LogDefaultKey = "default"
	LogErrorKey   = "error"
	LogUmpKey     = "ump"
)

type Logger struct {
	sync.Mutex
	name  string
	level int
	flag  int

	handler Handler

	quit chan struct{}
	msg  chan []byte

	bufs [][]byte

	wg sync.WaitGroup

	closed bool
}

//new a logger with specified handler and flag
func New(handler Handler, flag int, name string) *Logger {
	var l = new(Logger)

	l.level = LevelInfo
	l.handler = handler
	l.name = name

	l.flag = flag

	l.quit = make(chan struct{})
	l.closed = false

	l.msg = make(chan []byte, 1024)

	l.bufs = make([][]byte, 0, 16)

	l.wg.Add(1)
	go l.run()

	return l
}

//new a default logger with specified handler and flag: Ltime|Lfile|Llevel
func NewDefault(handler Handler) *Logger {
	return New(handler, Ltime|Lfile|Llevel, LogDefaultKey)
}

func newStdHandler() *StreamHandler {
	h, err := NewStreamHandler(os.Stdout)
	if err != nil {
		fmt.Printf("NewStreamHandler error:%s\n", err.Error())
	}
	return h
}

var std = NewDefault(newStdHandler())

func Close() {
	std.Close()
}

func (l *Logger) run() {
	defer l.wg.Done()
	for {
		select {
		case msg := <-l.msg:
			_, err := l.handler.Write(msg)
			if err != nil {
				fmt.Println(err.Error())
			}
			l.putBuf(msg)
		case <-l.quit:
			if len(l.msg) == 0 {
				return
			}
		}
	}
}

func (l *Logger) popBuf() []byte {
	l.Lock()
	var buf []byte
	if len(l.bufs) == 0 {
		buf = make([]byte, 0, 1024)
	} else {
		buf = l.bufs[len(l.bufs)-1]
		l.bufs = l.bufs[0 : len(l.bufs)-1]
	}
	l.Unlock()

	return buf
}

func (l *Logger) putBuf(buf []byte) {
	l.Lock()
	if len(l.bufs) < maxBufPoolSize {
		buf = buf[0:0]
		l.bufs = append(l.bufs, buf)
	}
	l.Unlock()
}

func (l *Logger) Close() {
	if l.closed {
		return
	}
	l.closed = true

	close(l.quit)
	l.wg.Wait()
	l.quit = nil

	err := l.handler.Close()
	if err != nil {
		fmt.Printf("close file error:%s\n", err.Error())
	}
}

//set log level, any log level less than it will not log
func (l *Logger) SetLevel(level int) {
	l.level = level
}

func (l *Logger) Level() int {
	return l.level
}

func (l *Logger) Name() string {
	return l.name
}

//a low interface, maybe you can use it for your special log format
//but it may be not exported later......
func (l *Logger) Output(callDepth int, level int, format string, v ...interface{}) {
	if l.name == LogUmpKey {
		return
	}

	if l.level > level {
		return
	}

	buf := l.popBuf()

	if l.flag&Ltime > 0 {
		now := time.Now().Format(TimeFormat)
		buf = append(buf, now...)
		buf = append(buf, " - "...)
	}

	if l.flag&Llevel > 0 {
		buf = append(buf, LevelName[level]...)
		buf = append(buf, " - "...)
	}

	if l.flag&Lfile > 0 {
		/*_, file, line, ok := runtime.Caller(callDepth)
		if !ok {
			file = "???"
			line = 0
		} else {
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					file = file[i+1:]
					break
				}
			}
		}
		*/
		theCallerInfo := retrieveCallInfo(callDepth + 1)

		buf = append(buf, theCallerInfo.fileName...)
		buf = append(buf, ":["...)

		buf = strconv.AppendInt(buf, int64(theCallerInfo.line), 10)
		buf = append(buf, "] - "...)

		//module name
		buf = append(buf, ":["...)
		buf = append(buf, theCallerInfo.packageName...)
		buf = append(buf, "]"...)
		//func name
		buf = append(buf, " "...)
		buf = append(buf, theCallerInfo.funcName...)
		buf = append(buf, " "...)
	}

	s := fmt.Sprintf(format, v...)

	buf = append(buf, s...)

	if s[len(s)-1] != '\n' {

		if runtime.GOOS == "windows" {
			buf = append(buf, '\r', '\n')
		} else {
			buf = append(buf, '\n')
		}
	}

	l.msg <- buf
}

//a low interface, maybe you can use it for your special log format
//but it may be not exported later......
func (l *Logger) OutputNoFomat(msg string) {

	buf := l.popBuf()
	s := msg
	buf = append(buf, s...)
	if s[len(s)-1] != '\n' {
		if runtime.GOOS == "windows" {
			buf = append(buf, '\r', '\n')
		} else {
			buf = append(buf, '\n')
		}
	}

	l.msg <- buf
}

func SetLevel(level int) {
	std.SetLevel(level)
}

func StdLogger() *Logger {
	return std
}

func GetLevel() int {
	return std.level
}

//全局变量
var GlobalSysLoggers = []*Logger{StdLogger()}
var GlobalSqlLogger *Logger = GlobalSysLoggers[0]

func (l *Logger) Write(p []byte) (n int, err error) {
	output(LevelInfo, "web", "api", string(p), 0)
	return len(p), nil
}

func escape(s string, filterEqual bool) string {
	dest := make([]byte, 0, 2*len(s))
	for i := 0; i < len(s); i++ {
		r := s[i]
		switch r {
		case '|':
			continue
		case '%':
			dest = append(dest, '%', '%')
		case '=':
			if !filterEqual {
				dest = append(dest, '=')
			}
		default:
			dest = append(dest, r)
		}
	}

	return string(dest)
}

func OutputSql(state string, format string, v ...interface{}) {
	l := GlobalSqlLogger
	buf := l.popBuf()

	if l.flag&Ltime > 0 {
		now := time.Now().Format(TimeFormat)
		buf = append(buf, now...)
		buf = append(buf, " - "...)
	}

	if l.flag&Llevel > 0 {
		buf = append(buf, state...)
		buf = append(buf, " - "...)
	}

	s := fmt.Sprintf(format, v...)

	buf = append(buf, s...)

	if s[len(s)-1] != '\n' {
		if runtime.GOOS == "windows" {
			buf = append(buf, '\r', '\n')
		} else {
			buf = append(buf, '\n')
		}
	}

	l.msg <- buf
}

type callInfo struct {
	packageName string
	fileName    string
	funcName    string
	line        int
}

func retrieveCallInfo(callDepth int) *callInfo {

	pc, file, line, ok := runtime.Caller(callDepth)
	if !ok {
		file = "???"
		line = 0
		return &callInfo{

			packageName: "???",

			fileName: "???",

			funcName: "???",

			line: 0,
		}
	}

	_, fileName := path.Split(file)

	parts := strings.Split(runtime.FuncForPC(pc).Name(), ".")

	pl := len(parts)

	packageName := ""

	funcName := parts[pl-1]

	if parts[pl-2][0] == '(' {

		funcName = parts[pl-2] + "." + funcName

		packageName = strings.Join(parts[0:pl-2], ".")

	} else {

		packageName = strings.Join(parts[0:pl-1], ".")

	}

	return &callInfo{

		packageName: packageName,

		fileName: fileName,

		funcName: funcName,

		line: line,
	}
}

func output(level int, msg string, reqId string, args ...interface{}) {
	for _, GlobalSysLogger := range GlobalSysLoggers {
		if GlobalSysLogger.Name() == LogUmpKey {
			continue
		}
		if level < GlobalSysLogger.Level() {
			continue
		}

		num := len(args) / 2
		var argsBuff bytes.Buffer
		for i := 0; i < num; i++ {
			argsBuff.WriteString(escape(fmt.Sprintf("%v=%v", args[i*2], args[i*2+1]), false))
			if (i+1)*2 != len(args) {
				argsBuff.WriteString("|")
			}
		}
		if len(args)%2 == 1 {
			argsBuff.WriteString(escape(fmt.Sprintf("%v", args[len(args)-1]), false))
		}

		content := fmt.Sprintf(` "%s" "%s" request_id=%s`,
			msg, argsBuff.String(), reqId)

		GlobalSysLogger.Output(3, level, content)
	}
}

func outputf(level int, msg string, reqId string, args ...interface{}) {
	for _, GlobalSysLogger := range GlobalSysLoggers {
		if GlobalSysLogger.Name() == LogUmpKey {
			continue
		}
		if level < GlobalSysLogger.Level() {
			continue
		}
		var argsBuff bytes.Buffer
		argsBuff.WriteString(escape(fmt.Sprintf(msg, args...), false))
		content := fmt.Sprintf(`"%s" request_id=%s`, argsBuff.String(), reqId)

		GlobalSysLogger.Output(3, level, content)
	}
}

func outputUmp(msg string) {
	for _, GlobalSysLogger := range GlobalSysLoggers {
		if GlobalSysLogger.Name() != LogUmpKey {
			continue
		}
		GlobalSysLogger.OutputNoFomat(msg)
	}
}

func Trace(reqId string, msg string, args ...interface{}) {
	output(LevelTrace, msg, reqId, args...)
}
func Debug(reqId string, msg string, args ...interface{}) {
	output(LevelDebug, msg, reqId, args...)
}
func Info(reqId string, msg string, args ...interface{}) {
	output(LevelInfo, msg, reqId, args...)
}
func Warn(reqId string, msg string, args ...interface{}) {
	output(LevelWarn, msg, reqId, args...)
}
func Error(reqId string, msg string, args ...interface{}) {
	output(LevelError, msg, reqId, args...)
}
func Fatal(reqId string, msg string, args ...interface{}) {
	output(LevelFatal, msg, reqId, args...)
}

func Tracef(reqId string, msg string, args ...interface{}) {
	outputf(LevelTrace, msg, reqId, args...)
}
func Debugf(reqId string, msg string, args ...interface{}) {
	outputf(LevelDebug, msg, reqId, args...)
}
func Debugx(reqId string, msg string, args ...interface{}) {
	outputf(LevelDebug, msg, reqId, args...)
}
func Infof(reqId string, msg string, args ...interface{}) {
	outputf(LevelInfo, msg, reqId, args...)
}
func Infox(reqId string, msg string, args ...interface{}) {
	outputf(LevelInfo, msg, reqId, args...)
}
func Warnf(reqId string, msg string, args ...interface{}) {
	outputf(LevelWarn, msg, reqId, args...)
}
func Warnx(reqId string, msg string, args ...interface{}) {
	outputf(LevelWarn, msg, reqId, args...)
}
func Errorf(reqId string, msg string, args ...interface{}) {
	outputf(LevelError, msg, reqId, args...)
}
func Errorx(reqId string, msg string, args ...interface{}) {
	outputf(LevelError, msg, reqId, args...)
}
func Fatalf(reqId string, msg string, args ...interface{}) {
	outputf(LevelFatal, msg, reqId, args...)
}

func PrintErr(msg string, args ...interface{}) {
	outputf(LevelError, msg, "0", args...)
}

func PrintInfo(msg string, args ...interface{}) {
	outputf(LevelInfo, msg, "0", args...)
}
func Ump(msg string) {
	outputUmp(msg)
}
