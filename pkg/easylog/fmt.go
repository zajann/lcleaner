package easylog

import (
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	logFmt  = "TIME, LEVEL, PID, MSG"
	timeFmt = "20060102150405"
)

func format(lv LogLevel, msg string) string {

	fmt := logFmt

	fmt = strings.Replace(fmt, "TIME", time.Now().Format(timeFmt), 1)
	fmt = strings.Replace(fmt, "LEVEL", levelFmt[lv], 1)
	fmt = strings.Replace(fmt, "PID", strconv.Itoa(os.Getpid()), 1)
	fmt = strings.Replace(fmt, "MSG", msg, 1)

	return fmt
}
