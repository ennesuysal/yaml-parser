package main

import "regexp"

const (
	singleLineRgx       = `([^:]+)\s*:\s*(.+)\s*$`
	continuingLineRgx   = `([^\s]+)\s*:\s*$`
	arrayElementRgx     = `^\s*((?:-\s*)+)\s*([^\s-:]+)\s*:?\s*([^\s]*)\s*`
	continuingStringRgx = `([^\s]+)\s*:\s*[>\|]\s*`
	continuingArrRgx    = `^\s*-\s*$`
)

func rgxShortcut(rgx string, txt string) ([][]string, error) {
	r, err := regexp.Compile(rgx)

	if err != nil {
		return nil, err
	}

	if !r.Match([]byte(txt)) {
		return nil, err
	}
	match := r.FindAllStringSubmatch(txt, -1)
	return match, nil
}

func trim(line string) (float32, string) {
	i := 0
	for ; i < len(line); i++ {
		if line[i] != ' ' && line[i] != '\t' {
			break
		}
	}

	return float32(i * 2), line[i:]
}
