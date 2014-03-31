package cfg

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
)

func ParseFile(path string) (map[string]string, error) {
	var buffer bytes.Buffer
	if _, err := os.Stat(path); err != nil {
		return nil, err
	}
	fl, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fl.Close()
	_, err = io.Copy(&buffer, src)
	if err != nil {
		return nil, err
	}
	m, err := parse(&buffer)
	if err != nil {
		return nil, errors.New(path + " :: " + err.Error())
	}
	return m, nil
}

func ParseString(val string) (map[string]string, error) {
	var buffer bytes.Buffer
	_, err := buffer.WriteString(val)
	if err != nil {
		return nil, err
	}
	return parse(&buffer)
}

func parse(buf *bytes.Buffer) (map[string]string, error) {
	insideVar := false
	var varName bytes.Buffer
	var varVal bytes.Buffer
	b := bufio.NewReader(buf)
	lc := 0
	mp := make(map[string]string, 0)
	for l, e := b.ReadString('\n'); e == nil; {
		a := []rune(s)
		if a[0] == '#' && !insideVar {
			continue
		}
		for i := 0; i < len(a); i++ {
			if !insideVar {
				switch a[i] {
				case '=':
					if varName.Len() < 1 {
						return nil, errors.New("Syntax error at line " + strconv.Itoa() + " (" + strconv.Itoa(lc) + "): " + l)
					}
					insideVar = true
				case '\\', '\'', '"', '#', '^', '&':
					return nil, errors.New("Syntax error at line " + strconv.Itoa() + " (" + strconv.Itoa(lc) + "): " + l)
				default:
					varName.WriteRune(a[i])
				}
			} else {
				if i == len(a)-1 && a[i] == '\\' {
					varVal.WriteRune('\n')
					continue
				} else if i == len(a) {
					insideVar = false
				}
				//TODO: parse "\" better
				varVal.WriteRune(a[i])
			}
		}
		lc++
		if !insideVar {
			mp[varName.String()] = varVal.String()
		}
	}
	return mp, nil
}
