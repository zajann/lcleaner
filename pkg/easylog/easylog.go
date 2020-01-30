package easylog

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

const (
	defaultLevel    = DEBUG
	defaultMaxSize  = 100 << (10 * 2) // MB
	defaultFilePath = "./"
)

var (
	logInit  bool = false
	level    LogLevel
	logger   *log.Logger
	fp       *os.File
	maxSize  int64
	filePath string
	fileName string
	mutex    sync.Mutex
)

func init() {
	level = defaultLevel
	maxSize = defaultMaxSize
	filePath = defaultFilePath
	fileName = fmt.Sprintf("%s.log", os.Args[0])
}

func Init(options ...func()) error {

	if logInit {
		return nil
	}

	for _, f := range options {
		f()
	}

	checkFile()

	if logger == nil {
		return errors.New("nil logger")
	}

	logInit = true

	return nil
}

func Debug(msg string, v ...interface{}) {
	checkFile()

	if level <= DEBUG {
		fmtMsg := format(DEBUG, msg)
		logger.Printf(fmtMsg, v...)
	}
}

func Info(msg string, v ...interface{}) {
	checkFile()

	if level <= INFO {
		fmtMsg := format(INFO, msg)
		logger.Printf(fmtMsg, v...)
	}
}

func Warn(msg string, v ...interface{}) {
	checkFile()

	if level <= WARN {
		fmtMsg := format(WARN, msg)
		logger.Printf(fmtMsg, v...)
	}
}

func Error(msg string, v ...interface{}) {
	checkFile()

	if level <= ERROR {
		fmtMsg := format(ERROR, msg)
		logger.Printf(fmtMsg, v...)
	}
}

func Fatal(msg string, v ...interface{}) {
	checkFile()

	if level <= FATAL {
		fmtMsg := format(FATAL, msg)
		logger.Fatalf(fmtMsg, v...)
	}
}

func checkFile() {
	mutex.Lock()
	defer mutex.Unlock()

	if fp == nil {
		if err := openFile(filePath, fileName); err != nil {
			panic(err)
		}
	}

	if isFileMax(fp) {
		closeFile()

		if err := renameFile(); err != nil {
			panic(err)
		}

		if err := openFile(filePath, fileName); err != nil {
			panic(err)
		}

		setNewLogger(fp)
	}

	if logger == nil {
		setNewLogger(fp)
	}
}

func renameFile() error {
	old := fmt.Sprintf("%s/%s", filePath, fileName)
	new := fmt.Sprintf("%s/%s.bak.%s", filePath, fileName, time.Now().Format("20060102150405"))

	if err := os.Rename(old, new); err != nil {
		return err
	}

	return nil
}

func setNewLogger(fp *os.File) {
	logger = log.New(fp, "", 0)
}

func closeFile() {
	if fp == nil {
		return
	}

	fp.Close()

	fp = nil
}

func isFileMax(fp *os.File) bool {
	info, err := fp.Stat()
	if err != nil {
		panic(err)
	}

	if info.Size() >= maxSize {
		return true
	}

	return false
}

func openFile(path string, name string) error {
	var err error

	fp, err = os.OpenFile(fmt.Sprintf("%s/%s", path, name), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	return nil
}

func SetLevel(lv LogLevel) func() {
	return func() {
		level = lv
	}
}

func SetMaxSize(size int) func() {
	return func() {
		maxSize = int64(size) << (10 * 2) // MB to Byte
	}
}

func SetFilePath(path string) func() {
	return func() {
		filePath = path
	}
}

func SetFileName(name string) func() {
	return func() {
		fileName = name
	}
}
