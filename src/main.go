package main

import (
	"flag"
	"os"
	"path"
	"path/filepath"
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
	pkg, err := gpkg.Init(pkgDef)
	if err != nil {
		log.Err.Fatal(err)
	}

	tags, err := pkg.GetAllTags()
	if err != nil {
		log.Err.Fatal(err)
	}
	pkg.TagName = tags[0]

	// pkg.TagName, err = pkg.GetLatestStableTag()
	// if err != nil {
	// 	log.Err.Fatal(err)
	// }

	pkg.Vars = map[string]string{
		"TAGNAME":      pkg.TagName,
		"VERSION":      pkg.ParseTagRegexp(pkg.TagName),
		"PLATFORM":     pkg.GetPlatform(),
		"ARCH":         pkg.GetArch(),
		"INSTALL_PATH": path.Join(config.G.InstallPath, pkg.Info.Repository, pkg.TagName),
	}

	err = setTagNameAsMain(pkg)
	// err = installPackage(pkg)
	if err != nil {
		log.Err.Fatal(err)
	}
}

func installPackage(pkg *gpkg.Gpkg) error {
	var err error

	pkg.Vars["FILE"], err = pkg.DownloadRelease(pkg.Vars)
	if err != nil {
		return err
	}

	// pkgInstallPath := path.Join(config.G.InstallPath, pkg.Info.Repository, pkg.TagName)
	pkgInstallPath := pkg.Vars["INSTALL_PATH"]
	if _, err := os.Stat(pkgInstallPath); err == nil {
		os.Remove(pkgInstallPath)
	}

	err = os.MkdirAll(pkgInstallPath, 0775)
	if err != nil {
		return err
	}
	// pkg.Vars["INSTALL_PATH"] = pkgInstallPath

	err = gpkg.RunCommands(pkg.Info.InstallSteps, pkg.Vars)
	if err != nil {
		return err
	}

	removedExistingSymlinks(pkg.Info.Bins)
	return pkg.SymlinkBinaryFiles(pkg.Vars)
}

func removedExistingSymlinks(bins []string) {
	for _, bin := range bins {
		_, filename := filepath.Split(bin)
		err := os.Remove(filepath.Join(config.G.BinPath, filename))
		if err != nil {
			log.Err.Println(err)
		}
	}
}

func setTagNameAsMain(pkg *gpkg.Gpkg) error {
	// tagPath := path.Join(config.G.InstallPath, pkg.Info.Repository, tag)
	tagPath := pkg.Vars["INSTALL_PATH"]
	if _, err := os.Stat(tagPath); err != nil {
		log.Err.Printf("tag name %s is not installed for package %s", pkg.TagName, pkg.Info.Repository)
		return err
	}

	// pkg.Vars["INSTALL_PATH"] = path.Join(config.G.InstallPath, pkg.Info.Repository, tag)
	removedExistingSymlinks(pkg.Info.Bins)
	pkg.SymlinkBinaryFiles(pkg.Vars)

	return nil
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
