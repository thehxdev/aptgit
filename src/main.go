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
		log.Wrn.Println("using empty metadata...")
	}
	subcmdName := subcmd.Name()

	pkg, err := gpkg.Init(filepath.Join(genv.G.Gpkgs, fmt.Sprintf("%s.json", fPackage)))
	if err != nil {
		log.Err.Fatal(err)
	}

	pkg.TagName = fTagName
	if pkg.TagName == "latest" {
		pkg.TagName, err = pkg.GetLatestStableTag()
		if err != nil {
			log.Err.Fatal(err)
		}
	}

	pkg.Vars = map[string]string{
		"TAGNAME":      pkg.TagName,
		"VERSION":      pkg.ParseTagRegexp(pkg.TagName),
		"PLATFORM":     pkg.GetPlatform(),
		"ARCH":         pkg.GetArch(),
		"INSTALL_PATH": path.Join(genv.G.InstallPath, pkg.Info.Repository, pkg.TagName),
	}

	switch subcmdName {
	case "install":
		if pkg.TagName == "" {
			log.Err.Fatal("-tag flag is empty")
		}
		if err = pkg.Install(); err != nil {
			log.Err.Fatal(err)
		}
	case "latest":
		latestTag, err := pkg.GetLatestStableTag()
		if err != nil {
			log.Err.Fatal(err)
		}
		fmt.Println(latestTag)
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
		if err = pkg.SetTagNameAsMain(); err != nil {
			log.Err.Fatal(err)
		}
	case "uninstall":
		if pkg.TagName == "" {
			log.Err.Fatal("-tag flag is empty")
		}
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
