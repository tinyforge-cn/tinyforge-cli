package commands

import (
	"fmt"

	"github.com/tinyforge-cn/cli/internal/version"
)

func VersionCommand() {
	fmt.Printf("tinyforge %s (%s)\n", version.Version, version.Commit)
}
