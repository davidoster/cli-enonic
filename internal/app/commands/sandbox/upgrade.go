package sandbox

import (
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/enonic/cli-enonic/internal/app/util"
	"github.com/urfave/cli"
	"os"
)

var Upgrade = cli.Command{
	Name:    "upgrade",
	Aliases: []string{"up"},
	Usage:   "Upgrades the distribution version.",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "version, v",
			Usage: "Distro version to upgrade to.",
		},
		cli.BoolFlag{
			Name:  "all, a",
			Usage: "List all distro versions.",
		},
	},
	Action: func(c *cli.Context) error {

		sandbox, _ := EnsureSandboxExists(c, "No sandboxes found, do you want to create one?", "Select sandbox:", true, false)
		if sandbox == nil {
			os.Exit(0)
		}
		version := ensureVersionCorrect(c.String("version"), c.Bool("all"))
		preventVersionDowngrade(sandbox, version)

		sandbox.Distro = formatDistroVersion(version, util.GetCurrentOs())
		writeSandboxData(sandbox)
		fmt.Fprintf(os.Stdout, "Sandbox '%s' distro upgraded to '%s'.\n", sandbox.Name, sandbox.Distro)

		return nil
	},
}

func preventVersionDowngrade(sandbox *Sandbox, newString string) {
	distroVer := parseDistroVersion(sandbox.Distro, false)
	oldVer, err := semver.NewVersion(distroVer)
	util.Fatal(err, fmt.Sprintf("Could not parse sandbox distro version from: %s", sandbox.Distro))
	newVer, err2 := semver.NewVersion(newString)
	util.Fatal(err2, fmt.Sprintf("Could not parse new distro version from: %s", newString))

	if oldVer.Compare(newVer) > 0 {
		fmt.Fprintf(os.Stdout, "Sandbox '%s' already has newer distro version: %s\n", sandbox.Name, sandbox.Distro)
		os.Exit(0)
	}
}
