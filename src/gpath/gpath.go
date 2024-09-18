package gpath

import (
	"fmt"
	"os"
	"path"

	"github.com/thehxdev/aptgit/log"
)

type Path string

// resolve `~` to users home directory
func Expand(p string) string {
	if p[0] == '~' {
		home := os.Getenv("HOME")
		if home == "" {
			log.Err.Fatal("HOME environment variable is not set")
		}
		return path.Join(home, p[1:])
	}
	return p
}

func Qoute(p string) string {
	return fmt.Sprintf("'%s'", p)
}

func Exist(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}
