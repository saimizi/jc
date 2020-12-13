package main

import (
	"errors"
	"flag"
	"jc/jcgzip"
	"jc/jclogger"
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
	JCLoggerErr = jclogger.NewErrLogger()
	JCLoggerWarn = jclogger.NewWarnLogger()
	JCLoggerInfo = jclogger.NewInfoLogger()
	JCLoggerDebug = jclogger.NewDebugLogger()
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

	if *strptrCompressCMD == "gzip" {
		config, err := jcgzip.New(*intptrCompressLevel)
		if err != nil {
			JCLoggerErr.Print(err)
			os.Exit(1)
		}

		var wg sync.WaitGroup
		wg.Add(len(infiles))
		for _, infile := range infiles {

			f := infile
			go func() {
				_, err = config.Compress(f)
				if err != nil {
					JCLoggerErr.Print(err)
				}
				wg.Done()
			}()
		}

		wg.Wait()
	}

}
