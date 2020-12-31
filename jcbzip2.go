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

// BZIP2Config :
type BZIP2Config struct {
	info *ConfigInfo
}

// Name : tar compress name
func (c BZIP2Config) Name() string {
	return "BZIP2Config"
}

//DeCompress : decompress function
func (c BZIP2Config) DeCompress(infile string) (string, error) {
	var err error
	var outfilename string

	JCLoggerDebug.Printf("BZIP2Config.DeCompress: %s", infile)

	if strings.HasSuffix(infile, "bz2") {
		cmd := exec.Command("bzip2", "-d", "-k", infile)
		err = RunCmd(cmd)
	} else {
		err = errors.New("suffix is not bz2")
	}

	if err == nil {
		outfilename = strings.TrimSuffix(infile, ".bz2")
		if c.info.moveto != "" {
			JCLoggerDebug.Printf("Move %s to %s", outfilename, c.info.moveto)
			_, base := FileNameParse(outfilename)
			cmd := exec.Command("mv", outfilename, c.info.moveto)
			RunCmd(cmd)
			outfilename = c.info.moveto + "/" + base
		}
	}

	JCLoggerDebug.Printf("After BZIP2Config.DeCompress: %s", outfilename)
	return outfilename, err
}

//Compress : compress function
func (c BZIP2Config) Compress(infile string) (string, error) {
	var err = error(nil)

	JCLoggerDebug.Printf("Compress %s with bzip2.\n", infile)
	fi, err := os.Stat(infile)
	if err != nil {
		return "", err
	}

	if fi.IsDir() {
		err = errors.New(infile + " is a directory and can not be compressed by bzip2.")
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

	cmd := exec.Command("bzip2",
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

func (c BZIP2Config) outFileName(infile string) (string, error) {
	var of string

	ts := Timestamp(c.info.timestampOption)

	if ts != "" {
		of = infile + "_" + ts + ".bz2"
	} else {
		of = infile + ".bz2"
	}

	return of, nil
}

func (c BZIP2Config) dumpConfig() {
	JCLoggerDebug.Printf("BZIP2Config.level: %d\n", c.info.level)
	JCLoggerDebug.Printf("BZIP2Config.timestampOption: %d\n", c.info.timestampOption)
	JCLoggerDebug.Printf("BZIP2Config.MoveTo: %s\n", c.info.moveto)
}

//SetTimestampOption : Set timestamp
func (c BZIP2Config) SetTimestampOption(option int) error {
	if option <= 3 && option >= 0 {
		c.info.timestampOption = option
		return nil
	}

	return fmt.Errorf("Invalid time stamp opton %d", option)
}

// VaildCompLevel : vaild compress level
func (c BZIP2Config) VaildCompLevel(level int) bool {
	return (level <= 9) && (level >= 1)
}

// SetCompLevel : Set compress level
func (c BZIP2Config) SetCompLevel(level int) bool {
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
func (c BZIP2Config) SetMoveTo(to string) error {

	err := CheckMoveTo(to)
	if err == nil {
		c.info.moveto = to
	}

	return err
}

// NewBZIP2Config : New BZIP2 config object
func NewBZIP2Config() (Config, error) {
	var err error

	info := ConfigInfo{level: 0, timestampOption: 0, moveto: ""}
	config := BZIP2Config{info: &info}

	return config, err
}
