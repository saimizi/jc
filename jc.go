package jc

import (
	"errors"
	"log"
)

var (
	JCLoggerErr   *log.Logger
	JCLoggerWarn  *log.Logger
	JCLoggerInfo  *log.Logger
	JCLoggerDebug *log.Logger
)

func init() {
	JCLoggerErr = NewErrLogger()
	JCLoggerWarn = NewWarnLogger()
	JCLoggerInfo = NewInfoLogger()
	JCLoggerDebug = NewDebugLogger()
}

type JCConfigInfo struct {
	level     int
	timestamp bool
	collect   bool
	movetopwd bool
}

type JCConfig interface {
	Compress(infile string) (string, error)
	EnableTimestamp()
	DisableTimestamp()
	SetCompLevel(level int) bool
}

func JCCompress(c JCConfig, infile string) (string, error) {
	var s string = ""
	var err error

	switch v := c.(type) {
	case JCGZIPConfig:
		s, err = v.Compress(infile)
	default:
		err = errors.New("Invalid compresser")
	}

	return s, err
}

func JCEnableTimestamp(c JCConfig) {
	switch v := c.(type) {
	case JCGZIPConfig:
		v.EnableTimestamp()
	}

}

func JCDisableTimestamp(c JCConfig) {
	switch v := c.(type) {
	case JCGZIPConfig:
		v.DisableTimestamp()
	}

}

func JCSetCompLevel(c JCConfig, level int) bool {
	var ret bool

	switch v := c.(type) {
	case JCGZIPConfig:
		ret = v.SetCompLevel(level)
	default:
		ret = false
	}

	return ret
}
