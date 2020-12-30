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

// GZIPConfig :
type GZIPConfig struct {
	info *ConfigInfo
}

// Name : tar compress name
func (c GZIPConfig) Name() string {
	return "GZIPConfig"
}

//DeCompress : decompress function
func (c GZIPConfig) DeCompress(infile string) (string, error) {
	var err error
	var outfilename string

	JCLoggerDebug.Printf("GZIPConfig.DeCompress: %s", infile)

	if strings.HasSuffix(infile, "gz") {
		cmd := exec.Command("gzip", "-d", "-k", infile)
		err = RunCmd(cmd)
	} else {
		err = errors.New("suffix is not gz")
	}

	if err == nil {
		outfilename = strings.TrimSuffix(infile, ".gz")
		if c.info.moveto != "" {
			JCLoggerDebug.Printf("Move %s to %s", outfilename, c.info.moveto)
			_, base := FileNameParse(outfilename)
			cmd := exec.Command("mv", outfilename, c.info.moveto)
			RunCmd(cmd)
			outfilename = c.info.moveto + "/" + base
		}
	}

	JCLoggerDebug.Printf("After GZIPConfig.DeCompress: %s", outfilename)
	return outfilename, err
}

//Compress : compress function
func (c GZIPConfig) Compress(infile string) (string, error) {
	var err = error(nil)

	JCLoggerDebug.Printf("Compress %s with gzip.\n", infile)
	fi, err := os.Stat(infile)
	if err != nil {
		return "", err
	}

	if fi.IsDir() {
		err = errors.New(infile + " is a directory and can not be compressed by gzip.")
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

	cmd := exec.Command("gzip",
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
			err = RunCmd(cmd)
			outName = c.info.moveto + "/" + base
		}
	}

	JCLoggerDebug.Printf("Compressed file: %s.", outName)
	return outName, err

}

func (c GZIPConfig) outFileName(infile string) (string, error) {
	var of string

	ts := Timestamp(c.info.timestampOption)

	if ts != "" {
		of = infile + "_" + ts + ".gz"
	} else {
		of = infile + ".gz"
	}

	return of, nil
}

func (c GZIPConfig) dumpConfig() {
	JCLoggerDebug.Printf("GZIPConfig.level: %d\n", c.info.level)
	JCLoggerDebug.Printf("GZIPConfig.timestampOption: %d\n", c.info.timestampOption)
	JCLoggerDebug.Printf("GZIPConfig.MoveTo: %s\n", c.info.moveto)
}

//SetTimestampOption : Set timestamp
func (c GZIPConfig) SetTimestampOption(option int) error {
	if option <= 3 && option >= 0 {
		c.info.timestampOption = option
		return nil
	}

	return fmt.Errorf("Invalid time stamp opton %d", option)
}

// VaildCompLevel : vaild compress level
func (c GZIPConfig) VaildCompLevel(level int) bool {
	return (level <= 9) && (level >= 1)
}

// SetCompLevel : Set compress level
func (c GZIPConfig) SetCompLevel(level int) bool {
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
func (c GZIPConfig) SetMoveTo(to string) error {

	err := CheckMoveTo(to)
	if err == nil {
		c.info.moveto = to
	}

	return err
}

// NewGZIPConfig : New GZIP config object
func NewGZIPConfig() (Config, error) {
	var err error

	info := ConfigInfo{level: 0, timestampOption: 0, moveto: ""}
	config := GZIPConfig{info: &info}

	return config, err
}
