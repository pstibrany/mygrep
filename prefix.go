package main;

import "fmt"
import "os"

func NewPrefixingOutput(includeFilename, includeLineNumber bool, output Output) Output {
    result := output

    if includeLineNumber {
        result = &addLineNumberPrefix{output: result}
    }
    
    if includeFilename {
        result = &addFilenamePrefix{output: result}
    }
    
    return result
}

type addLineNumberPrefix struct {
    output Output
}

func (this *addLineNumberPrefix) Reset(filename string) {
    this.output.Reset(filename)
}

func (this *addLineNumberPrefix) ProcessLine(lineNumber int, line []byte, matches [][]int) error {
    _, err := os.Stdout.WriteString(fmt.Sprintf("%d:", lineNumber))
    if err != nil {
        return err
    }
    return this.output.ProcessLine(lineNumber, line, matches)
}

type addFilenamePrefix struct {
    prefix string
    output Output
}

func (this *addFilenamePrefix) Reset(filename string) {
    this.prefix = filename + ":"
    this.output.Reset(filename)
}

func (this *addFilenamePrefix) ProcessLine(lineNumber int, line []byte, matches [][]int) error {
    _, err := os.Stdout.WriteString(this.prefix)
    if err != nil {
        return err
    }
    return this.output.ProcessLine(lineNumber, line, matches)
}
