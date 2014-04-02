package cfg

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
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
	_, err = io.Copy(&buffer, fl)
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
	for {
		l, e := b.ReadString('\n')
		if e != nil {
			break
		} else {
			fmt.Println("!RIOT! ", l)
			a := []rune(l)
			if a[0] == '#' && !insideVar {
				continue
			}
			for i := 0; i < len(a); i++ {
				fmt.Printf(" %s (%v) ", string(a[i]), lc)
				if !insideVar {
					fmt.Print(insideVar, " ")
					switch a[i] {
					case '=':
						fmt.Print(" is '=' ")
						if varName.Len() < 1 {
							return nil, errors.New("Syntax error at line " + strconv.Itoa(lc) + " (" + strconv.Itoa(i) + "): " + l)
						}
						insideVar = true
					case '\\', '\'', '"', '#', '^', '&':
						return nil, errors.New("Syntax error at line " + strconv.Itoa(lc) + " (" + strconv.Itoa(i) + "): " + l)
					default:
						varName.WriteRune(a[i])
					}
				} else {
					fmt.Print(insideVar)
					if i == len(a)-1 && a[i] == '\\' {
						varVal.WriteRune('\n')
						fmt.Print("\n")
						continue
					} else if i == len(a)-1 {
						fmt.Print("inside var will be cause because EOL\n")
						insideVar = false
						break
					}
					//TODO: parse "\" better
					varVal.WriteRune(a[i])
				}
			}
			lc++
			if !insideVar {
				v00 := varName.String()
				v01 := varVal.String()
				fmt.Println("v00", v00, "v01", v01)
				mp[strings.TrimSpace(v00)] = strings.TrimSpace(v01)
				varName.Truncate(0)
				varVal.Truncate(0)
			}
		}

	}
	return mp, nil
}
