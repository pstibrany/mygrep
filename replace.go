package main;

import "bytes"
import "regexp"

type ReplaceOutput struct {
    re *regexp.Regexp
    template []byte
    output Output
}

func NewReplaceOutput(re *regexp.Regexp, replaceString string, output Output) Output {
    if replaceString == "" {
        return output
    }
    
    return &ReplaceOutput{re: re, template: []byte(replaceString), output: output}
}

func (this *ReplaceOutput) Reset(filename string) {
    this.output.Reset(filename)
}

func (this *ReplaceOutput) ProcessLine(lineNumber int, line []byte, matches [][]int) error {
    newLine := &bytes.Buffer{}
    
    var newMatches [][]int
    if matches != nil {
        newMatches = make([][]int, len(matches))
    }
    
    start := 0
    
    for ix, match := range matches {
        b := match[0]
        e := match[1]
        
        newLine.Write(line[start:b])
        
        replaced := replaceMatch(this.re, this.template, line, match)

        newMatch := make([]int, 2)
        newMatch[0] = newLine.Len()
        newLine.Write(replaced)
        newMatch[1] = newLine.Len()
        
        newMatches[ix] = newMatch
        
        start = e
    }
    
    if start < len(line) {
        newLine.Write(line[start:])
    }
    
    return this.output.ProcessLine(lineNumber, newLine.Bytes(), newMatches)
}

func replaceMatch(re *regexp.Regexp, templateString, line []byte, matches []int) []byte {
    matchedLine := line[matches[0]:matches[1]]
    
    aMatches := make([]int, len(matches))
    for i := 0; i < len(matches); i ++ {
        aMatches[i] = matches[i] - matches[0]
    }
    
    return re.Expand(nil, templateString, matchedLine, aMatches)
}
