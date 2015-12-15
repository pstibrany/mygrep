package main;

type Output interface {
    /* matches may be nil, if line was not matched */
    ProcessLine(lineNumber int, line []byte, matches [][]int) error
    Reset(filename string)
}

type matchingOutput struct {
    printMatchedLines bool
    printUnmatchedLines bool
    
    output Output
}

func (this *matchingOutput) ProcessLine(lineNumber int, line []byte, matches [][]int) error {
    if (matches != nil && this.printMatchedLines) {
        return this.output.ProcessLine(lineNumber, line, matches)
    }
    
    if (matches == nil && this.printUnmatchedLines) {
        return this.output.ProcessLine(lineNumber, line, matches)
    }
    
    return nil
}

func (this *matchingOutput) Reset(filename string) {
    this.output.Reset(filename)
}

func NewMatchingOutput(printMatchedLines, printUnmatchedLines bool, output Output) Output {
    return &matchingOutput{printMatchedLines: printMatchedLines, printUnmatchedLines: printUnmatchedLines, output: output}
}