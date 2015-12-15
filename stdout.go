package main;

import "fmt"
import "os"

/* This output prints all lines. Lines with matches have matches highlighted. */
type stdout struct {
    includeFilename bool
    includeLineNumber bool
    highlight bool
    
    filename string
    lastPrintedLine int
}

func NewStandardOutput(highlight bool) Output {
    return &stdout{highlight: highlight}
}

func (this *stdout) ProcessLine(lineNumber int, line []byte, matches [][]int) error {
    if lineNumber <= this.lastPrintedLine {
        return nil
    }

    this.lastPrintedLine = lineNumber
    
    if !this.highlight {
        _, err := os.Stdout.Write(line)
        return err
    }
    
    start := 0
    for _, match := range matches {
        b := match[0]
        e := match[1]
        
        if _, err := os.Stdout.WriteString(fmt.Sprintf("%s\033[31m%s\033[0m", line[start:b], line[b:e])); err != nil {
            return err
        }
        
        start = e
    }
    
    _, err := os.Stdout.Write(line[start:])
    
    return err
}

func (this *stdout) Reset(filename string) {
    this.filename = filename
    this.lastPrintedLine = 0
}