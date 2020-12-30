package jc

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// XZConfig :
type XZConfig struct {
	info *ConfigInfo
}

// Name : tar compress name
func (c XZConfig) Name() string {
	return "XZConfig"
}

//DeCompress : decompress function
func (c XZConfig) DeCompress(infile string) (string, error) {
	var err error
	var outfilename string

	JCLoggerDebug.Printf("XZConfig.DeCompress: %s", infile)

	if strings.HasSuffix(infile, "xz") {
		cmd := exec.Command("xz", "-d", "-k", infile)
		err = RunCmd(cmd)
	} else {
		err = errors.New("suffix is not xz")
	}

	if err == nil {
		outfilename = strings.TrimSuffix(infile, ".xz")
		if c.info.moveto != "" {
			JCLoggerDebug.Printf("Move %s to %s", outfilename, c.info.moveto)
			_, base := FileNameParse(outfilename)
			cmd := exec.Command("mv", outfilename, c.info.moveto)
			RunCmd(cmd)
			outfilename = c.info.moveto + "/" + base
		}
	}

	JCLoggerDebug.Printf("After XZConfig.DeCompress: %s", outfilename)
	return outfilename, err
}

//Compress : compress function
func (c XZConfig) Compress(infile string) (string, error) {
	var err = error(nil)

	JCLoggerDebug.Printf("Compress %s with xz.\n", infile)
	fi, err := os.Stat(infile)
	if err != nil {
		return "", err
	}

	if fi.IsDir() {
		err = errors.New(infile + " is a directory and can not be compressed by xz.")
		return "", err
	}

	outName, err := c.outFileName(infile)
	if err != nil {
		JCLoggerErr.Print(err)
		return "", err
	}

	outf, err := os.OpenFile(outName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		JCLoggerErr.Print(err)
		return "", err
	}
	defer outf.Close()
	outFileWriter := bufio.NewWriter(outf)

	JCLoggerDebug.Printf("Compress %s to %s\n", infile, outName)
	c.dumpConfig()

	cmd := exec.Command("xz",
		fmt.Sprintf("-%d", c.info.level),
		"--keep",
		"--stdout",
		infile)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		JCLoggerErr.Print(err)
		return "", err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		JCLoggerErr.Print(err)
		return "", err
	}

	finished := make(chan bool)
	go func() {
		buffer := new(bytes.Buffer)
		buffer.ReadFrom(stdout)
		buffer.WriteTo(outFileWriter)
		outFileWriter.Flush()
		finished <- true
	}()

	errBuf := make(chan []byte)
	go func() {
		b, _ := ioutil.ReadAll(stderr)
		errBuf <- b
	}()
	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("%s", <-errBuf)
	}

	<-finished

	if err == nil {
		if c.info.moveto != "" {
			_, base := FileNameParse(outName)
			cmd := exec.Command("mv", outName, c.info.moveto)
			RunCmd(cmd)
			outName = c.info.moveto + "/" + base
		}
	}

	JCLoggerDebug.Printf("Compressed file: %s.", outName)
	return outName, err

}

func (c XZConfig) outFileName(infile string) (string, error) {
	var of string

	ts := Timestamp(c.info.timestampOption)

	if ts != "" {
		of = infile + "_" + ts + ".xz"
	} else {
		of = infile + ".xz"
	}

	return of, nil
}

func (c XZConfig) dumpConfig() {
	JCLoggerDebug.Printf("XZConfig.level: %d\n", c.info.level)
	JCLoggerDebug.Printf("XZConfig.timestampOption: %d\n", c.info.timestampOption)
	JCLoggerDebug.Printf("XZConfig.MoveTo: %s\n", c.info.moveto)
}

//SetTimestampOption : Set timestamp
func (c XZConfig) SetTimestampOption(option int) error {
	if option <= 3 && option >= 0 {
		c.info.timestampOption = option
		return nil
	}

	return fmt.Errorf("Invalid time stamp opton %d", option)
}

// VaildCompLevel : vaild compress level
func (c XZConfig) VaildCompLevel(level int) bool {
	return (level <= 9) && (level >= 0)
}

// SetCompLevel : Set compress level
func (c XZConfig) SetCompLevel(level int) bool {
	info := c.info
	ret := false
	for {
		if !c.VaildCompLevel(level) {
			break
		}

		(*info).level, ret = level, true
		break
	}

	return ret
}

// SetMoveTo : set move to directory name
func (c XZConfig) SetMoveTo(to string) error {

	err := CheckMoveTo(to)
	if err == nil {
		c.info.moveto = to
	}

	return err
}

// NewXZConfig : New XZ config object
func NewXZConfig() (Config, error) {
	var err error

	info := ConfigInfo{level: 6, timestampOption: 0, moveto: ""}
	config := XZConfig{info: &info}

	return config, err
}
