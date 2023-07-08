package turbine

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/fatih/camelcase"
)

// Ident is an identifier. TODO.
type Ident struct {
	Exported   string
	Initial    string
	Name       string
	Unexported string
	Words      []string
}

func newIdent(name string) *Ident {
	var id = Ident{Name: name}
	var words []string

	for _, d := range strings.Split(name, "_") {
		for _, w := range camelcase.Split(d) {
			words = append(words, w)
		}
	}

	var initials []string

	for i := range words {
		words[i] = strings.ToLower(words[i])

		var r, _ = utf8.DecodeRuneInString(words[i])

		if r == utf8.RuneError {
			panic(name)
		}

		initials = append(initials, strings.ToLower(fmt.Sprintf("%c", r)))
	}

	id.Words = words

	var r, n = utf8.DecodeRune([]byte(name))

	if unicode.IsUpper(r) {
		id.Exported = name
		id.Initial = fmt.Sprintf("%c", unicode.ToLower(r))
		id.Unexported = fmt.Sprintf("%c%v", unicode.ToLower(r), name[n:]) // TODO: Lowercase first word
	} else {
		id.Exported = fmt.Sprintf("%c%v", unicode.ToUpper(r), name[n:]) // TODO: Lowercase first word
		id.Initial = fmt.Sprintf("%c", r)
		id.Unexported = name
	}

	return &id
}

// TODO
func (i *Ident) String() string {
	return i.Name
}
