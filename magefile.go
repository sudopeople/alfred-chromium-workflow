//go:build mage
// +build mage

package main

import (
	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
	"github.com/magefile/mage/sh"
)

var (
	// Default target to run when none is specified
	Default = Test
	Aliases = map[string]interface{}{
		"cover": Coverage,
		"c":     Coverage,
		"t":     Test,
		"l":     Lint,
		"f":     Formatter,
		"b":     Builder,
	}

	test = sh.RunCmd("go", "test", "-v")
)

// Lint run golint and golang-ci-lint
func Lint() error {
	if err := sh.Run("golint", "-set_exit_status", "./..."); err != nil {
		return err
	}
	return sh.Run("golangci-lint", "run", "-c", ".golangci.toml")
}

// Test run unit tests
func Test() error {
	return test("-coverprofile=coverage.out", "./...")
}

// Coverage run unit tests and open coverage
func Coverage() error {
	mg.Deps(Test)
	return sh.Run("go", "tool", "cover", "-html=coverage.out")
}

// Clean remove cache files
func Clean() error { return sh.Run("mage", "-clean") }

func Formatter() error { return sh.Run("gofmt", "-w", "src", "main.go") }

func Builder() error {
	return sh.RunWith(map[string]string{
		"GOOS":        "darwin",
		"GOARCH":      "amd64",
		"CGO_ENABLED": "1",
	}, "go", "build", "-o", "alfred-chromium-workflow", ".")
}
