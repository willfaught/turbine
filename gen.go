package turbine

import (
	"bytes"
	"fmt"
	"go/format"
	"regexp"
	"text/template"
)

const identifierContent = "[[:alnum:]](?:_?[[:alnum:]])*?"

var (
	regexpEscape      = regexp.MustCompile("ESCAPE_([[:alnum:]]+)")
	regexpEscapeX     = regexp.MustCompile("ESCAPEX_(" + identifierContent + ")_ENDESCAPE")
	regexpEscaped     = regexp.MustCompile(`\{\{\/\*GOO:ESCAPE\*\/\}\}`)
	regexpEscapedX    = regexp.MustCompile(`\{\{\/\*GOO:ESCAPEX\*\/\}\}`)
	regexpFields      = regexp.MustCompile("FIELDS_(" + identifierContent + ")_ENDFIELDS")
	regexpGenerate    = regexp.MustCompile(`(?m:^\/\/go:generate.*$)`)
	regexpIdentifier  = regexp.MustCompile("__(" + identifierContent + ")__")
	regexpKeys        = regexp.MustCompile("KEYS_(" + identifierContent + ")_ENDKEYS")
	regexpLine        = regexp.MustCompile(`\/\/\/(.*)`)
	regexpLonghand    = regexp.MustCompile("(?m:^__X_(" + identifierContent + ")__$)")
	regexpMethods     = regexp.MustCompile("METHODS_(" + identifierContent + ")_ENDMETHODS")
	regexpMultiLine   = regexp.MustCompile(`(?s:\/\*\*(.*?)\*\*\/)`)
	regexpShorthand   = regexp.MustCompile("(?m:^__(" + identifierContent + ")__$)")
	regexpUnderscoreR = regexp.MustCompile("_")
	//regexpPackage     = regexp.MustCompile(`(?m:^)\w+(?m:$)`)
)

var (
	textDollarOne       = []byte("$1")
	textDollarOneBraces = []byte("{{$1}}")
	textDot             = []byte(".")
	textEmpty           = []byte("")
	textPlaceholder     = []byte("{{/*GOO:ESCAPE*/}}")
	textPlaceholderX    = []byte("{{/*GOO:ESCAPEX*/}}")
	textSpace           = []byte(" ")
	textUnderscoreB     = []byte("_")
	textUnderscores     = []byte("__")
)

var aliases = map[string][]byte{
	"ASSIGN":   []byte(":="),
	"DOT":      []byte("{{.}}"),
	"EMPTY":    []byte(""),
	"LBRACES":  []byte("{{"),
	"LCOMMENT": []byte("{{/*"),
	"RBRACES":  []byte("}}"),
	"RCOMMENT": []byte("*/}}"),
	"TRIM":     []byte("-"),
}

var funcs = map[*regexp.Regexp][]byte{
	regexp.MustCompile("ACTION_(.+?)_ENDACTION"):                       []byte("{{$1}}"),
	regexp.MustCompile("BLOCK_(.+?)_BEGIN_(.+?)_ENDBLOCK"):             []byte(`{{block "$1" $2}}`),
	regexp.MustCompile("FIELDX_(.+?)_ENDFIELD"):                        []byte(".$1"),
	regexp.MustCompile("FIELD_([[:alnum:]]+)"):                         []byte(".$1"),
	regexp.MustCompile("GROUPX_(.+?)_ENDGROUP"):                        []byte("($1)"),
	regexp.MustCompile("GROUP_([[:alnum:]]+)"):                         []byte("($1)"),
	regexp.MustCompile("IFE_(.+?)_THEN_(.+?)_ELSE_(.+?)_ENDIF"):        []byte("{{if $1}}$2{{else}}$3{{end}}"),
	regexp.MustCompile("IF_(.+?)_THEN_(.+?)_ENDIF"):                    []byte("{{if $1}}$2{{end}}"),
	regexp.MustCompile("INIT_([[:alnum:]]+)_(.+?)_ENDINIT"):            []byte("$$$1 := $2"),
	regexp.MustCompile("KEYX_(.+?)_ENDKEY"):                            []byte(".$1"),
	regexp.MustCompile("KEY_([[:alnum:]]+)"):                           []byte(".$1"),
	regexp.MustCompile("METHODX_(.+?)_ENDMETHOD"):                      []byte(".$1"),
	regexp.MustCompile("METHOD_([[:alnum:]]+)"):                        []byte(".$1"),
	regexp.MustCompile("OMITX_.+_ENDOMIT"):                             []byte(""),
	regexp.MustCompile("OMIT_[[:alnum:]]+"):                            []byte(""),
	regexp.MustCompile("RANGEE_(.+?)_BEGIN_(.+?)_ELSE_(.+?)_ENDRANGE"): []byte("{{range $1}}$2{{else}}$3{{end}}"),
	regexp.MustCompile("RANGE_(.+?)_BEGIN_(.+?)_ENDRANGE"):             []byte("{{range $1}}$2{{end}}"),
	regexp.MustCompile("RAWX_(.+?)_ENDRAW"):                            []byte("`$1`"),
	regexp.MustCompile("RAW_([[:alnum:]]+)"):                           []byte("`$1`"),
	regexp.MustCompile("RUNEO_([0-7]{3})"):                             []byte(`\$1`),
	regexp.MustCompile("RUNEU4_([0-9A-Fa-f]{4})"):                      []byte(`\u$1`),
	regexp.MustCompile("RUNEU8_([0-9A-Fa-f]{8})"):                      []byte(`\U$1`),
	regexp.MustCompile("RUNEX_([0-9A-Fa-f]{2})"):                       []byte(`\x$1`),
	regexp.MustCompile("RUNE_([[:alnum:]]+)"):                          []byte("'$1'"),
	regexp.MustCompile("STRINGX_(.+?)_ENDSTRING"):                      []byte(`"$1"`),
	regexp.MustCompile("STRING_([[:alnum:]]+)"):                        []byte(`"$1"`),
	regexp.MustCompile("TEMPLATE_(.+?)_BEGIN_(.+?)_ENDTEMPLATE"):       []byte(`{{template "$1" $2}}`),
	regexp.MustCompile("UNICODE_([[:alnum:]]+)"):                       []byte(`'\u$1'`),
	regexp.MustCompile("VAR_([[:alnum:]]+)"):                           []byte("$$$1"),
	regexp.MustCompile("WITHE_(.+?)_BEGIN_(.+?)_ELSE_(.+?)_ENDWITH"):   []byte("{{with $1}}$2{{else}}$3{{end}}"),
	regexp.MustCompile("WITH_(.+?)_BEGIN_(.+?)_ENDWITH"):               []byte("{{with $1}}$2{{end}}"),
}

var specials = map[*regexp.Regexp]func([]byte) []byte{
	regexpFields:  funcFields,
	regexpKeys:    funcKeys,
	regexpMethods: funcMethods,
}

var symbols = map[string][]byte{
	"AMPERSAND":   []byte("&"),
	"APOSTROPHE":  []byte("'"),
	"ASTERISK":    []byte("*"),
	"AT":          []byte("@"),
	"BACKSLASH":   []byte(`\`),
	"CARET":       []byte("^"),
	"COLON":       []byte(":"),
	"COMMA":       []byte(","),
	"DOLLAR":      []byte("$"),
	"EQUALS":      []byte("="),
	"EXCLAMATION": []byte("!"),
	"GRAVE":       []byte("`"),
	"GREATER":     []byte(">"),
	"LBRACE":      []byte("{"),
	"LBRACKET":    []byte("["),
	"LESS":        []byte("<"),
	"LPAREN":      []byte("("),
	"MINUS":       []byte("-"),
	"NUMBER":      []byte("#"),
	"PERCENT":     []byte(`%`),
	"PERIOD":      []byte("."),
	"PIPE":        []byte("|"),
	"PLUS":        []byte("+"),
	"QUESTION":    []byte("?"),
	"QUOTATION":   []byte(`"`),
	"RBRACE":      []byte("}"),
	"RBRACKET":    []byte("]"),
	"RPAREN":      []byte(")"),
	"SEMICOLON":   []byte(";"),
	"SLASH":       []byte("/"),
	"SPACE":       []byte(" "),
	"TILDE":       []byte("~"),
	"UNDERSCORE":  []byte("_"),
}

// Execute runs text as a text/template Template with name and data.
func Execute(name string, text []byte, data interface{}) ([]byte, error) {
	var t, err = template.New(name).Parse(string(text))

	if err != nil {
		var nameformat string

		if name != "" {
			nameformat = " " + name
		}

		return nil, &Error{error: fmt.Errorf("cannot parse%s: %v", nameformat, err), Data: data, Text: string(text)}
	}

	var b bytes.Buffer

	if err = t.Execute(&b, data); err != nil {
		var nameformat string

		if name != "" {
			nameformat = " " + name
		}

		return nil, &Error{error: fmt.Errorf("cannot execute%s: %v", nameformat, err), Data: data, Text: string(text)}
	}

	return b.Bytes(), nil
}

// Format formats text with go/format.
func Format(text []byte) ([]byte, error) {
	var b, err = format.Source(text)

	if err != nil {
		return nil, &Error{error: fmt.Errorf("cannot format: %v", err), Text: string(text)}
	}

	return b, nil
}

// Generate calls Process, Execute, and then Format.
func Generate(name string, text []byte, data interface{}) ([]byte, error) {
	var err error

	if text, err = Execute(name, Process(text), data); err != nil {
		return nil, err
	}

	if text, err = Format(text); err != nil {
		err.(*Error).Data = data

		return nil, err
	}

	return text, nil
}

// Process strips generate comments, renames identifiers of the form __X__ to
// {{.X}}, and replaes single-line comments beginning with a slash and
// multi-line comments beginning and ending with an asterisk replaced with their
// content.
func Process(text []byte) []byte {
	text = regexpGenerate.ReplaceAll(text, textEmpty)
	text = regexpIdentifier.ReplaceAllFunc(text, preprocess)
	text = regexpLine.ReplaceAll(text, textDollarOne)
	text = regexpMultiLine.ReplaceAll(text, textDollarOne)

	return text
}

func funcEscape(bs []byte) []byte {
	return textPlaceholder
}

func funcEscapeX(bs []byte) []byte {
	return textPlaceholderX
}

func funcFields(bs []byte) []byte {
	return append([]byte("."), regexpUnderscoreR.ReplaceAll(regexpFields.FindSubmatch(bs)[1], textDot)...)
}

func funcKeys(bs []byte) []byte {
	return append([]byte("."), regexpUnderscoreR.ReplaceAll(regexpKeys.FindSubmatch(bs)[1], textDot)...)
}

func funcMethods(bs []byte) []byte {
	return append([]byte("."), regexpUnderscoreR.ReplaceAll(regexpMethods.FindSubmatch(bs)[1], textDot)...)
}

func preprocess(bs []byte) []byte {
	var escapes = regexpEscape.FindAll(bs, -1)
	var escapesX = regexpEscapeX.FindAll(bs, -1)

	bs = regexpEscape.ReplaceAllFunc(bs, funcEscape)
	bs = regexpEscapeX.ReplaceAllFunc(bs, funcEscapeX)
	bs = regexpLonghand.ReplaceAll(bs, textDollarOne)
	bs = regexpShorthand.ReplaceAll(bs, textDollarOneBraces)

	for call, f := range specials {
		bs = call.ReplaceAllFunc(bs, f)
	}

	for call, result := range funcs {
		bs = call.ReplaceAll(bs, result)
	}

	for short, long := range aliases {
		bs = bytes.Replace(bs, []byte(short), long, -1)
	}

	bs = bytes.TrimPrefix(bs, textUnderscores)
	bs = bytes.TrimSuffix(bs, textUnderscores)
	bs = bytes.Replace(bs, textUnderscoreB, textSpace, -1)

	for name, symbol := range symbols {
		bs = bytes.Replace(bs, []byte(name), symbol, -1)
	}

	bs = regexpEscaped.ReplaceAllFunc(bs, func([]byte) []byte {
		var bs = escapes[0]

		escapes = escapes[1:]

		return regexpEscape.ReplaceAll(bs, textDollarOne)
	})

	bs = regexpEscapedX.ReplaceAllFunc(bs, func([]byte) []byte {
		var bs = escapesX[0]

		escapesX = escapesX[1:]

		return regexpEscapeX.ReplaceAll(bs, textDollarOne)
	})

	return bs
}

// Error has the corresponding template information.
type Error struct {
	error
	Data interface{} // The template data.
	Text string      // The template text.
}
