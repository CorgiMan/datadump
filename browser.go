package datadump

import (
	"fmt"
	"os/exec"
	"runtime"
)

var commands = map[string]string{
	"windows": "start",
	"darwin":  "open",
	"linux":   "xdg-open",
}

var Version = "0.1.0"

// Open calls the OS default program for uri
func OpenInBrowser(uri string) error {
	run, ok := commands[runtime.GOOS]
	if !ok {
		return fmt.Errorf("don't know how to open things on %s platform", runtime.GOOS)
	}

	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/C", run, "", uri)
		return cmd.Start()
	} else {
		cmd := exec.Command(run, uri)
		return cmd.Start()
	}
}
