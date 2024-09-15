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
	pdef, err := gpkg.ReadDefinitionFile(pkgDef)
	if err != nil {
		log.Err.Fatal(err)
	}

	latestTag, err := pdef.GetLatestTag()
	if err != nil {
		log.Err.Fatal(err)
	}

	savedFilePath, err := pdef.DownloadRelease(latestTag)
	if err != nil {
		log.Err.Fatal(err)
	}

	err = pdef.RunInstallSteps(map[string]string{
		"FILE":         savedFilePath,
		"INSTALL_PATH": path.Join(config.G.InstallPath, pdef.Repository),
	})
	if err != nil {
		log.Err.Fatal(err)
	}

	err = pdef.LinkBinaryFiles()
	if err != nil {
		log.Err.Fatal(err)
	}
}

func parseFlags() {
	flag.StringVar(&configPath, "c", "", "Path to aptgit config file")
	flag.StringVar(&pkgDef, "def", "", "Package definition path")
	flag.Parse()
}

func slapWindowsUsers() {
	if runtime.GOOS == "windows" {
		log.Err.Fatal("install Linux Bruh!")
	}
}
