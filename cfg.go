package cfg

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
)

func StringSlice(str string, split ...string) []string {
	sp := ","
	if len(split) > 0 {
		sp = split[0]
	}
	out := strings.Split(str, sp)
	for k, v := range out {
		out[k] = strings.TrimSpace(v)
	}
	return out
}

func IntSlice(str string, split ...string) ([]int, error) {
	var err error
	sp := ","
	if len(split) > 0 {
		sp = split[0]
	}
	temp := strings.Split(str, sp)
	out := make([]int, len(temp))
	for k, v := range temp {
		v = strings.TrimSpace(v)
		out[k], err = strconv.Atoi(v)
		if err != nil {
			break
		}
	}

	return out, err
}

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
	convnl(buf, buf)
	buf.WriteRune('\n')
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
			a := []rune(l)
			if a[0] == '#' && !insideVar {
				continue
			}
			for i := 0; i < len(a); i++ {
				if !insideVar {
					switch a[i] {
					case '=':
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
					if i == len(a)-1 {
						// last char is always a \n
						insideVar = false
						if i > 0 {
							if a[i-1] == '\\' {
								varVal.WriteRune('\n')
								insideVar = true
							}
						}
						continue
					}
					//TODO: parse "\" better
					if a[i] == '\\' {
						if i > 0 {
							if a[i-1] == '\\' {
								varVal.WriteRune(a[i])
							}
						}
					} else {
						varVal.WriteRune(a[i])
					}
				}
			}
			lc++
			if !insideVar {
				v00 := varName.String()
				v01 := varVal.String()
				//fmt.Println("v00", v00, "v01", v01)
				mp[strings.TrimSpace(v00)] = strings.Trim(v01, " ")
				varName.Truncate(0)
				varVal.Truncate(0)
			}
		}

	}
	return mp, nil
}

func convnl(input io.Reader, output io.Writer) {
	r := bufio.NewReader(input)
	w := bufio.NewWriter(output)
	var prev rune
	for {
		cur, _, err := r.ReadRune()
		if err != nil {
			if prev == '\r' {
				w.WriteRune(prev)
			}
			break
		}
		if cur == '\r' {
			if prev == '\r' {
				w.WriteRune(prev)
			}
			prev = cur
			continue
		}
		if cur != '\n' && prev == '\r' {
			w.WriteRune(prev)
		}
		prev = cur
		w.WriteRune(cur)
	}
	w.Flush()
}
