package internal

import (
	"runtime/debug"
	"strings"
)

// Version represents the current version of dasel.
// The real version number is injected at build time using ldflags.
var Version = "development"

func init() {
	// Version is set by ldflags on build.
	if Version != "development" {
		// Strip the "v" prefix to normalize version output across build methods.
		// The Dockerfile passes the raw git tag (e.g. "v3.10.1") while Homebrew
		// passes the bare semver (e.g. "3.10.1").
		Version = strings.TrimPrefix(Version, "v")
		return
	}

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}

	// https://github.com/golang/go/issues/29228
	if info.Main.Version == "(devel)" || info.Main.Version == "" {
		return
	}

	Version += "-" + info.Main.Version
}
