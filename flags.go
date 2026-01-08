package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	bold      = "\033[1m"
	underline = "\033[4m"
	reset     = "\033[0m"
)

type Flags struct {
	Directory string
	Version   bool
	Help      bool
}

func (f *Flags) Usage() {
	fmt.Print(customUsage())
}

func (f *Flags) define() {
	flag.StringVar(&f.Directory, "directory", ".", "Root directory to search (default: current directory)")
	flag.StringVar(&f.Directory, "d", ".", "Root directory to search (default: current directory)")

	flag.BoolVar(&f.Version, "v", false, "Print version")
	flag.BoolVar(&f.Version, "version", false, "Print version")

	flag.BoolVar(&f.Help, "h", false, "Print help")
	flag.BoolVar(&f.Help, "help", false, "Print help")
}

func customUsage() string {
	type flagGroup struct {
		names    []string
		defValue string
	}

	groups := make(map[string]*flagGroup)

	flag.VisitAll(func(f *flag.Flag) {
		if val, ok := groups[f.Usage]; !ok {
			groups[f.Usage] = &flagGroup{
				names:    []string{f.Name},
				defValue: f.DefValue,
			}
		} else {
			val.names = append(val.names, f.Name)
		}
	})

	var b strings.Builder

	fmt.Fprintf(&b, "%s%sUsage:%s %s [OPTIONS]\n\n", underline, bold, reset, filepath.Base(os.Args[0]))
	fmt.Fprintf(&b, "%s%sOptions:%s\n", underline, bold, reset)

	// TODO: Preserve insertion order.
	// Currently map iteration is random, causing inconsistent help output.
	for usage, fgroup := range groups {
		fmt.Fprintf(&b, "  ")
		for i, name := range fgroup.names {
			if i > 0 {
				fmt.Fprintf(&b, ", ")
			}

			prefix := "-"
			if len(name) > 1 {
				prefix = "--"
			}

			fmt.Fprintf(&b, "%s%s", prefix, name)
		}

		fmt.Fprintf(&b, "\n\t%s", usage)
		b.WriteString("\n")
	}

	return b.String()
}

func ParseFlags() *Flags {
	f := Flags{}
	f.define()

	flag.Usage = func() {
		fmt.Print(customUsage())
	}

	flag.Parse()
	return &f
}
