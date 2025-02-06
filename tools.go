//go:build tools

// This file will never be built, but `go mod tidy` will see the packages
// imported here as dependencies and not remove them from `go.mod`.

package main

import (
	_ "github.com/kevinburke/go-bindata"
	// adding glog to override nested indirect dependency
	_ "github.com/golang/glog"
)
