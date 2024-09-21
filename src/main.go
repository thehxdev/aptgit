package main

import (
	"flag"
	"fmt"
	"path"
	"path/filepath"
	"runtime"

	"github.com/thehxdev/aptgit/genv"
	"github.com/thehxdev/aptgit/gpkg"
	"github.com/thehxdev/aptgit/log"
)

func init() {
	slapWindowsUsers()
	parseFlags()
	genv.Init()
}

func main() {
	allMds, err := gpkg.ReadMdFile()
	if err != nil {
		log.Wrn.Println(err)
		log.Wrn.Println("failed to read packages metadata file. using default empy metadata...")
	}

	pkg, err := gpkg.Init(filepath.Join(genv.G.Gpkgs, fmt.Sprintf("%s.json", fPackage)))
	if err != nil {
		log.Err.Fatal(err)
	}

	pkg.TagName = fTagName
	if pkg.TagName == "" {
		log.Err.Fatal("-tag flag is empty")
	}

	if pkg.TagName == "latest" {
		pkg.TagName, err = pkg.GetLatestStableTag()
		if err != nil {
			log.Err.Fatal(err)
		}
		if subcmd.Name() == "latest" {
			fmt.Println(pkg.TagName)
			return
		}
	}

	pkg.Vars = map[string]string{
		"TAGNAME":      pkg.TagName,
		"VERSION":      pkg.ParseTagRegexp(pkg.TagName),
		"PLATFORM":     pkg.GetPlatform(),
		"ARCH":         pkg.GetArch(),
		"INSTALL_PATH": path.Join(genv.G.InstallPath, pkg.Info.Repository, pkg.TagName),
	}

	switch subcmd.Name() {
	case "install":
		if err = pkg.Install(); err != nil {
			log.Err.Fatal(err)
		}
	case "list-all":
		tags, err := pkg.GetAllTags()
		if err != nil {
			log.Err.Fatal(err)
		}
		for _, t := range tags {
			fmt.Println(t)
		}
	case "global":
		if err = pkg.SetTagNameAsMain(); err != nil {
			log.Err.Fatal(err)
		}
	case "uninstall":
		if err = pkg.Uninstall(); err != nil {
			log.Err.Println(err)
		}
	default:
		printUsage()
		log.Err.Fatal("no valid operation")
	}

	allMds[pkg.Info.Repository] = pkg.MainTag
	err = gpkg.WriteMdFile(allMds)
	if err != nil {
		log.Err.Fatal(err)
	}
}

func registerFlagSet(name string) {
	flagSets[name] = flag.NewFlagSet(name, flag.ExitOnError)
}

func slapWindowsUsers() {
	if runtime.GOOS == "windows" {
		log.Err.Fatal("install Linux Bruh!")
	}
}
