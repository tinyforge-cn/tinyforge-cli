package commands

import (
	"fmt"

	"github.com/bytecafe-run/cli/internal/version"
)

func VersionCommand() {
	fmt.Printf("bytecafe %s (%s)\n", version.Version, version.Commit)
}
