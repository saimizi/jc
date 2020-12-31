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
	"strings"
	"sync"
)

var (
	//JCLoggerErr is log function
	JCLoggerErr *log.Logger

	//JCLoggerWarn is log function
	JCLoggerWarn *log.Logger

	//JCLoggerInfo is log function
	JCLoggerInfo *log.Logger

	//JCLoggerDebug is log function
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
		"xz",
		"txz",
		"bzip2",
		"tbz2",
	}

	for _, r := range validCmd {
		if r == cmd {
			return true
		}
	}

	return false
}

func checkDecompressInFiles(files []string) ([]string, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("no input files")
	}

	m := make(map[string]bool)
	var nfiles []string
	for _, f := range files {
		if !strings.HasSuffix(f, "gz") &&
			!strings.HasSuffix(f, "tar") &&
			!strings.HasSuffix(f, "bz2") &&
			!strings.HasSuffix(f, "xz") &&
			!strings.HasSuffix(f, "tar.bz2") &&
			!strings.HasSuffix(f, "tar.xz") &&
			!strings.HasSuffix(f, "tar.gz") {

			return nil, fmt.Errorf(f + " : Invalid suffix\n")
		}

		_, err := os.Stat(f)
		if err != nil {
			return nil, errors.New(f + " is not found")
		}

		if m[f] == false {
			nfiles = append(nfiles, f)
			m[f] = true
		}
	}

	return nfiles, nil
}

func checkInFiles(files []string) ([]string, bool, error) {
	var haveSameName bool

	if len(files) == 0 {
		return nil, haveSameName, fmt.Errorf("no input files")
	}

	if files == nil {
		return nil, haveSameName, errors.New("no target file spcified")
	}

	m := make(map[string]bool)
	n := make(map[string]int)
	for _, f := range files {
		_, err := os.Stat(f)
		if err != nil {
			return nil, haveSameName, errors.New(f + " is not found")
		}

		m[f] = true
		base := filepath.Base(f)
		n[base]++

		JCLoggerDebug.Printf("f: %s base:%s\n", f, base)
	}

	var nfiles []string
	for f := range m {
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

	err := jc.CheckMoveTo(to)

	if err == nil {
		return to, nil
	}

	return "", err
}

func getConfig(infile string) jc.Config {

	var j jc.Config

	JCLoggerDebug.Printf("getConfig: check %s\n", infile)

	if infile == "" {
		return nil
	}

	if strings.HasSuffix(infile, "gz") {
		j, _ = jc.NewGZIPConfig()
	}

	if strings.HasSuffix(infile, "tar") {
		j, _ = jc.NewTARConfig()
	}

	if strings.HasSuffix(infile, "xz") {
		j, _ = jc.NewXZConfig()
	}

	if strings.HasSuffix(infile, "bz2") {
		j, _ = jc.NewBZIP2Config()
	}

	if j != nil {
		JCLoggerDebug.Printf("getConfig: find compresser %s\n", j.Name())
	} else {
		JCLoggerDebug.Printf("getConfig: Not find any compresser\n")
	}
	return j

}

// JCDecompressOne : decompress function
func JCDecompressOne(infile string, to string) (string, error) {
	var err error
	f := infile

	var tmpfiles []string

	defer func() {
		for _, tf := range tmpfiles {
			JCLoggerDebug.Printf("Remove tmpfile %s\n", tf)
			os.Remove(tf)
		}
	}()

	for j := getConfig(f); j != nil; j = getConfig(f) {

		if j != nil && f != infile {
			tmpfiles = append(tmpfiles, f)
		}

		if to != "" {
			j.SetMoveTo(to)
		}

		JCLoggerDebug.Printf("Decompress %s\n", f)
		f, err = j.DeCompress(f)
	}
	return infile, err
}

// JCDecompress : decompress function
func JCDecompress(infiles []string, to string) error {

	var err error
	var wg sync.WaitGroup

	if to != "" {
		JCLoggerWarn.Printf("Move decompressed files may overried local files")
	}

	wg.Add(len(infiles))
	for _, inf := range infiles {

		f := inf
		go func() {
			_, err = JCDecompressOne(f, to)
			if err != nil {
				JCLoggerErr.Print(err)
			}
			wg.Done()
		}()
	}

	wg.Wait()
	return err
}

//JCCompressOne : a one step compress function
func JCCompressOne(c jc.Config, infiles []string) error {
	var err error
	var wg sync.WaitGroup

	wg.Add(len(infiles))
	for _, infile := range infiles {

		f := infile
		go func() {
			_, err = jc.Compress(c, f)
			if err != nil {
				JCLoggerErr.Print(err)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	return err
}

//JCCompressTwo is a two step compress function
func JCCompressTwo(c1 jc.Config, c2 jc.Config, infiles []string) error {
	var err error
	var wg sync.WaitGroup

	wg.Add(len(infiles))
	for _, infile := range infiles {

		f := infile
		go func() {

			f1, err := jc.Compress(c1, f)
			if err == nil {
				_, err = jc.Compress(c2, f1)
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

// JCCollectionCompress is a compress function for collection
func JCCollectionCompress(c2 jc.Config,
	pkgname string,
	infiles []string,
	moveto string,
	timestampOption int,
	noParentDir bool) error {

	if pkgname == "" {
		return fmt.Errorf("Package name is null")
	}

	_, err := os.Stat(pkgname)
	if err == nil {
		return fmt.Errorf(" %s exists and can not be used as pakcage name", pkgname)
	}

	if len(infiles) == 0 {
		return fmt.Errorf("No input files")
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
		err = jc.RunCmd(cmd)
		if err != nil {
			return err
		}
	}

	var f1, f2 string
	for {
		var c1 jc.Config

		c1, err = jc.NewTARConfig()
		if err != nil {
			JCLoggerErr.Print(err)

		}

		err = jc.SetTimestampOption(c1, timestampOption)
		if err != nil {
			JCLoggerErr.Print(err)
			os.Exit(1)
		}

		if noParentDir == true {
			f1, err = jc.CompressMultiFiles(c1, pkgname, pkgpath)
			if err != nil {
				JCLoggerErr.Print(err)
				break
			}
		} else {
			f1, err = jc.Compress(c1, pkgpath)
			if err != nil {
				JCLoggerErr.Print(err)
				break
			}
		}

		JCLoggerDebug.Printf("c2: %v\n", c2)
		if c2 != nil {
			f2, err = jc.Compress(c2, f1)
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
			err = jc.RunCmd(cmd)
		} else {
			JCLoggerDebug.Printf("mv %s to %s.\n", f2, moveto)
			cmd := exec.Command("mv", f2, moveto)
			_, errbuf, tmperr := jc.RunCmdBuffer(cmd)
			if tmperr != nil {
				err = fmt.Errorf("%s", errbuf)
			}
		}

		break
	}

	return err
}

func main() {
	boolptrDecompress := flag.Bool("d", false, "Decompress.")
	strptrMoveTo := flag.String("C", ".", "Move the compressed file to specified dir.")
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
		fmt.Fprintf(os.Stderr, "\nUsage: jc [Options] <File|Dir> [File|Dir]...\n")
		fmt.Fprintf(os.Stderr, "Avaible options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "Avaible compress commands:\n")
		fmt.Fprintf(os.Stderr, "\t gzip (.gz)\n")
		fmt.Fprintf(os.Stderr, "\t bzip2 (.bz2)\n")
		fmt.Fprintf(os.Stderr, "\t xz (.xz)\n")
		fmt.Fprintf(os.Stderr, "\t tar (.tar)\n")
		fmt.Fprintf(os.Stderr, "\t tgz (.tar.gz)\n")
		fmt.Fprintf(os.Stderr, "\t txz (.tar.xz)\n")
		fmt.Fprintf(os.Stderr, "\t tbz2 (.tar.bz2)\n")
		fmt.Fprintf(os.Stderr, "\n")
	}

	flag.Parse()

	infiles := flag.Args()

	if *boolptrDecompress {
		JCLoggerDebug.Printf("Decompress")
		infiles, err := checkDecompressInFiles(infiles)
		if err != nil {
			JCLoggerErr.Println(err)
			os.Exit(1)
		}

		to, err := checkMoveTo(*strptrMoveTo)
		if err != nil {
			JCLoggerErr.Print(err)
			os.Exit(1)
		}
		err = JCDecompress(infiles, to)
		if err != nil {
			JCLoggerErr.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	infiles, haveSameName, err := checkInFiles(infiles)
	if err != nil {
		JCLoggerErr.Println(err)
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
		var c jc.Config

		if haveSameName {
			JCLoggerErr.Print("Can not collect files that have the same name.")
			os.Exit(1)
		}

		if (*strptrCompressCMD == "gzip") || (*strptrCompressCMD == "tgz") {
			c, err = jc.NewGZIPConfig()
			if err != nil {
				JCLoggerErr.Print(err)
				os.Exit(1)
			}
		}

		if (*strptrCompressCMD == "xz") || (*strptrCompressCMD == "txz") {
			c, err = jc.NewXZConfig()
			if err != nil {
				JCLoggerErr.Print(err)
				os.Exit(1)
			}
		}

		if (*strptrCompressCMD == "bz2") || (*strptrCompressCMD == "tbz2") {
			c, err = jc.NewBZIP2Config()
			if err != nil {
				JCLoggerErr.Print(err)
				os.Exit(1)
			}
		}

		if c == nil {
			JCLoggerErr.Printf("Invalid compressor %s", *strptrCompressCMD)
			os.Exit(1)
		}

		if !c.SetCompLevel(*intptrCompressLevel) {
			JCLoggerWarn.Printf("Set compress level %d to %s failed",
				*intptrCompressLevel,
				c.Name())
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
		c, err := jc.NewGZIPConfig()
		if err != nil {
			JCLoggerErr.Print(err)
			os.Exit(1)
		}

		if !c.SetCompLevel(*intptrCompressLevel) {
			JCLoggerWarn.Printf("Set compress level %d to %s failed",
				*intptrCompressLevel,
				c.Name())
		}

		err = jc.SetTimestampOption(c, *intptrTimestamp)
		if err != nil {
			JCLoggerErr.Print(err)
			os.Exit(1)
		}

		to, err := checkMoveTo(*strptrMoveTo)
		if err == nil {
			if haveSameName && (*intptrTimestamp < 3) {
				JCLoggerErr.Print("Can not move compressed files that have the same name")
				os.Exit(1)
			}

			err = jc.SetMoveTo(c, to)
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

	if *strptrCompressCMD == "bzip2" {
		c, err := jc.NewBZIP2Config()
		if err != nil {
			JCLoggerErr.Print(err)
			os.Exit(1)
		}

		if !c.SetCompLevel(*intptrCompressLevel) {
			JCLoggerWarn.Printf("Set compress level %d to %s failed",
				*intptrCompressLevel,
				c.Name())
		}

		err = jc.SetTimestampOption(c, *intptrTimestamp)
		if err != nil {
			JCLoggerErr.Print(err)
			os.Exit(1)
		}

		to, err := checkMoveTo(*strptrMoveTo)
		if err == nil {
			if haveSameName && (*intptrTimestamp < 3) {
				JCLoggerErr.Print("Can not move compressed files that have the same name")
				os.Exit(1)
			}

			err = jc.SetMoveTo(c, to)
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

	if *strptrCompressCMD == "xz" {
		c, err := jc.NewXZConfig()
		if err != nil {
			JCLoggerErr.Print(err)
			os.Exit(1)
		}

		if !c.SetCompLevel(*intptrCompressLevel) {
			JCLoggerWarn.Printf("Set compress level %d to %s failed",
				*intptrCompressLevel,
				c.Name())
		}

		err = jc.SetTimestampOption(c, *intptrTimestamp)
		if err != nil {
			JCLoggerErr.Print(err)
			os.Exit(1)
		}

		to, err := checkMoveTo(*strptrMoveTo)
		if err == nil {
			if haveSameName && (*intptrTimestamp < 3) {
				JCLoggerErr.Print("Can not move compressed files that have the same name")
				os.Exit(1)
			}

			err = jc.SetMoveTo(c, to)
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

		err = jc.SetTimestampOption(c, *intptrTimestamp)
		if err != nil {
			JCLoggerErr.Print(err)
			os.Exit(1)
		}

		to, err := checkMoveTo(*strptrMoveTo)
		if err == nil {
			if haveSameName && (*intptrTimestamp < 3) {
				JCLoggerErr.Print("Can not move compressed files that have the same name")
				os.Exit(1)
			}

			err = jc.SetMoveTo(c, to)
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

	if *strptrCompressCMD == "tgz" ||
		*strptrCompressCMD == "txz" ||
		*strptrCompressCMD == "tbz2" {
		c1, err := jc.NewTARConfig()
		if err != nil {
			JCLoggerErr.Print(err)
			os.Exit(1)
		}

		var c2 jc.Config
		if *strptrCompressCMD == "tgz" {
			c2, err = jc.NewGZIPConfig()
		} else if *strptrCompressCMD == "txz" {
			c2, err = jc.NewXZConfig()
		} else if *strptrCompressCMD == "tbz2" {
			c2, err = jc.NewBZIP2Config()
		}

		if !c2.SetCompLevel(*intptrCompressLevel) {
			JCLoggerWarn.Printf("Set compress level %d to %s failed",
				*intptrCompressLevel,
				c2.Name())
		}

		if err != nil {
			JCLoggerErr.Print(err)
			os.Exit(1)
		}

		err = jc.SetTimestampOption(c1, *intptrTimestamp)
		if err != nil {
			JCLoggerErr.Print(err)
			os.Exit(1)
		}

		to, err := checkMoveTo(*strptrMoveTo)
		if err == nil {
			if haveSameName && (*intptrTimestamp < 3) {
				JCLoggerErr.Print("Can not move compressed files that have the same name")
				os.Exit(1)
			}

			err = jc.SetMoveTo(c2, to)
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
