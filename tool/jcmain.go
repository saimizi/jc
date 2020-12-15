package main

import (
	"errors"
	"flag"
	"fmt"
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

	if len(files) == 0 {
		return fmt.Errorf("No input files.")
	}

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

func JCCollectionCompress(c2 jc.JCConfig,
	pkgname string,
	infiles []string,
	level int,
	timestampOption int) error {

	if pkgname == "" {
		return fmt.Errorf("Pakcage name is null.")
	}

	if len(infiles) == 0 {
		return fmt.Errorf("No input files.")
	}

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

	pkgpath := tmpdir + "/" + pkgname
	err = os.Mkdir(pkgpath, 0755)
	if err != nil {
		return err
	}

	for _, tp := range infiles {
		cmd := exec.Command("cp", "-r", tp, pkgpath+"/")
		err = cmd.Run()
		if err != nil {
			return err
		}
	}

	var f1, f2 string
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

		f1, err = jc.JCCompress(c1, pkgpath)
		if err == nil {
			f2, err = jc.JCCompress(c2, f1)
			if err != nil {
				JCLoggerErr.Print(err)
			}
			os.Remove(f1)

		} else {
			JCLoggerErr.Print(err)
		}

		break
	}

	cmd := exec.Command("mv", f2, ".")
	err = cmd.Run()
	if err != nil {
		return err
	}
	return err
}

func main() {
	//boolptrMoveToPWD := flag.Bool("w", false, "Move the compressed file to current dir.")
	strptrCompressCMD := flag.String("c", "gzip", "Compress command.")
	intptrCompressLevel := flag.Int("l", 6, "Compress level.")
	strptrCollect := flag.String("C", "", "Collect all files to create a tarball.")
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

	if *strptrCollect != "" {
		var c jc.JCConfig

		if *strptrCompressCMD == "gzip" {
			c, err = jc.NewGZIPConfig(*intptrCompressLevel)
			if err != nil {
				JCLoggerErr.Print(err)
				os.Exit(1)
			}
		}

		err = JCCollectionCompress(c,
			*strptrCollect,
			infiles,
			*intptrCompressLevel,
			*intptrTimestamp)
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
