package main

import (
	"flag"
	"path"
	"runtime"

	"github.com/thehxdev/aptgit/config"
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
	if err := config.ReadConfig(configPath); err != nil {
		log.Err.Fatal(err)
	}
}

func main() {
	err := installPackage(pkgDef)
	if err != nil {
		log.Err.Fatal(err)
	}
}

func installPackage(defPath string) error {
	pdef, err := gpkg.ReadDefinitionFile(defPath)
	if err != nil {
		return err
	}

	latestTag, err := pdef.GetLatestStableTag()
	if err != nil {
		return err
	}

	pkgVars := map[string]string{
		"TAGNAME":  latestTag,
		"VERSION":  pdef.ParseTagRegexp(latestTag),
		"PLATFORM": pdef.GetPlatform(),
		"ARCH":     pdef.GetArch(),
	}

	savedFilePath, err := pdef.DownloadRelease(pkgVars)
	if err != nil {
		return err
	}

	pkgVars["FILE"] = savedFilePath
	pkgVars["INSTALL_PATH"] = path.Join(config.G.InstallPath, pdef.Repository, latestTag)

	err = pdef.RunCommands(pdef.InstallSteps, pkgVars)
	if err != nil {
		return err
	}

	return pdef.SymlinkBinaryFiles(pkgVars)
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
