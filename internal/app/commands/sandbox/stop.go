package sandbox

import (
	"github.com/urfave/cli"
	"fmt"
	"os"
)

var Stop = cli.Command{
	Name:  "stop",
	Usage: "Stop the sandbox started in detached mode.",
	Action: func(c *cli.Context) error {

		sData := readSandboxesData()
		if sData.Running == "" || sData.PID == 0 {
			fmt.Fprintln(os.Stderr, "No sandbox is currently running in detached mode.")
			os.Exit(0)
		}

		stopDistro(sData.PID)
		writeRunningSandbox("", 0)

		fmt.Fprintf(os.Stderr, "Sandbox %s stopped\n", sData.Running)

		return nil
	},
}
