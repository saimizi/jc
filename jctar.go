package jc

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
)

type JCTARConfig struct {
	info *JCConfigInfo
}

func (c JCTARConfig) Compress(infile string) (string, error) {
	var err = error(nil)
	var cmd *exec.Cmd

	parent, base := c.InFile(infile)
	outName, err := c.OutFileName(infile)
	if err != nil {
		JCLoggerErr.Print(err)
		return "", err
	}

	JCLoggerInfo.Printf("Tar %s to %s\n", infile, outName)
	c.DumpConfig()
	JCLoggerDebug.Printf("parent: %s, base: %s outName: %s", parent, base, outName)

	if parent == "." {
		cmd = exec.Command("tar", "-cf", outName, infile)
	} else {
		cmd = exec.Command("tar", "-C", parent, "-cf", outName, base)
	}

	err = cmd.Run()

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

func (c JCTARConfig) InFile(infile string) (string, string) {
	n := len(infile) - 1
	if n >= 0 && infile[n] == '/' {
		infile = infile[:n]
	}

	parent := filepath.Dir(infile)
	base := filepath.Base(infile)
	return parent, base
}

func (c JCTARConfig) DumpConfig() {
	info := *(c.info)
	JCLoggerInfo.Printf("JCTARConfig.level: %d\n", info.level)
	JCLoggerInfo.Printf("JCTARConfig.timestampOption: %d\n", info.timestampOption)
	JCLoggerInfo.Printf("JCTARConfig.collect: %v\n", info.collect)
}

func (c JCTARConfig) JCSetTimestampOption(option int) error {
	if option <= 3 && option >= 0 {
		c.info.timestampOption = option
		return nil
	}

	return errors.New(fmt.Sprintf("Invalid time stamp option %d.", option))
}

func (c JCTARConfig) EnableCollect() {
	info := c.info
	(*info).collect = true
}

func (c JCTARConfig) DisableCollect() {
	info := c.info
	(*info).collect = false
}

func (c JCTARConfig) SetCompLevel(level int) bool {
	return true
}

func NewTARConfig() (JCConfig, error) {
	var info JCConfigInfo
	config := JCTARConfig{info: &info}

	var j JCConfig = config

	return j, nil
}
