package turbine

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/fatih/camelcase"
)

type Name string

func (n Name) Initial() Name {
	r, _ := utf8.DecodeRune([]byte(n))
	if r == utf8.RuneError {
		panic(n)
	}
	if unicode.IsUpper(r) {
		r = unicode.ToLower(r)
	}
	return Name(fmt.Sprintf("%c", r))
}

func (n Name) Initials() Name {
	var rs []rune
	for _, d := range strings.Split(string(n), "_") {
		for _, w := range camelcase.Split(d) {
			r, _ := utf8.DecodeRuneInString(w)
			if r == utf8.RuneError {
				panic(w)
			}
			rs = append(rs, unicode.ToLower(r))
		}
	}
	return Name(string(rs))
}

func (n Name) IsExported() bool {
	r, _ := utf8.DecodeRuneInString(string(n))
	return unicode.IsUpper(r)
}

func (n Name) Export() Name {
	r, x := utf8.DecodeRuneInString(string(n))
	if r == utf8.RuneError {
		panic(n)
	}
	return Name(fmt.Sprintf("%c%v", unicode.ToUpper(r), n[x:]))
}

func (n Name) Unexport() Name {
	r, x := utf8.DecodeRuneInString(string(n))
	if r == utf8.RuneError {
		panic(n)
	}
	return Name(fmt.Sprintf("%c%v", unicode.ToLower(r), n[x:]))
}
