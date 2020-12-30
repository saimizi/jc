package jc

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
)

// TARConfig : tar compresser config
type TARConfig struct {
	info *ConfigInfo
}

// Name : tar compress name
func (c TARConfig) Name() string {
	return "TARConfig"
}

// DeCompress : decompress function
func (c TARConfig) DeCompress(infile string) (string, error) {
	var err error
	var outfilename string

	JCLoggerDebug.Printf("TARConfig.DeCompress: %s", infile)

	if strings.HasSuffix(infile, "tar") {
		parent, _ := FileNameParse(infile)
		cmd := exec.Command("tar", "-x", "-C", parent, "-f", infile)
		err = RunCmd(cmd)
	} else {
		err = errors.New("suffix is not tar")
	}

	if err == nil {
		outfilename = strings.TrimSuffix(infile, ".tar")
	}

	JCLoggerDebug.Printf("After TARConfig.DeCompress: %s", outfilename)
	return outfilename, err
}

// Compress : compress function
func (c TARConfig) Compress(infile string) (string, error) {
	var err = error(nil)
	var cmd *exec.Cmd

	JCLoggerDebug.Printf("Compress %s with tar.\n", infile)

	parent, base := FileNameParse(infile)
	outName, err := c.outFileName(infile)
	if err != nil {
		JCLoggerErr.Print(err)
		return "", err
	}

	JCLoggerDebug.Printf("Tar %s to %s\n", infile, outName)
	c.dumpConfig()

	if parent == "." {
		JCLoggerDebug.Printf("tar -cf  %s  %s\n", outName, infile)
		cmd = exec.Command("tar", "-cf", outName, infile)
	} else {
		JCLoggerDebug.Printf("tar -C  %s -cf  %s  %s\n", parent, outName, base)
		cmd = exec.Command("tar", "-C", parent, "-cf", outName, base)
	}

	err = cmd.Run()
	if err == nil {
		if c.info.moveto != "" {
			JCLoggerDebug.Printf("Move %s to %s.", outName, c.info.moveto)
			_, base := FileNameParse(outName)
			cmd := exec.Command("mv", outName, c.info.moveto)
			err = cmd.Run()
			outName = c.info.moveto + "/" + base
		}
	}

	JCLoggerDebug.Printf("Compressed file: %s.", outName)
	return outName, err
}

//CompressMultiFiles : compress multiple file function
func (c TARConfig) CompressMultiFiles(pkgname string, infileDir string) (string, error) {
	var err = error(nil)
	var cmd *exec.Cmd

	JCLoggerDebug.Printf("Compress files in %s with tar.\n", infileDir)

	outName, err := c.outFileName(infileDir + "/" + pkgname)
	if err != nil {
		JCLoggerErr.Print(err)
		return "", err
	}

	JCLoggerDebug.Printf("Tar %s/* to %s\n", infileDir, outName)
	c.dumpConfig()

	var arg []string

	arg = append(arg, "-C")
	arg = append(arg, infileDir)
	arg = append(arg, "-cf")
	arg = append(arg, outName)

	/*
		filepath.Walk(infileDir,
			func(path string, info os.FileInfo, err error) error {
				if path != infileDir {
					_, base := FileNameParse(path)
					arg = append(arg, base)
				}
				return nil
			})
	*/

	files, _ := ioutil.ReadDir(infileDir)
	for _, f := range files {
		arg = append(arg, f.Name())
	}

	JCLoggerDebug.Printf("tar %s\n", arg)
	cmd = exec.Command("tar", arg...)

	err = cmd.Run()
	if err == nil {
		if c.info.moveto != "" {
			JCLoggerDebug.Printf("Move %s to %s.", outName, c.info.moveto)
			_, base := FileNameParse(outName)
			cmd := exec.Command("mv", outName, c.info.moveto)
			err = cmd.Run()
			outName = c.info.moveto + "/" + base
		}
	}

	JCLoggerDebug.Printf("Compressed file: %s.", outName)
	return outName, err
}

func (c TARConfig) outFileName(infile string) (string, error) {
	var of string

	n := len(infile) - 1
	if n >= 0 && infile[n] == '/' {
		infile = infile[:n]
	}

	ts := Timestamp(c.info.timestampOption)
	if ts != "" {
		of = infile + "_" + ts + ".tar"
	} else {
		of = infile + ".tar"
	}

	return of, nil
}

func (c TARConfig) dumpConfig() {
	JCLoggerDebug.Printf("JCTARConfig.level: %d\n", c.info.level)
	JCLoggerDebug.Printf("JCTARConfig.timestampOption: %d\n", c.info.timestampOption)
	JCLoggerDebug.Printf("JCGZIPConfig.MoveTo: %s\n", c.info.moveto)
}

// SetTimestampOption : set time stamp option
func (c TARConfig) SetTimestampOption(option int) error {
	if option <= 3 && option >= 0 {
		c.info.timestampOption = option
		return nil
	}

	return fmt.Errorf("Invalid time stamp option %d", option)
}

// SetCompLevel : set compress function
func (c TARConfig) SetCompLevel(level int) bool {
	return true
}

// SetMoveTo : set move to directory name
func (c TARConfig) SetMoveTo(to string) error {

	err := CheckMoveTo(to)
	if err == nil {
		c.info.moveto = to
	}

	return err
}

// NewTARConfig : create tar compressor config
func NewTARConfig() (Config, error) {
	var info ConfigInfo
	config := TARConfig{info: &info}

	var j Config = config

	return j, nil
}
