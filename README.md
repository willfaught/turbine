# turbine

[![Go Reference](https://pkg.go.dev/badge/go.dev/pkg/github.com/willfaught/turbine.svg)](https://pkg.go.dev/go.dev/pkg/github.com/willfaught/turbine)

Turbine generates Go files from other Go files. Go file templates are transformed like this:

1. Line comments matching the regular expression `^\/\/go:generate.*$` are stripped.
1. Text matching the regular expression `__[[:alnum:]](?:_?[[:alnum:]])*?__` is transformed like this:
    1. Text matching the regular expressions `ESCAPE_([[:alnum:]]+)` and `ESCAPEX_([[:alnum:]](?:_?[[:alnum:]])*?)_ENDESCAPE` is replaced with the regular expression result `$1` without any other steps in this sublist affecting it.
    1. Text matching the regular expression `__([[:alnum:]](?:_?[[:alnum:]])*?)__` is replaced with the regular exxpression result `{{$1}}`.
    1. Text matching the regular expression `__X_([[:alnum:]](?:_?[[:alnum:]])*?)__` is replaced with the regular exxpression result `$1`.
    1. Text of the form `FIELDS_f_ENDFIELDS` is replaced with `{{.f}}`, `FIELDS_f_g_ENDFIELDS` is replaced with `{{.f.g}}`, and so on. `FIELDS` can also be `KEYS` and `METHODS`.
    1. Various functions and aliases are interpreted. These are undocumented at this time. See `funcs` and `aliases` in the code for more information.
    1. The `__` prefix and suffix are stripped.
    1. All underscores are replaced with spaces.
    1. Various symbols are interpreted. These are undocumented at this time. See `symbols` in the code for more information.
1. Line comments matching the regular expression `\/\/\/.*$` are stripped.
1. General comments beginning and ending with an asterisk (`/**` and `**/`) are stripped.
1. The code is executed as a Go text/template.Template with a supplied data context.
1. The code is formatted.

A "make" subcommand provides a general-purpose method of generating Go files from any Go file templates and data.

Several other subcommands provide useful Go file templates.

A common scenario is to generate a Go file based on information about a declaration in another Go file. The turbine package has wys to facilitate getting that information and using it as a template context.

*Note: This document is woefully incomplete.*

## Usage

```output
$ turbine -h
Usage of turbine:
  make
    	Generate a file from a template and data.
  mock
    	Generate an interface mock.
  pack
    	Generate a constant declaration that contains a file.
  stub
    	Generate an interface stub.
  wrap
    	Generate an interface wrapper.
```

```output
$ turbine make -h
Usage of turbine make:
  -data string
    	The template data as JSON.
  -identifier string
    	The declaration identifier.
  -input string
    	The input file. Defaults to standard input.
  -output string
    	The output file. Defaults to standard output.
  -package string
    	The declaration package.
```

```output
$ turbine mock -h
Usage of turbine mock:
  -identifier string
    	The declaration identifier.
  -output string
    	The output file. Defaults to standard output.
  -package string
    	The declaration package.
```

```output
$ turbine pack -h
Usage of turbine pack:
  -identifier string
    	The declaration identifier.
  -input string
    	The input file. Defaults to standard input.
  -output string
    	The output file. Defaults to standard output.
  -package string
    	The declaration package.
```

```output
$ turbine stub -h
Usage of turbine stub:
  -identifier string
    	The declaration identifier.
  -output string
    	The output file. Defaults to standard output.
  -package string
    	The declaration package.
```

```output
$ turbine wrap -h
Usage of turbine wrap:
  -identifier string
    	The declaration identifier.
  -output string
    	The output file. Defaults to standard output.
  -package string
    	The declaration package.
```
