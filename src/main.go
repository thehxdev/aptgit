package main

import (
	"flag"
	"path"
	"runtime"

	"github.com/thehxdev/aptgit/genv"
	"github.com/thehxdev/aptgit/gpkg"
	"github.com/thehxdev/aptgit/log"
)

var (
	configPath string
	pkgDef     string
)

func init() {
	slapWindowsUsers()
	parseFlags()
	genv.Init()
}

func main() {
	pkg, err := gpkg.Init(pkgDef)
	if err != nil {
		log.Err.Fatal(err)
	}

	// tags, err := pkg.GetAllTags()
	// if err != nil {
	// 	log.Err.Fatal(err)
	// }
	// pkg.TagName = tags[0]

	pkg.TagName, err = pkg.GetLatestStableTag()
	if err != nil {
		log.Err.Fatal(err)
	}

	pkg.Vars = map[string]string{
		"TAGNAME":      pkg.TagName,
		"VERSION":      pkg.ParseTagRegexp(pkg.TagName),
		"PLATFORM":     pkg.GetPlatform(),
		"ARCH":         pkg.GetArch(),
		"INSTALL_PATH": path.Join(genv.G.InstallPath, pkg.Info.Repository, pkg.TagName),
	}

	// err = pkg.SetTagNameAsMain()
	err = pkg.Install()
	if err != nil {
		log.Err.Fatal(err)
	}
}

func parseFlags() {
	flag.StringVar(&configPath, "c", "", "Path to aptgit config file")
	flag.StringVar(&pkgDef, "p", "", "Path to package definition file")
	flag.Parse()
}

func slapWindowsUsers() {
	if runtime.GOOS == "windows" {
		log.Err.Fatal("install Linux Bruh!")
	}
}
