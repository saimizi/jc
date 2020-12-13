package jc

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"os/exec"
)

const (
	MaxCompressLevel int = 11
	MinCompressLevel int = 0
)

type JCGZIPConfig struct {
	info *JCConfigInfo
}

func (c JCGZIPConfig) Compress(infile string) (string, error) {
	var err = error(nil)

	outName, err := c.OutFileName(infile)
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

	JCLoggerInfo.Printf("Compress %s to %s\n", infile, outName)
	c.DumpConfig()

	compLevel := "--best"
	info := (*c.info)

	if info.level < 5 {
		compLevel = "--fast"
	}

	cmd := exec.Command("gzip", "--keep", "--stdout", compLevel, infile)

	stdout, err := cmd.StdoutPipe()
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
		JCLoggerDebug.Print("gorutine fished")
		finished <- true
	}()

	if err := cmd.Run(); err != nil {
		JCLoggerErr.Print(err)
	}

	JCLoggerDebug.Print("Wait gorutine ")
	<-finished
	JCLoggerDebug.Print("Wait finished")

	return outName, err

}

func (c JCGZIPConfig) OutFileName(infile string) (string, error) {
	var of string

	of = infile + ".gz"

	return of, nil
}

func (c JCGZIPConfig) DumpConfig() {
	info := *(c.info)
	JCLoggerInfo.Printf("JCGZIPConfig.level: %d\n", info.level)
	JCLoggerInfo.Printf("JCGZIPConfig.timestamp: %v\n", info.timestamp)
	JCLoggerInfo.Printf("JCGZIPConfig.collect: %v\n", info.collect)
}

func (c JCGZIPConfig) EnableTimestamp() {
	info := c.info
	(*info).timestamp = true
}

func (c JCGZIPConfig) DisableTimestamp() {
	info := c.info
	(*info).timestamp = false
}

func (c JCGZIPConfig) EnableCollect() {
	info := c.info
	(*info).collect = true
}

func (c JCGZIPConfig) DisableCollect() {
	info := c.info
	(*info).collect = false
}

func vaildCompLevel(level int) bool {
	return (level <= MaxCompressLevel) && (level >= MinCompressLevel)
}

func (c JCGZIPConfig) SetCompLevel(level int) bool {
	info := c.info
	ret := false
	for {
		if !vaildCompLevel(level) {
			break
		}

		(*info).level, ret = level, true
		break
	}

	return ret
}

func NewGZIPConfig(level int) (JCConfig, error) {
	if !vaildCompLevel(level) {
		return nil, errors.New("Invalid compress level.")
	}

	info := JCConfigInfo{level: level, timestamp: false, collect: false, movetopwd: false}

	config := JCGZIPConfig{info: &info}

	var j JCConfig = config

	return j, nil
}
