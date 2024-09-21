package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	flagSets = make(map[string]*flag.FlagSet)
	subcmd   *flag.FlagSet
	fPackage string
	fTagName string
)

func parseFlags() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	if os.Args[1] == "help" {
		printUsage()
		os.Exit(0)
	} else if os.Args[1] == "version" {
		fmt.Println("aptgit", VERSION_STRING)
		os.Exit(0)
	}

	registerFlagSet("install")
	registerFlagSet("uninstall")
	registerFlagSet("global")
	registerFlagSet("list-all")
	registerFlagSet("latest")

	for _, n := range []string{"install", "global", "uninstall"} {
		flagSets[n].StringVar(&fTagName, "tag", "", "Package tag name")
	}

	for _, n := range []string{"install", "global", "list-all", "latest", "uninstall"} {
		flagSets[n].StringVar(&fPackage, "p", "", "Package name")
	}

	subcmd = whichSubcmd()
	if subcmd == nil {
		printUsage()
		os.Exit(1)
	}
	subcmd.Parse(os.Args[2:])
}

func printUsage() {
	fmt.Fprint(os.Stderr, `aptgit Usage:
    install -p <package> -tag <tag name>
        install a package
		you can use '-tag latest' to install latest stable tag

    global -p <package> -tag <tag name>
        set global version of a pacakge

    latest -p <pacakge>
        get latest tag name of a package

    list-all -p <pacakge>
        list all tag names available to install

    uninstall -p <pacakge> -tag <tag name>
        uninstall a package
        you can use '-tag all' to uninstall all the versions installed

    help
        show this help message

    version
        show version information
`)
}

func whichSubcmd() *flag.FlagSet {
	if f, ok := flagSets[os.Args[1]]; ok {
		return f
	}
	return nil
}
