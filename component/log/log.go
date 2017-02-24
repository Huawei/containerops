/*
Copyright 2014 Huawei Technologies Co., Ltd. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package log

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/containerops/configure"
)

const (
	PanicLevel = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
)

const (
	nocolor = 0
	red     = 31
	green   = 32
	yellow  = 33
	blue    = 34
	gray    = 37
)

type Logger struct {
	*logrus.Logger
}

var (
	log *Logger
)

func New() *Logger {
	if log == nil {
		log = new(Logger)
		log.Logger = logrus.New()

		log.Formatter = &MyFormatter{}

		switch strings.ToUpper(strings.TrimSpace(configure.GetString("log.level"))) {
		case "PANIC":
			log.Level = logrus.PanicLevel
		case "FATAL":
			log.Level = logrus.FatalLevel
		case "ERROR":
			log.Level = logrus.ErrorLevel
		case "WARN", "WARNING":
			log.Level = logrus.WarnLevel
		case "INFO":
			log.Level = logrus.InfoLevel
		case "DEBUG":
			log.Level = logrus.DebugLevel
		default:
			log.Level = logrus.DebugLevel
		}

		logFile := getLogFile(strings.TrimSpace(configure.GetString("log.file")))
		log.Out = logFile

	}
	return log
}

func getLogFile(name string) *os.File {
	if name == "" {
		return os.Stdout
	}
	var f *os.File
	fileInfo, err := os.Stat(name)
	if err == nil {
		if fileInfo.IsDir() {
			name = name + string(os.PathSeparator) + "runtime.log"
			return getLogFile(name)
		} else {
			var flag int
			flag = os.O_RDWR | os.O_APPEND
			f, err = os.OpenFile(name, flag, 0)
		}
	} else if os.IsNotExist(err) {
		d := path.Dir(name)
		_, err = os.Stat(d)
		if os.IsNotExist(err) {
			os.MkdirAll(d, 0755)
		}
		f, err = os.Create(name)
	}
	if err != nil {
		f = os.Stdout
		fmt.Println(err)
	}
	return f
}

type MyFormatter struct {
}

func (f *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	keys := make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// if entry.Buffer != nil {
	// 	b = entry.Buffer
	// } else {
	b = &bytes.Buffer{}
	// }

	levelColor := getPrintColored(entry)
	levelText := strings.ToUpper(entry.Level.String())

	depth := getDepth()
	if depth > 0 {
		depth = depth - 1
	}
	pc, file, line, _ := runtime.Caller(depth)

	pkgandmodel := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	if len(pkgandmodel) < 2 {
		pkgandmodel = append(pkgandmodel, pkgandmodel[0])
	}

	fmtStr := ""
	if logrus.IsTerminal() {
		fmtStr = "\x1b[%dm[%s] %s%d %s [%s] %-4s\x1b[0m"
		fmt.Fprintf(b, fmtStr, levelColor, entry.Time.Format(time.RFC3339), file+":", line, pkgandmodel[0], pkgandmodel[1], levelText)
	} else {
		fmtStr = "[%s] %s%d %s [%s] %-4s"
		fmt.Fprintf(b, fmtStr, entry.Time.Format(time.RFC3339), file+":", line, pkgandmodel[0], pkgandmodel[1], levelText)
	}

	for _, k := range keys {
		v := entry.Data[k]
		if logrus.IsTerminal() {
			fmt.Fprintf(b, " \x1b[%dm%s\x1b[0m=", levelColor, k)
		} else {
			fmt.Fprintf(b, " %s", k)
		}

		f.appendValue(b, v)
	}

	fmt.Fprintf(b, " %-44s \n", entry.Message)

	return b.Bytes(), nil
}

func getPrintColored(entry *logrus.Entry) int {
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel:
		levelColor = gray
	case logrus.WarnLevel:
		levelColor = yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = red
	default:
		levelColor = blue
	}

	return levelColor
}

func getDepth() int {
	depth := 0
	hasLogPkg := false
	for true {
		_, file, _, ok := runtime.Caller(depth)
		if !ok {
			return depth
		}

		if strings.Index(file, "/github.com/Sirupsen/logrus") != -1 {
			hasLogPkg = true
		}

		if strings.Index(file, "/github.com/Sirupsen/logrus") == -1 && hasLogPkg {
			return depth
		}

		depth++
	}

	return depth
}

func (f *MyFormatter) needsQuoting(text string) bool {
	if len(text) == 0 {
		return true
	}
	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.') {
			return true
		}
	}
	return false
}

func (f *MyFormatter) appendKeyValue(b *bytes.Buffer, key string, value interface{}) {
	b.WriteString(key)
	b.WriteByte('=')
	f.appendValue(b, value)
	b.WriteByte(' ')
}

func (f *MyFormatter) appendValue(b *bytes.Buffer, value interface{}) {
	switch value := value.(type) {
	case string:
		if !f.needsQuoting(value) {
			b.WriteString(value)
		} else {
			fmt.Fprintf(b, "%s%v%s", "`", value, "`")
		}
	case error:
		errmsg := value.Error()
		if !f.needsQuoting(errmsg) {
			b.WriteString(errmsg)
		} else {
			fmt.Fprintf(b, "%s%v%s", "`", errmsg, "`")
		}
	default:
		fmt.Fprint(b, value)
	}
}
