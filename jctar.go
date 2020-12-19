package jc

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type JCTARConfig struct {
	info *JCConfigInfo
}

func (c JCTARConfig) Compress(infile string) (string, error) {
	var err = error(nil)
	var cmd *exec.Cmd

	JCLoggerDebug.Printf("Compress %s with tar.\n", infile)

	parent, base := JCFileNameParse(infile)
	outName, err := c.OutFileName(infile)
	if err != nil {
		JCLoggerErr.Print(err)
		return "", err
	}

	JCLoggerDebug.Printf("Tar %s to %s\n", infile, outName)
	c.DumpConfig()

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
			_, base := JCFileNameParse(outName)
			cmd := exec.Command("mv", outName, c.info.moveto)
			err = cmd.Run()
			outName = c.info.moveto + "/" + base
		}
	}

	JCLoggerDebug.Printf("Compressed file: %s.", outName)
	return outName, err
}

func (c JCTARConfig) CompressMultiFiles(pkgname string, infileDir string) (string, error) {
	var err = error(nil)
	var cmd *exec.Cmd

	JCLoggerDebug.Printf("Compress files in %s with tar.\n", infileDir)

	outName, err := c.OutFileName(pkgname)
	if err != nil {
		JCLoggerErr.Print(err)
		return "", err
	}

	JCLoggerDebug.Printf("Tar %s/* to %s\n", infileDir, outName)
	c.DumpConfig()

	var arg []string

	arg = append(arg, "-C")
	arg = append(arg, infileDir)
	arg = append(arg, "-cf")
	arg = append(arg, outName)

	filepath.Walk(infileDir,
		func(path string, info os.FileInfo, err error) error {
			if path != infileDir {
				_, base := JCFileNameParse(path)
				arg = append(arg, base)
			}
			return nil
		})

	JCLoggerDebug.Printf("tar %s\n", arg)
	cmd = exec.Command("tar", arg...)

	err = cmd.Run()
	if err == nil {
		if c.info.moveto != "" {
			JCLoggerDebug.Printf("Move %s to %s.", outName, c.info.moveto)
			_, base := JCFileNameParse(outName)
			cmd := exec.Command("mv", outName, c.info.moveto)
			err = cmd.Run()
			outName = c.info.moveto + "/" + base
		}
	}

	JCLoggerDebug.Printf("Compressed file: %s.", outName)
	return outName, err
}

func (c JCTARConfig) OutFileName(infile string) (string, error) {
	var of string

	n := len(infile) - 1
	if n >= 0 && infile[n] == '/' {
		infile = infile[:n]
	}

	ts := JCTimestamp(c.info.timestampOption)
	if ts != "" {
		of = infile + "_" + ts + ".tar"
	} else {
		of = infile + ".tar"
	}

	return of, nil
}

func (c JCTARConfig) DumpConfig() {
	JCLoggerDebug.Printf("JCTARConfig.level: %d\n", c.info.level)
	JCLoggerDebug.Printf("JCTARConfig.timestampOption: %d\n", c.info.timestampOption)
	JCLoggerDebug.Printf("JCGZIPConfig.MoveTo: %s\n", c.info.moveto)
}

func (c JCTARConfig) SetTimestampOption(option int) error {
	if option <= 3 && option >= 0 {
		c.info.timestampOption = option
		return nil
	}

	return errors.New(fmt.Sprintf("Invalid time stamp option %d.", option))
}

func (c JCTARConfig) SetCompLevel(level int) bool {
	return true
}

func (c JCTARConfig) SetMoveTo(to string) error {

	err := JCCheckMoveTo(to)
	if err == nil {
		c.info.moveto = to
	}

	return err
}

func NewTARConfig() (JCConfig, error) {
	var info JCConfigInfo
	config := JCTARConfig{info: &info}

	var j JCConfig = config

	return j, nil
}
