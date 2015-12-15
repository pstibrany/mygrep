package main;

func NewContextOutput(contextLines, contextBefore, contextAfter int, output Output, co Output) Output {
    if contextBefore < 0 {
        contextBefore = 0
    }
    
    if contextAfter < 0 {
        contextAfter = 0
    }
    
    if contextLines > 0 {
        if contextBefore == 0 {
            contextBefore = contextLines
        }
        
        if contextAfter == 0 {
            contextAfter = contextLines
        }
    }
    
    if contextBefore == 0 && contextAfter == 0 {
        return output
    }
    
    return &contextOutput{contextBefore: contextBefore, contextAfter: contextAfter, output: output, contextOutput: co}
}

/* Output with context */
type contextOutput struct {
    contextBefore int
    contextAfter int

    lines [][]byte
    linesToPrint int
    output Output
    contextOutput Output // used to output the context lines
}

func (this *contextOutput) Reset(filename string) {
    this.lines = make([][]byte, this.contextBefore)
    this.linesToPrint = 0
    this.output.Reset(filename)
}

func (this *contextOutput) ProcessLine(lineNumber int, line []byte, matches [][]int) error {
    if len(matches) > 0 && this.contextBefore > 0 {
        for i := 0; i < this.contextBefore; i++ {
            if this.lines[i] != nil {
                if err := this.contextOutput.ProcessLine(lineNumber-this.contextBefore+i, this.lines[i], nil); err != nil {
                    return err
                }
            }
        }
        
        this.linesToPrint = this.contextAfter
    }

    this.addLineToContext(line)

    if err := this.output.ProcessLine(lineNumber, line, matches); err != nil {
        return err;
    }
    
    if len(matches) == 0 && this.linesToPrint > 0 {
        this.linesToPrint --

        return this.contextOutput.ProcessLine(lineNumber, line, matches)        
    }
    
    return nil
}

func (this *contextOutput) addLineToContext(line []byte) {
    if this.contextBefore > 0 {
        for ix := 1; ix < this.contextBefore; ix ++ {
            this.lines[ix - 1] = this.lines[ix]
        }
        
        last := this.contextBefore - 1;
        this.lines[last] = nil
        this.lines[last] = append(this.lines[last], line...)
    }
}
