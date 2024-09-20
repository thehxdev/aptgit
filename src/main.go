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

var (
	flagSets = make(map[string]*flag.FlagSet)
	subcmd   *flag.FlagSet
	fPackage string
	fTagName string
)

func init() {
	slapWindowsUsers()
	parseFlags()
	genv.Init()
}

func main() {
	allMds, err := gpkg.ReadMdFile()
	if err != nil {
		log.Err.Fatal(err)
	}

	pkg, err := gpkg.Init(filepath.Join(genv.G.Gpkgs, fmt.Sprintf("%s.json", fPackage)))
	if err != nil {
		log.Err.Fatal(err)
	}

	pkg.TagName = fTagName
	if pkg.TagName == "" || pkg.TagName == "latest" {
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
		err = pkg.Install()
		if err != nil {
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
		if pkg.TagName == "" {
			log.Err.Fatal("-tag flag is empty")
		}
		err = pkg.SetTagNameAsMain()
		if err != nil {
			log.Err.Fatal(err)
		}
	default:
		(flag.Usage)()
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
