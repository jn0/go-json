package json

import (
	"fmt"
	"strconv"
	"strings"
)

type SyntaxError error
type NoValue SyntaxError
type BadValue SyntaxError
type BadTail SyntaxError
type MissedValue SyntaxError

var digits = map[byte]bool{
	'0': true,
	'1': true,
	'2': true,
	'3': true,
	'4': true,
	'5': true,
	'6': true,
	'7': true,
	'8': true,
	'9': true,
}

func isDigit(c byte) bool {
	res, ok := digits[c]
	return ok && res
}

func parseObject(s string) (v JsonValue, t string, e error) {
	s = strings.TrimSpace(s)
	if s == "" {
		e = NoValue(fmt.Errorf("No value for object"))
		return
	}
	if s[0] != '{' {
		t = s
		e = SyntaxError(fmt.Errorf("Not an object %+q", s))
		return
	}
	v = new(JsonObject)
	t = strings.TrimSpace(s[1:])
	ok := false
	for t != "" {
		if t[0] != '"' {
			e = SyntaxError(fmt.Errorf("%+q bad name at '%c'", t, t[0]))
			return
		}
		xi, name, sok := getString(t)
		if !sok {
			e = SyntaxError(fmt.Errorf("Bad name %+q at %+q", name, t))
			return
		}

		t = strings.TrimSpace(t[xi+2:])
		if t == "" || t[0] != ':' {
			e = SyntaxError(fmt.Errorf("%+q no colon after name %q", t, name))
			return
		}

		t = strings.TrimSpace(t[1:])
		if t == "" {
			e = SyntaxError(fmt.Errorf("%+q no value for name %q", t, name))
			return
		}
		xv, xt, xe := ParseValue(t)
		if xe != nil {
			e = xe
			t = xt
			return
		}
		v.Insert(name, xv)

		t = strings.TrimSpace(xt)
		if t == "" {
			break
		}
		if t[0] == '}' {
			t = strings.TrimSpace(t[1:])
			ok = true
			break
		}
		if t[0] == ',' {
			t = strings.TrimSpace(t[1:])
			if t == "" {
				e = MissedValue(fmt.Errorf("after comma"))
				return
			}
			continue
		}
		e = BadTail(fmt.Errorf("%+q bad tail", t))
		return
	}
	if !ok {
		e = SyntaxError(fmt.Errorf("No closing brace in object %+q", s))
	}
	return
}
func parseArray(s string) (v JsonValue, t string, e error) {
	s = strings.TrimSpace(s)
	if s == "" {
		e = NoValue(fmt.Errorf("No value for array"))
		return
	}
	if s[0] != '[' {
		t = s
		e = SyntaxError(fmt.Errorf("Not an array %+q", s))
		return
	}
	v = new(JsonArray)
	t = strings.TrimSpace(s[1:])
	ok := false
	for t != "" {
		xv, xt, xe := ParseValue(t)
		if xe != nil {
			e = xe
			t = xt
			return
		}
		v.Append(xv)

		t = strings.TrimSpace(xt)
		if t == "" {
			break
		}
		if t[0] == ']' {
			t = strings.TrimSpace(t[1:])
			ok = true
			break
		}
		if t[0] == ',' {
			t = strings.TrimSpace(t[1:])
			if t == "" {
				e = MissedValue(fmt.Errorf("after comma"))
				return
			}
			continue
		}
		e = BadTail(fmt.Errorf("%+q bad tail", t))
		return
	}
	if !ok {
		e = SyntaxError(fmt.Errorf("No closing bracket in array %+q", s))
	}
	return
}
func getString(s string) (pos int, res string, ok bool) {
	escape, unicode, res, hex := false, false, "", ""
	var c rune
	for pos, c = range s[1:] {
		if unicode {
			hex += string(c)
			if len(hex) == 4 {
				v, e := strconv.ParseInt(hex, 16, 64)
				if e != nil {
					panic(e)
				}
				res += string(rune(v))
				hex = ""
				unicode = false
			}
			continue
		}
		if escape {
			switch c {
			case 'b':
				res += "\b"
			case 'f':
				res += "\f"
			case 'n':
				res += "\n"
			case 'r':
				res += "\r"
			case 't':
				res += "\t"
			case 'u':
				unicode = true
			default:
				res += string(c) // slash, backslash, quote - relaxed...
			}
			escape = false
			continue
		}
		if c == '\\' {
			escape = true
			continue
		}
		if c == '"' {
			ok = true
			break
		}
		res += string(c)
	}
	return
}
func parseString(s string) (v JsonValue, t string, e error) {
	s = strings.TrimSpace(s)
	if s == "" {
		e = NoValue(fmt.Errorf("No value for string"))
		return
	}
	if s[0] != '"' {
		e = SyntaxError(fmt.Errorf("Not a string %+q", s))
		return
	}
	i, r, ok := getString(s)
	if !ok {
		e = SyntaxError(fmt.Errorf("Bad string %+q", s))
		return
	}
	t = s[i+2:]
	v = new(JsonString)
	v.Set(r)
	return
}
func parseNumber(s string) (v JsonValue, t string, e error) {
	s = strings.TrimSpace(s)
	if s == "" {
		e = NoValue(fmt.Errorf("No value for number"))
		return
	}
	isFloat := false
	intPart := ""
	frac := ""
	t = s
	if t[0] == '+' || t[0] == '-' {
		intPart = string(t[0])
		t = t[1:]
	}
	for t != "" && isDigit(t[0]) {
		intPart += string(t[0])
		t = t[1:]
	}
	if t != "" && t[0] == '.' {
		isFloat = true
		t = t[1:]
	}
	for isFloat && t != "" && isDigit(t[0]) {
		frac += string(t[0])
		t = t[1:]
	}
	if isFloat {
		v = new(JsonFloat)
		e = v.Parse(intPart + "." + frac)
	} else {
		v = new(JsonInt)
		e = v.Parse(intPart)
	}
	return
}
func parseBool(s string) (v JsonValue, t string, e error) {
	s = strings.TrimSpace(s)
	if s == "" {
		e = NoValue(fmt.Errorf("No value for bool"))
		return
	}
	if strings.HasPrefix(s, "true") {
		v = new(JsonBool)
		v.Set(true)
		t = strings.TrimSpace(s[4:])
		return
	}
	if strings.HasPrefix(s, "false") {
		v = new(JsonBool)
		v.Set(false)
		t = strings.TrimSpace(s[5:])
		return
	}
	return nil, s, BadValue(fmt.Errorf("%+q is neither 'true' nor 'false'", s))
}
func parseNull(s string) (v JsonValue, t string, e error) {
	s = strings.TrimSpace(s)
	if s == "" {
		e = NoValue(fmt.Errorf("No value for null"))
		return
	}
	if strings.HasPrefix(s, "null") {
		// v = new(JsonObject)
		t = strings.TrimSpace(s[4:])
		return
	}
	return nil, s, BadValue(fmt.Errorf("%+q is not 'null'", s))
}

func ParseValue(s string) (v JsonValue, t string, e error) {
	s = strings.TrimSpace(s)
	if s == "" {
		e = NoValue(fmt.Errorf("No value at all"))
		return
	}
	switch s[0] {
	case '{':
		v, t, e = parseObject(s)
	case '[':
		v, t, e = parseArray(s)
	case '"':
		v, t, e = parseString(s)
	case '-', '+', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		v, t, e = parseNumber(s)
	case 't', 'f':
		v, t, e = parseBool(s)
	case 'n':
		v, t, e = parseNull(s)
	default:
		e = BadValue(fmt.Errorf("%+q is not a value for '%c'=%#v", s, s[0], s[0]))
		return
	}
	return
}
