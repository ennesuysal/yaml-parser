package main

import "regexp"

func rgxShortcut(rgx string, txt string) ([][][]byte, error) {
	r, err := regexp.Compile(rgx)

	if err != nil {
		return nil, err
	}

	if !r.Match([]byte(txt)) {
		return nil, err
	}
	match := r.FindAllSubmatch([]byte(txt), -1)
	return match, nil
}

func trim(line string) (int, string) {
	i := 0
	for ; i < len(line); i++ {
		if line[i] != ' ' && line[i] != '\t' {
			break
		}
	}

	indent := i * 2

	if analyze(line) == (arrayElement{}) {
		indent += 1
	}
	if analyze(line) == (continuingArr{}) {
		indent += 1
	}

	return indent, line[i:]
}
