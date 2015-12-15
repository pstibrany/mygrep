package main;

import "bytes"

type MatchesOnlyOutput struct {
    output Output
}

func NewMatchesOnlyOutput(matchesOnly bool, output Output) Output {
    if !matchesOnly {
        return output
    }
    
    return &MatchesOnlyOutput{output: output}
}

func (this *MatchesOnlyOutput) Reset(filename string) {
    this.output.Reset(filename)
}

func (this *MatchesOnlyOutput) ProcessLine(lineNumber int, line []byte, matches [][]int) error {
    newLine := &bytes.Buffer{}
    
    var newMatches [][]int
    if matches != nil {
        newMatches = make([][]int, len(matches))
    }
    
    for ix, match := range matches {
        b := match[0]
        e := match[1]
        
        newMatch := make([]int, len(match))
        for i := 0; i < len(match); i ++ {
            newMatch[i] = match[i] - match[0] + newLine.Len()
        }
        
        newLine.Write(line[b:e])
        
        newMatches[ix] = newMatch
    }
    
    return this.output.ProcessLine(lineNumber, newLine.Bytes(), newMatches)
}
