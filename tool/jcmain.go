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
	"path/filepath"
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

func checkInFiles(files []string) ([]string, bool, error) {
	var haveSameName bool

	if len(files) == 0 {
		return nil, haveSameName, fmt.Errorf("No input files.")
	}

	if files == nil {
		return nil, haveSameName, errors.New("No target file spcified")
	}

	m := make(map[string]bool)
	n := make(map[string]int)
	for _, f := range files {
		_, err := os.Stat(f)
		if err != nil {
			return nil, haveSameName, errors.New(f + " is not found")
		} else {
			m[f] = true
			base := filepath.Base(f)
			n[base]++

			JCLoggerDebug.Printf("f: %s base:%s\n", f, base)
		}
	}

	var nfiles []string
	for f, _ := range m {
		nfiles = append(nfiles, f)
	}

	for f, count := range n {
		if count > 1 {
			haveSameName = true
			JCLoggerDebug.Printf("%d files have name %s\n", count, f)
		}
	}

	return nfiles, haveSameName, nil
}

func checkMoveTo(to string) (string, error) {

	n := len(to) - 1
	if n >= 0 && to[n] == '/' {
		to = to[:n]
	}

	err := jc.JCCheckMoveTo(to)

	if err == nil {
		return to, nil
	}

	return "", err
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
	moveto string,
	timestampOption int,
	noParentDir bool) error {

	if pkgname == "" {
		return fmt.Errorf("Pakcage name is null.")
	}

	_, err := os.Stat(pkgname)
	if err == nil {
		return fmt.Errorf(" %s exists and can not be used as pakcage name.", pkgname)
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
		err = jc.JCRunCmd(cmd)
		if err != nil {
			return err
		}
	}

	var f1, f2 string
	for {
		var c1 jc.JCConfig

		c1, err = jc.NewTARConfig()
		if err != nil {
			JCLoggerErr.Print(err)

		}

		err = jc.JCSetTimestampOption(c1, timestampOption)
		if err != nil {
			JCLoggerErr.Print(err)
			os.Exit(1)
		}

		if noParentDir == true {
			f1, err = jc.JCCompressMultiFiles(c1, pkgname, pkgpath)
			if err != nil {
				JCLoggerErr.Print(err)
				break
			}
		} else {
			f1, err = jc.JCCompress(c1, pkgpath)
			if err != nil {
				JCLoggerErr.Print(err)
				break
			}
		}

		JCLoggerDebug.Printf("c2: %v\n", c2)
		if c2 != nil {
			f2, err = jc.JCCompress(c2, f1)
			if err != nil {
				JCLoggerErr.Print(err)
			}
			os.Remove(f1)
		} else {
			f2 = f1
		}

		if moveto == "" {
			JCLoggerDebug.Printf("mv %s to current directory \n", f2)
			cmd := exec.Command("mv", f2, ".")
			err = jc.JCRunCmd(cmd)
		} else {
			JCLoggerDebug.Printf("mv %s to %s.\n", f2, moveto)
			cmd := exec.Command("mv", f2, moveto)
			_, errbuf, tmperr := jc.JCRunCmdBuffer(cmd)
			if tmperr != nil {
				err = fmt.Errorf("%s", errbuf)
			}
		}

		break
	}

	return err
}

func main() {
	strptrMoveTo := flag.String("C", "", "Move the compressed file to specified dir.")
	strptrCompressCMD := flag.String("c", "tgz", "Compress command.")
	intptrCompressLevel := flag.Int("l", 6, "Compress level.")
	strptrCollect := flag.String("a", "", "Collect all files to create a tarball.")
	strptrCollectNoParentDir := flag.String("A", "", "Collect all files to create a tarball.")
	intptrTimestamp := flag.Int("t", 0, "Append time stamp to compressed file\n"+
		"0: none\n"+
		"1: Year to day\n"+
		"2: Year to seconds\n"+
		"3: nanoseconds\n")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: jc [Options] <File|Dir> [File|Dir]...\n")
		fmt.Fprintf(os.Stderr, "\nAvaible options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	infiles := flag.Args()

	infiles, haveSameName, err := checkInFiles(infiles)
	if err != nil {
		flag.Usage()
		os.Exit(1)
	}

	jcCmd := *strptrCompressCMD

	if !checkCompressCmd(jcCmd) {
		JCLoggerErr.Printf("Compress command %s is invalid\n", jcCmd)
		os.Exit(1)
	} else {
		JCLoggerDebug.Printf("Using %s\n", jcCmd)
	}

	if *strptrCollect != "" || *strptrCollectNoParentDir != "" {
		var c jc.JCConfig

		if haveSameName {
			JCLoggerErr.Print("Can not collect files that have the same name.")
			os.Exit(1)
		}

		if (*strptrCompressCMD == "gzip") || (*strptrCompressCMD == "tgz") {
			c, err = jc.NewGZIPConfig(*intptrCompressLevel)
			if err != nil {
				JCLoggerErr.Print(err)
				os.Exit(1)
			}
		}

		var noParentDir bool
		var pkgname string

		if *strptrCollectNoParentDir != "" {
			noParentDir = true
			pkgname = *strptrCollectNoParentDir
		} else {
			pkgname = *strptrCollect
		}

		err = JCCollectionCompress(c,
			pkgname,
			infiles,
			*strptrMoveTo,
			*intptrTimestamp,
			noParentDir)

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

		to, err := checkMoveTo(*strptrMoveTo)
		if err == nil {
			if haveSameName && (*intptrTimestamp < 3) {
				JCLoggerErr.Print("Can not move compressed files that have the same name.")
				os.Exit(1)
			}

			err = jc.JCSetMoveTo(c, to)
			if err != nil {
				JCLoggerErr.Print(err)
				os.Exit(1)
			}
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

		to, err := checkMoveTo(*strptrMoveTo)
		if err == nil {
			if haveSameName && (*intptrTimestamp < 3) {
				JCLoggerErr.Print("Can not move compressed files that have the same name.")
				os.Exit(1)
			}

			err = jc.JCSetMoveTo(c, to)
			if err != nil {
				JCLoggerErr.Print(err)
				os.Exit(1)
			}
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

		to, err := checkMoveTo(*strptrMoveTo)
		if err == nil {
			if haveSameName && (*intptrTimestamp < 3) {
				JCLoggerErr.Print("Can not move compressed files that have the same name.")
				os.Exit(1)
			}

			err = jc.JCSetMoveTo(c2, to)
			if err != nil {
				JCLoggerErr.Print(err)
				os.Exit(1)
			}
		}

		err = JCCompressTwo(c1, c2, infiles)
		if err != nil {
			JCLoggerErr.Print(err)
			os.Exit(1)
		}
		os.Exit(0)
	}

}
