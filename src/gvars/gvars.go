package gvars

import (
	"fmt"
	"strings"
)

const DELIM byte = '%'

func ResolveAll(s string, vars map[string]string) string {
	sGvars := findGvars(s)
	for _, sgvar := range sGvars {
		if resolved, ok := vars[sgvar]; ok {
			gvar := fmt.Sprintf("%c%s%c", DELIM, sgvar, DELIM)
			s = strings.ReplaceAll(s, gvar, resolved)
		}
	}
	return s
}

func findGvars(s string) []string {
	vars := make([]string, 0)
	for i := 1; i < len(s); i++ {
		if s[i-1] == DELIM && s[i] != DELIM {
			for j := i; j < len(s); j++ {
				if s[j] == DELIM {
					vars = append(vars, s[i:j])
					break
				}
			}
		}
	}
	return vars
}
