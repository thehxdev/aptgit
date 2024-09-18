package genv

import (
	"os"
	"path"
	"path/filepath"

	"github.com/thehxdev/aptgit/gpath"
	"github.com/thehxdev/aptgit/log"
)

type Env struct {
	Home string
	InstallPath string
	DownloadPath string
	BinPath string
	Gpkgs string
	LockFile string
}

var G *Env = &Env{}

func Init() {
	home := os.Getenv("APTGIT_HOME")
	if home == "" {
		home = gpath.Expand("~/.aptgit")
	}
	if !path.IsAbs(home) {
		log.Err.Fatal("APTGIT_HOME environment varibale must be an absolute path")
	}
	if !gpath.Exist(home) {
		log.Err.Fatalf("APTGIT_HOME (%s) does not exist", home)
	}

	G.Home = home
	G.InstallPath = path.Join(home, "installed")
	G.DownloadPath = path.Join(home, "downloads")
	G.BinPath = path.Join(home, "bin")
	G.Gpkgs = path.Join(home, "gpkgs")
	G.LockFile = filepath.Join(home, "aptgit.lock")
}
