package main

import "os"
import "flag"
import "fmt"
import "regexp"
import "io"
import "bufio"
import "path/filepath"
import "compress/gzip"
import "compress/bzip2"
import "runtime/pprof"

func main() {
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	memprofile := flag.String("memprofile", "", "write memory profile to this file")

	replaceString := flag.String("r", "", "Replacement string. Use $N, ${N} or $name for referencing groups")
	invertMatch := flag.Bool("v", false, "Invert match")
	contextBefore := flag.Int("B", 0, "Print this number of lines before the match.")
	contextAfter := flag.Int("A", 0, "Print this number of lines after the match.")
	contextLines := flag.Int("C", 0, "Print this number of lines before and after the match.")
	printLineNumbers := flag.Bool("l", false, "Print line numbers")
	printVersion := flag.Bool("version", false, "Print version")
	matchesOnly := flag.Bool("m", false, "Print only matched parts of line. Useful when using replacement string.")
	disableHighlights := flag.Bool("nh", false, "Disable highlights. Highlights are enabled by default")
	disableFilenames := flag.Bool("nf", false, "Disable filename prefixes.")
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			fmt.Fprintln(os.Stderr)
		} else {
			pprof.StartCPUProfile(f)
			defer pprof.StopCPUProfile()
		}
	}

	if *printVersion {
		fmt.Fprintln(os.Stderr, BuildId)
		return
	}

	args := flag.Args()
	if len(args) < 1 {
		usage()
		return
	}

	re, err := regexp.Compile(flag.Args()[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to compile regular expression: %s\n", err)
		os.Exit(1)
	}

	allFiles := flag.Args()[1:]

	var output Output = NewStandardOutput(!*disableHighlights)
	output = NewPrefixingOutput(len(allFiles) > 1 && !*disableFilenames, *printLineNumbers, output)

	// origOutput is simply output that prints any line (possibly with prefix and linenumber)
	origOutput := output

	if *invertMatch {
		output = NewMatchingOutput(false, true, output)
	} else {
		output = NewMatchingOutput(true, false, output)
	}

	output = NewContextOutput(*contextLines, *contextBefore, *contextAfter, output, origOutput)

	if *matchesOnly {
		output = NewMatchesOnlyOutput(output)
	}
	output = NewReplaceOutput(re, *replaceString, output)

	if len(allFiles) == 0 {
		err := grepReader(re, "", os.Stdin, output)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	} else {
		for _, file := range allFiles {
			err := grepFile(re, file, output)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		} else {
			pprof.WriteHeapProfile(f)
			f.Close()
		}
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s <regex> [files]\n", os.Args[0])
	flag.PrintDefaults()
}

func grepFile(re *regexp.Regexp, filename string, output Output) (writeError error) {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil
	}

	defer f.Close()

	var reader io.Reader = f
	reader, err = prepareReader(reader, filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open %s: %s\n", filename, err)
		return nil
	}
	return grepReader(re, filename, reader, output)
}

func prepareReader(reader io.Reader, filename string) (io.Reader, error) {
	ext := filepath.Ext(filename)
	switch {
	case ext == ".gz":
		return gzip.NewReader(reader)
	case ext == ".bz2":
		return bzip2.NewReader(reader), nil
	default:
		return reader, nil
	}
}

func grepReader(re *regexp.Regexp, filename string, reader io.Reader, output Output) error {
	output.Reset(filename)

	bufReader := bufio.NewReaderSize(reader, 65535)
	var err error = nil
	var line []byte = nil
	var lineNumber int = 0

	for err == nil {
		line, err = bufReader.ReadSlice('\n')
		var toProcess []byte = line

		if err == bufio.ErrBufferFull {
			var buffer []byte = nil
			for ; err == bufio.ErrBufferFull; line, err = bufReader.ReadSlice('\n') {
				buffer = append(buffer, line...)
			}
			buffer = append(buffer, line...)
			toProcess = buffer
		}

		lineNumber++

		matches := re.FindAllSubmatchIndex(toProcess, -1)
		err2 := output.ProcessLine(lineNumber, toProcess, matches)

		if err2 != nil {
			return err2
		}
	}

	if err != io.EOF {
		fmt.Fprintln(os.Stderr, err)
	}

	return nil
}
