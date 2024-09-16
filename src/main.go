package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/thehxdev/aptgit/genv"
	"github.com/thehxdev/aptgit/gpkg"
	"github.com/thehxdev/aptgit/log"
)

var (
	fPackage          string
	fTagName          string
	fInstall          bool
	fSetGlobalVersion bool
	fGetLatestTag     bool
	fGetAllTags       bool
)

func init() {
	slapWindowsUsers()
	parseFlags()
	genv.Init()
}

func main() {
	if len(os.Args) < 2 {
		(flag.Usage)()
		os.Exit(1)
	}

	pkg, err := gpkg.Init(filepath.Join(genv.G.Gpkgs, fmt.Sprintf("%s.json", fPackage)))
	if err != nil {
		log.Err.Fatal(err)
	}
	pkg.TagName = fTagName

	pkg.Vars = map[string]string{
		"TAGNAME":      pkg.TagName,
		"VERSION":      pkg.ParseTagRegexp(pkg.TagName),
		"PLATFORM":     pkg.GetPlatform(),
		"ARCH":         pkg.GetArch(),
		"INSTALL_PATH": path.Join(genv.G.InstallPath, pkg.Info.Repository, pkg.TagName),
	}

	if fInstall {
		if pkg.TagName == "" {
			log.Inf.Println("-tag flag is empty. using latest stable version of the package...")
			pkg.TagName, err = pkg.GetLatestStableTag()
			if err != nil {
				log.Err.Fatal(err)
			}
		}
		err = pkg.Install()
		if err != nil {
			log.Err.Fatal(err)
		}
	} else if fGetLatestTag {
		latest, err := pkg.GetLatestStableTag()
		if err != nil {
			log.Err.Fatal(err)
		}
		fmt.Println(latest)
	} else if fGetAllTags {
		tags, err := pkg.GetAllTags()
		if err != nil {
			log.Err.Fatal(err)
		}
		for _, t := range tags {
			fmt.Println(t)
		}
	} else if fSetGlobalVersion {
		if pkg.TagName == "" {
			log.Err.Fatal("-tag flag is empty")
		}
		err = pkg.SetTagNameAsMain()
		if err != nil {
			log.Err.Fatal(err)
		}
	} else {
		(flag.Usage)()
		log.Err.Fatal("no valid operation")
	}
}

func parseFlags() {
	flag.StringVar(&fPackage, "p", "", "Package name")
	flag.StringVar(&fTagName, "tag", "", "Package tag name")
	flag.BoolVar(&fInstall, "install", false, "Install a package")
	flag.BoolVar(&fSetGlobalVersion, "set-global", false, "Set package global version")
	flag.BoolVar(&fGetAllTags, "list-all", false, "List all available tags for a package")
	flag.BoolVar(&fGetLatestTag, "latest", false, "Get latest tag of a package")
	flag.Parse()
}


func slapWindowsUsers() {
	if runtime.GOOS == "windows" {
		log.Err.Fatal("install Linux Bruh!")
	}
}
