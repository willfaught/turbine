package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/build"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/willfaught/turbine/internal/turbine"
)

var suffixes = []string{
	"_386",
	"_amd64",
	"_amd64p32",
	"_android",
	"_arm",
	"_arm64",
	"_arm64be",
	"_armbe",
	"_darwin",
	"_dragonfly",
	"_freebsd",
	"_linux",
	"_mips",
	"_mips64",
	"_mips64le",
	"_mips64p32",
	"_mips64p32le",
	"_mipsle",
	"_nacl",
	"_netbsd",
	"_openbsd",
	"_plan9",
	"_ppc",
	"_ppc64",
	"_ppc64le",
	"_s390",
	"_s390x",
	"_solaris",
	"_sparc",
	"_sparc64",
	"_test",
	"_windows",
	"_zos",
}

func cmdDecl(args []string, name, template string) error {
	var f = flag.NewFlagSet("turbine "+name, flag.ExitOnError)

	var (
		flagidentifier = f.String("identifier", "", "The declaration identifier.")
		flagoutput     = f.String("output", "", "The output file. Defaults to standard output.")
		flagpackage    = f.String("package", "", "The declaration package.")
	)

	if err := f.Parse(args); err != nil {
		return fmt.Errorf("cannot parse program arguments: %v", err)
	}

	var d, err = turbine.NewDecl(*flagpackage, *flagidentifier)

	if err != nil {
		return fmt.Errorf("cannot parse declaration: %v", err)
	}

	b, err := turbine.Generate(name, []byte(template), d)

	if err != nil {
		return fmt.Errorf("cannot generate %s: %v", name, err)
	}

	return write(*flagoutput, b)
}

func cmdMain(args []string) error {
	var f = flag.NewFlagSet("turbine", flag.ExitOnError)

	f.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of turbine:\n")
		f.PrintDefaults()

		for _, c := range []struct {
			name, usage string
		}{
			{"make", "Generate a file from a template and data."},
			{"mock", "Generate an interface mock."},
			{"pack", "Generate a constant declaration that contains a file."},
			{"stub", "Generate an interface stub."},
			{"wrap", "Generate an interface wrapper."},
		} {
			var indent string

			if len(c.name) > 2 {
				indent = "\n    "
			}

			fmt.Fprintf(os.Stderr, "  %s%s\t%s\n", c.name, indent, c.usage)
		}
	}

	if err := f.Parse(args); err != nil {
		return fmt.Errorf("cannot parse program arguments: %v", err)
	}

	var cmd func([]string) error

	switch a := f.Arg(0); a {
	case "make":
		cmd = cmdMake

	case "mock":
		cmd = cmdMock

	case "pack":
		cmd = cmdPack

	case "stub":
		cmd = cmdStub

	case "wrap":
		cmd = cmdWrap

	case "":
		return fmt.Errorf("no command")

	default:
		return fmt.Errorf("invalid command: %s", a)
	}

	return cmd(f.Args()[1:])
}

func cmdMake(args []string) error {
	var f = flag.NewFlagSet("turbine make", flag.ExitOnError)

	var (
		flagdata       = f.String("data", "", "The template data as JSON.")
		flagidentifier = f.String("identifier", "", "The declaration identifier.")
		flaginput      = f.String("input", "", "The input file. Defaults to standard input.")
		flagoutput     = f.String("output", "", "The output file. Defaults to standard output.")
		flagpackage    = f.String("package", "", "The declaration package.")
	)

	if err := f.Parse(args); err != nil {
		return fmt.Errorf("cannot parse program arguments: %v", err)
	}

	var data interface{}
	var err error

	if *flagpackage == "" || *flagidentifier == "" {
		if err = json.Unmarshal([]byte(*flagdata), &data); err != nil {
			return fmt.Errorf("cannot parse json: %v", err)
		}
	} else {
		if data, err = turbine.NewDecl(*flagpackage, *flagidentifier); err != nil {
			return fmt.Errorf("cannot parse declaration: %v", err)
		}
	}

	b, err := read(*flaginput)

	if err != nil {
		return err
	}

	if *flaginput == "" {
		*flaginput = "standard input"
	}

	if b, err = turbine.Generate(*flaginput, b, data); err != nil {
		return fmt.Errorf("cannot generate: %v", err)
	}

	return write(*flagoutput, b)
}

func cmdMock(args []string) error {
	return cmdDecl(args, "mock", mock)
}

func cmdPack(args []string) error {
	var f = flag.NewFlagSet("turbine pack", flag.ExitOnError)

	var (
		flagidentifier = f.String("identifier", "", "The declaration identifier.")
		flaginput      = f.String("input", "", "The input file. Defaults to standard input.")
		flagoutput     = f.String("output", "", "The output file. Defaults to standard output.")
		flagpackage    = f.String("package", "", "The declaration package.")
	)

	if err := f.Parse(args); err != nil {
		return fmt.Errorf("cannot parse program arguments: %v", err)
	}

	var b, err = read(*flaginput)

	if err != nil {
		return err
	}

	if *flaginput == "" {
		*flaginput = "standard input"
	}

	if *flagidentifier == "" {
		var id = strings.TrimSuffix(filepath.Base(*flagoutput), filepath.Ext(*flagoutput))

		for _, s := range suffixes {
			id = strings.TrimSuffix(id, s)
		}

		id = regexp.MustCompile("[^a-zA-Z0-9]+").ReplaceAllString(id, "")

		if !regexp.MustCompile("(?m:^)[a-zA-Z]").MatchString(id) {
			id = "resource" + id
		}

		*flagidentifier = id
	}

	if *flagpackage == "" {
		var d = filepath.Dir(*flagoutput)

		if p, err := build.ImportDir(d, 0); err == nil {
			*flagpackage = p.Name
		} else if b := filepath.Base(d); regexp.MustCompile(`(?m:^)\w+(?m:$)`).MatchString(b) {
			*flagpackage = b
		} else {
			*flagpackage = "main"
		}
	}

	var data = struct {
		Package, Identifier, Data string
	}{
		*flagpackage, *flagidentifier, fmt.Sprintf("%#v", string(b)),
	}

	if b, err = turbine.Generate(*flaginput, []byte(pack), data); err != nil {
		return fmt.Errorf("cannot generate: %v", err)
	}

	return write(*flagoutput, b)
}

func cmdStub(args []string) error {
	return cmdDecl(args, "stub", stub)
}

func cmdWrap(args []string) error {
	return cmdDecl(args, "wrap", stub)
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("turbine: ")

	if err := cmdMain(os.Args[1:]); err != nil {
		log.Fatalln("error:", err)
	}
}

func read(path string) ([]byte, error) {
	var b []byte
	var err error

	if path == "" {
		if b, err = ioutil.ReadAll(os.Stdin); err != nil {
			return nil, fmt.Errorf("cannot read standard input: %v", err)
		}
	} else {
		if b, err = ioutil.ReadFile(path); err != nil {
			return nil, fmt.Errorf("cannot read %s: %v", path, err)
		}
	}

	return b, nil
}

func write(path string, b []byte) error {
	var err error

	if path == "" {
		if _, err = os.Stdout.Write(b); err != nil {
			return fmt.Errorf("cannot write standard output: %v", err)
		}
	} else {
		if err = ioutil.WriteFile(path, b, 0600); err != nil {
			return fmt.Errorf("cannot write %s: %v", path, err)
		}
	}

	return nil
}
