package jc

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"
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
	level           int
	timestampOption int
	moveto          string
}

type JCConfig interface {
	Compress(infile string) (string, error)
	SetTimestampOption(option int) error
	SetCompLevel(level int) bool
	SetMoveTo(to string) error
}

func JCCompress(c JCConfig, infile string) (string, error) {
	var s string = ""
	var err error

	switch v := c.(type) {
	case JCGZIPConfig:
		s, err = v.Compress(infile)
	case JCTARConfig:
		s, err = v.Compress(infile)
	default:
		err = errors.New("Invalid compresser")
	}

	return s, err
}

func JCSetTimestampOption(c JCConfig, option int) error {
	switch v := c.(type) {
	case JCGZIPConfig:
		return v.SetTimestampOption(option)
	case JCTARConfig:
		return v.SetTimestampOption(option)
	}

	return errors.New("UnKnown JC config")

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

func JCSetMoveTo(c JCConfig, to string) error {
	var err error

	switch v := c.(type) {
	case JCGZIPConfig:
		err = v.SetMoveTo(to)
	case JCTARConfig:
		err = v.SetMoveTo(to)
	default:
		err = nil
	}

	return err
}

func JCTimestamp(option int) string {
	t := time.Now()

	switch option {
	case 1:
		return fmt.Sprintf("%d%d%d",
			t.Year(),
			t.Month(),
			t.Day())
	case 2:
		return fmt.Sprintf("%d%d%d_%d%d%d",
			t.Year(),
			t.Month(),
			t.Day(),
			t.Hour(),
			t.Minute(),
			t.Second())
	case 3:
		return fmt.Sprintf("%d%d%d_%d%d%d_%d",
			t.Year(),
			t.Month(),
			t.Day(),
			t.Hour(),
			t.Minute(),
			t.Second(),
			t.Nanosecond())
	}

	return ""
}

func JCCheckMoveTo(to string) error {

	JCLoggerDebug.Printf("MoveTo %s\n", to)

	if to == "" {
		return fmt.Errorf("MoveTo Directory is not specified.")
	}

	fi, err := os.Stat(to)
	if err != nil {
		return fmt.Errorf("MoveTo does not exist.")
	}

	if !fi.IsDir() {
		return fmt.Errorf("MoveTo %s is not a directory.", to)
	}

	return nil
}
