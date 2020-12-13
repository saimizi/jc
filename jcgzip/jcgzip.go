package jcgzip

import (
	"bufio"
	"bytes"
	"errors"
	"jc/jclogger"
	"log"
	"os"
	"os/exec"
)

const (
	MaxCompressLevel int = 11
	MinCompressLevel int = 0
)

type JCGZIPConfig struct {
	level     int
	timestamp bool
	collect   bool
	movetopwd bool
}

var (
	JCLoggerErr   *log.Logger
	JCLoggerWarn  *log.Logger
	JCLoggerInfo  *log.Logger
	JCLoggerDebug *log.Logger
)

func init() {
	JCLoggerErr = jclogger.NewErrLogger()
	JCLoggerWarn = jclogger.NewWarnLogger()
	JCLoggerInfo = jclogger.NewInfoLogger()
	JCLoggerDebug = jclogger.NewDebugLogger()
}

func (ptrConfig *JCGZIPConfig) OutFileName(infile string) (string, error) {
	var of string

	of = infile + ".gz"

	return of, nil
}

func (ptrConfig *JCGZIPConfig) DumpConfig() {
	config := *ptrConfig
	JCLoggerInfo.Printf("JCGZIPConfig.level: %d\n", config.level)
	JCLoggerInfo.Printf("JCGZIPConfig.timestamp: %v\n", config.timestamp)
	JCLoggerInfo.Printf("JCGZIPConfig.collect: %v\n", config.collect)
}

func (ptrConfig *JCGZIPConfig) Compress(infile string) (string, error) {
	var err = error(nil)

	outName, err := ptrConfig.OutFileName(infile)
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

	config := *ptrConfig
	JCLoggerInfo.Printf("Compress %s to %s\n", infile, outName)
	ptrConfig.DumpConfig()

	compLevel := "--best"

	if config.level < 5 {
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

func (ptrConfig *JCGZIPConfig) EnableTimestamp() {
	(*ptrConfig).timestamp = true
}

func (ptrConfig *JCGZIPConfig) DisableTimestamp() {
	(*ptrConfig).timestamp = false
}

func (ptrConfig *JCGZIPConfig) EnableCollect() {
	(*ptrConfig).collect = true
}

func (ptrConfig *JCGZIPConfig) DisableCollect() {
	(*ptrConfig).collect = false
}

func vaildCompLevel(level int) bool {
	return (level <= MaxCompressLevel) && (level >= MinCompressLevel)
}

func (ptrConfig *JCGZIPConfig) VaildCompLevel(level int) bool {
	return vaildCompLevel(level)
}

func (ptrConfig *JCGZIPConfig) SetCompLevel(level int) bool {
	ret := false

	for {
		if !ptrConfig.VaildCompLevel(level) {
			break
		}

		(*ptrConfig).level, ret = level, true

		break
	}

	return ret
}

func New(level int) (*JCGZIPConfig, error) {
	if !vaildCompLevel(level) {
		return nil, errors.New("Invalid compress level.")
	}

	config := JCGZIPConfig{level: 11, timestamp: false, collect: false, movetopwd: false}
	config.level = level

	return &config, nil
}
