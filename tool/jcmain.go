package main

import (
	"errors"
	"flag"
	"jc"
	"log"
	"os"
	"sync"
)

var (
	JCLoggerErr   *log.Logger
	JCLoggerWarn  *log.Logger
	JCLoggerInfo  *log.Logger
	JCLoggerDebug *log.Logger
)

func init() {
	JCLoggerErr = jc.NewErrLogger()
	JCLoggerWarn = jc.NewWarnLogger()
	JCLoggerInfo = jc.NewInfoLogger()
	JCLoggerDebug = jc.NewDebugLogger()
}

func checkCompressCmd(cmd string) bool {
	validCmd := [...]string{
		"gzip",
	}

	for _, r := range validCmd {
		if r == cmd {
			return true
		}
	}

	return false
}

func checkInFiles(files []string) error {

	if files == nil {
		return errors.New("No target file spcified")
	}

	for _, f := range files {
		_, err := os.Stat(f)
		if err != nil {
			return errors.New(f + " is not found")
		}

	}

	return nil
}

func main() {
	//boolptrMoveToPWD := flag.Bool("w", false, "Move the compressed file to current dir.")
	strptrCompressCMD := flag.String("c", "gzip", "Compress command")
	intptrCompressLevel := flag.Int("l", 6, "Compress level")

	flag.Parse()

	jcCmd := *strptrCompressCMD

	if !checkCompressCmd(jcCmd) {
		JCLoggerErr.Printf("Compress command %s is invalid\n", jcCmd)
		os.Exit(1)
	} else {
		JCLoggerInfo.Printf("Using %s\n", jcCmd)
	}

	infiles := flag.Args()

	err := checkInFiles(infiles)
	if err != nil {
		JCLoggerErr.Print(err)
		os.Exit(1)
	}

	var c jc.JCConfig

	if *strptrCompressCMD == "gzip" {
		c, err = jc.NewGZIPConfig(*intptrCompressLevel)
		if err != nil {
			JCLoggerErr.Print(err)
			os.Exit(1)
		}
	}

	var wg sync.WaitGroup
	wg.Add(len(infiles))
	for _, infile := range infiles {

		f := infile
		go func() {
			_, err = jc.JCCompress(c, f)
			if err != nil {
				JCLoggerErr.Print(err)
			}
			wg.Done()
		}()
	}
	wg.Wait()

}
