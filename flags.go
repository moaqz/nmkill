package main

import "flag"

type Flags struct {
	Directory string
	Version   bool
	Debug     bool
}

func (f *Flags) define() {
	flag.StringVar(&f.Directory, "directory", ".", "Set the directory from which to begin searching.")
	flag.StringVar(&f.Directory, "d", ".", "Set the directory from which to begin searching.")

	flag.BoolVar(&f.Debug, "debug", false, "Enable debug mode")

	flag.BoolVar(&f.Version, "v", false, "Print version information and quit")
	flag.BoolVar(&f.Version, "version", false, "Print version information and quit")
}

func ParseFlags() *Flags {
	flags := Flags{}

	flags.define()
	flag.Parse()

	return &flags
}
