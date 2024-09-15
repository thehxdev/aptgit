package main

import (
	"flag"
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
	pdef.DownloadLatest(config.G.DownloadPath)
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
