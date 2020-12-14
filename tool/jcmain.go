package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"jc"
	"log"
	"os"
	"os/exec"
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
		"tar",
		"tgz",
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

func JCCompressOne(c jc.JCConfig, infiles []string) error {
	var err error
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

	return err
}

func JCCompressTwo(c1 jc.JCConfig, c2 jc.JCConfig, infiles []string) error {
	var err error
	var wg sync.WaitGroup

	wg.Add(len(infiles))
	for _, infile := range infiles {

		f := infile
		go func() {

			f1, err := jc.JCCompress(c1, f)
			if err == nil {
				_, err = jc.JCCompress(c2, f1)
				if err != nil {
					JCLoggerErr.Print(err)
				}
				os.Remove(f1)

			} else {
				JCLoggerErr.Print(err)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	return err
}

func JCCollectionCompress(c2 jc.JCConfig, infiles []string, level int, timestampOption int) error {
	tmpdir, err := ioutil.TempDir(".", "jcpkg_")
	if err != nil {
		return err
	}

	defer func() {
		err = os.RemoveAll(tmpdir)
		if err != nil {
			JCLoggerErr.Print(err)
		}
	}()

	for _, tp := range infiles {
		cmd := exec.Command("cp", "-r", tp, tmpdir)
		err = cmd.Run()
		if err != nil {
			return err
		}
	}

	JCLoggerDebug.Printf("%s\n", tmpdir)
	t := []string{tmpdir}

	for {
		c1, err := jc.NewTARConfig()
		if err != nil {
			JCLoggerErr.Print(err)

		}

		err = jc.JCSetTimestampOption(c1, timestampOption)
		if err != nil {
			JCLoggerErr.Print(err)
			os.Exit(1)
		}

		err = JCCompressTwo(c1, c2, t)
		if err != nil {
			JCLoggerErr.Print(err)
			os.Exit(1)
		}

		break
	}
	return err
}

func main() {
	//boolptrMoveToPWD := flag.Bool("w", false, "Move the compressed file to current dir.")
	strptrCompressCMD := flag.String("c", "gzip", "Compress command.")
	intptrCompressLevel := flag.Int("l", 6, "Compress level.")
	boolptrCollect := flag.Bool("C", false, "Collect all files to create a tarball.")
	intptrTimestamp := flag.Int("t", 0, "Append time stamp to compressed file\n"+
		"0: none\n"+
		"1: Year to day\n"+
		"2: Year to seconds\n"+
		"3: Year to nanoseconds\n")

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

	if *boolptrCollect {
		var c jc.JCConfig

		if *strptrCompressCMD == "gzip" {
			c, err = jc.NewGZIPConfig(*intptrCompressLevel)
			if err != nil {
				JCLoggerErr.Print(err)
				os.Exit(1)
			}
		}

		err = JCCollectionCompress(c, infiles, *intptrCompressLevel, *intptrTimestamp)
		if err != nil {
			JCLoggerErr.Print(err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	if *strptrCompressCMD == "gzip" {
		c, err := jc.NewGZIPConfig(*intptrCompressLevel)
		if err != nil {
			JCLoggerErr.Print(err)
			os.Exit(1)
		}

		err = jc.JCSetTimestampOption(c, *intptrTimestamp)
		if err != nil {
			JCLoggerErr.Print(err)
			os.Exit(1)
		}

		err = JCCompressOne(c, infiles)
		if err != nil {
			JCLoggerErr.Print(err)
			os.Exit(1)
		}
		os.Exit(0)

	}

	if *strptrCompressCMD == "tar" {
		c, err := jc.NewTARConfig()
		if err != nil {
			JCLoggerErr.Print(err)
			os.Exit(1)
		}

		err = jc.JCSetTimestampOption(c, *intptrTimestamp)
		if err != nil {
			JCLoggerErr.Print(err)
			os.Exit(1)
		}

		err = JCCompressOne(c, infiles)
		if err != nil {
			JCLoggerErr.Print(err)
			os.Exit(1)
		}

		os.Exit(0)
	}

	if *strptrCompressCMD == "tgz" {
		c1, err := jc.NewTARConfig()
		if err != nil {
			JCLoggerErr.Print(err)
			os.Exit(1)
		}

		c2, err := jc.NewGZIPConfig(*intptrCompressLevel)
		if err != nil {
			JCLoggerErr.Print(err)
			os.Exit(1)
		}

		err = jc.JCSetTimestampOption(c1, *intptrTimestamp)
		if err != nil {
			JCLoggerErr.Print(err)
			os.Exit(1)
		}

		err = JCCompressTwo(c1, c2, infiles)
		if err != nil {
			JCLoggerErr.Print(err)
			os.Exit(1)
		}
		os.Exit(0)
	}

}
