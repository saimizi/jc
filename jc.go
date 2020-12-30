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
	// JCLoggerErr : Error log handler
	JCLoggerErr *log.Logger
	// JCLoggerWarn : Warning log handler
	JCLoggerWarn *log.Logger

	// JCLoggerInfo : Information log handler
	JCLoggerInfo *log.Logger

	// JCLoggerDebug : Debug log handler
	JCLoggerDebug *log.Logger
)

func init() {
	JCLoggerErr = NewErrLogger()
	JCLoggerWarn = NewWarnLogger()
	JCLoggerInfo = NewInfoLogger()
	JCLoggerDebug = NewDebugLogger()
}

// ConfigInfo : Common config info
type ConfigInfo struct {
	level              int
	timestampOption    int
	moveto             string
	showOutputFileSize bool
}

// Config : interface for JC comprocesser
type Config interface {
	Name() string
	Compress(infile string) (string, error)
	DeCompress(infile string) (string, error)
	SetTimestampOption(option int) error
	SetCompLevel(level int) bool
	SetMoveTo(to string) error
}

// DeCompress : common decompress interface for JC
func DeCompress(c Config, infile string) (string, error) {
	var s string = ""
	var err error

	switch v := c.(type) {
	case GZIPConfig:
		s, err = v.DeCompress(infile)
	case TARConfig:
		s, err = v.DeCompress(infile)
	case XZConfig:
		s, err = v.DeCompress(infile)
	default:
		err = errors.New("Invalid decompresser")
	}

	return s, err
}

// Compress : common  compress interface for JC
func Compress(c Config, infile string) (string, error) {
	var s string = ""
	var err error

	switch v := c.(type) {
	case GZIPConfig:
		s, err = v.Compress(infile)
	case TARConfig:
		s, err = v.Compress(infile)
	case XZConfig:
		s, err = v.Compress(infile)
	default:
		err = errors.New("Invalid compresser")
	}

	return s, err
}

//CompressMultiFiles : Compress multi files
func CompressMultiFiles(c Config, pkgname string, infileDir string) (string, error) {
	var s string = ""
	var err error

	switch v := c.(type) {
	case TARConfig:
		s, err = v.CompressMultiFiles(pkgname, infileDir)
	default:
		err = errors.New("Invalid compresser")
	}

	return s, err
}

// SetTimestampOption : set time stamp
func SetTimestampOption(c Config, option int) error {
	switch v := c.(type) {
	case GZIPConfig:
		return v.SetTimestampOption(option)
	case TARConfig:
		return v.SetTimestampOption(option)
	case XZConfig:
		return v.SetTimestampOption(option)
	}

	return errors.New("UnKnown JC config")

}

// SetCompLevel : set compress level
func SetCompLevel(c Config, level int) bool {
	var ret bool

	switch v := c.(type) {
	case GZIPConfig:
		ret = v.SetCompLevel(level)
	case XZConfig:
		ret = v.SetCompLevel(level)
	default:
		ret = false
	}

	return ret
}

// SetMoveTo : set move to directory name
func SetMoveTo(c Config, to string) error {
	var err error

	switch v := c.(type) {
	case GZIPConfig:
		err = v.SetMoveTo(to)
	case XZConfig:
		err = v.SetMoveTo(to)
	case TARConfig:
		err = v.SetMoveTo(to)
	default:
		err = nil
	}

	return err
}

// Timestamp : create time stamp name
func Timestamp(option int) string {
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

// CheckMoveTo : check move to directory name
func CheckMoveTo(to string) error {
	if to == "" {
		return fmt.Errorf("MoveTo Directory is not specified")
	}

	fi, err := os.Stat(to)
	if err != nil {
		return fmt.Errorf("MoveTo does not exist")
	}

	if !fi.IsDir() {
		return fmt.Errorf("MoveTo %s is not a directory", to)
	}

	return nil
}

// RunCmd : run shell command helper
func RunCmd(cmd *exec.Cmd) error {
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

// RunCmdBuffer : run shell command helper
func RunCmdBuffer(cmd *exec.Cmd) ([]byte, []byte, error) {
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

// FileNameParse : file name parser
func FileNameParse(infile string) (string, string) {
	n := len(infile) - 1
	if n >= 0 && infile[n] == '/' {
		infile = infile[:n]
	}

	parent := filepath.Dir(infile)
	base := filepath.Base(infile)
	return parent, base
}
