package genv

import (
	"os"
	"path"

	"github.com/thehxdev/aptgit/gpath"
	"github.com/thehxdev/aptgit/log"
)

type Env struct {
	Home string

	// These other paths MUST be relative to `Home` field
	BinPath      string
	InstallPath  string
	DownloadPath string
	Gpkgs        string
}

var G *Env = &Env{}

func Init() {
	G.Home = os.Getenv("APTGIT_HOME")
	if G.Home == "" {
		G.Home = gpath.Expand("~/.aptgit")
	}
	if !path.IsAbs(G.Home) {
		log.Err.Fatal("APTGIT_HOME environment varibale must be an absolute path")
	}
	if _, err := os.Stat(G.Home); err != nil {
		log.Err.Fatalf("APTGIT_HOME (%s) does not exist", G.Home)
	}
	G.InstallPath = path.Join(G.Home, "installs")
	G.DownloadPath = path.Join(G.Home, "downloads")
	G.BinPath = path.Join(G.Home, "bin")
	G.Gpkgs = path.Join(G.Home, "gpkgs")
}
