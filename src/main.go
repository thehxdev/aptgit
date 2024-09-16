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
	installPackage(pkgDef)
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

	savedFilePath, err := pdef.DownloadRelease(latestTag)
	if err != nil {
		return err
	}

	err = pdef.RunInstallSteps(map[string]string{
		"FILE":         savedFilePath,
		"INSTALL_PATH": path.Join(config.G.InstallPath, pdef.Repository),
	})
	if err != nil {
		return err
	}

	return pdef.SymlinkBinaryFiles()
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
