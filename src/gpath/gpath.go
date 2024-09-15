package gpath

import (
	"os"
	"path"
)

type Path string

// resolve `~` to users home directory
func Expand(p string) string {
	if p[0] == '~' {
		home := os.Getenv("HOME")
		return path.Join(home, p[1:])
	}
	return p
}
