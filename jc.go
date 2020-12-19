package jc

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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
	level              int
	timestampOption    int
	moveto             string
	showOutputFileSize bool
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

func JCCompressMultiFiles(c JCConfig, pkgname string, infileDir string) (string, error) {
	var s string = ""
	var err error

	switch v := c.(type) {
	case JCTARConfig:
		s, err = v.CompressMultiFiles(pkgname, infileDir)
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
		return fmt.Sprintf("%d",
			t.Nanosecond())
	}

	return ""
}

func JCCheckMoveTo(to string) error {
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

func JCRunCmd(cmd *exec.Cmd) error {
	var err error
	r, _ := cmd.StderrPipe()

	cmd.Start()
	s, _ := ioutil.ReadAll(r)
	err = cmd.Wait()

	if err != nil {
		err = fmt.Errorf("%s", s)
	}

	return err
}

func JCRunCmdBuffer(cmd *exec.Cmd) ([]byte, []byte, error) {
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	errBuf := make(chan []byte)
	go func() {
		b, _ := ioutil.ReadAll(stderr)
		errBuf <- b
	}()

	outBuf := make(chan []byte)
	go func() {
		b, _ := ioutil.ReadAll(stdout)
		outBuf <- b
	}()
	err := cmd.Run()

	return <-outBuf, <-errBuf, err
}

func JCFileNameParse(infile string) (string, string) {
	n := len(infile) - 1
	if n >= 0 && infile[n] == '/' {
		infile = infile[:n]
	}

	parent := filepath.Dir(infile)
	base := filepath.Base(infile)
	return parent, base
}
