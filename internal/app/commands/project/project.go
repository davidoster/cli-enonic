package project

import (
	"github.com/urfave/cli"
	"github.com/enonic/xp-cli/internal/app/util"
	"os"
	"github.com/enonic/xp-cli/internal/app/commands/sandbox"
	"fmt"
	"os/exec"
	"path"
	"bytes"
)

func All() []cli.Command {
	return []cli.Command{
		Create,
		Sandbox,
		Clean,
		Build,
		Deploy,
		Install,
	}
}

type ProjectData struct {
	Sandbox string `toml:"sandbox"`
}

func readProjectData() ProjectData {
	file := util.OpenOrCreateDataFile(".enonic", true)
	defer file.Close()

	var data ProjectData
	util.DecodeTomlFile(file, &data)
	return data
}

func writeProjectData(data ProjectData) {
	file := util.OpenOrCreateDataFile(".enonic", false)
	defer file.Close()

	util.EncodeTomlFile(file, data)
}

func getOsGradlewFile() string {
	gradlewFile := "gradlew"
	switch util.GetCurrentOs() {
	case "windows":
		gradlewFile += ".bat"
	case "mac", "linux":
		gradlewFile = "./" + gradlewFile

	}
	return gradlewFile
}

func ensureDirHasGradleFile() {
	dir, err := os.Getwd()
	util.Fatal(err, "Could not get current dir")

	if _, err := os.Stat(path.Join(dir, getOsGradlewFile())); os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, "Not a valid project folder")
		os.Exit(0)
	}
}

func ensureProjectDataExists(c *cli.Context, noBoxMessage string) ProjectData {

	ensureDirHasGradleFile()

	projectData := readProjectData()
	badSandbox := projectData.Sandbox == "" || !sandbox.Exists(projectData.Sandbox)
	argExist := c != nil && c.NArg() > 0
	if badSandbox || argExist {
		sbox := sandbox.EnsureSandboxNameExists(c, noBoxMessage, "Select a sandbox to use:")
		projectData.Sandbox = sbox.Name
		if badSandbox {
			writeProjectData(projectData)
			fmt.Fprintf(os.Stderr, "Sandbox '%s' set as default for this project. You can change it using 'project sandbox' command at any time.\n", projectData.Sandbox)
		}
	}
	return projectData
}

func runGradleTask(projectData ProjectData, task, message string) {

	sandboxData := sandbox.ReadSandboxData(projectData.Sandbox)

	javaHome := fmt.Sprintf("-Dorg.gradle.java.home=%s", sandbox.GetDistroJdkPath(sandboxData.Distro))
	xpHome := fmt.Sprintf("-Dxp.home=%s", sandbox.GetSandboxHomePath(projectData.Sandbox))

	var stderr bytes.Buffer
	fmt.Fprint(os.Stderr, message)
	command := getOsGradlewFile()
	cmd := exec.Command(command, task, javaHome, xpHome)
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "\n"+stderr.String())
	} else {
		fmt.Fprintln(os.Stderr, "Done")
	}
}
