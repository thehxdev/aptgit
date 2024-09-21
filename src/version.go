package main

import "fmt"


const (
	V_MAJOR = 0
	V_MINOR = 0
	V_PATCH = 1
)

var VERSION_STRING = fmt.Sprintf("v%d.%d.%d", V_MAJOR, V_MINOR, V_PATCH)
