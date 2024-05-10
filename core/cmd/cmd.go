package cmd

import (
	"fmt"
	"guide/core"
)

func UseCmd(code string) error {
	osType := core.ShowSys()
	switch osType {
	case "linux":
		err := core.LinuxC(code)
		if err != nil {
			return fmt.Errorf("ERROR: %s", err)
		}
		return nil
	case "windows":
		err := core.WinC(code)
		if err != nil {
			return fmt.Errorf("ERROR: %s", err)
		}
		return nil
	default:
		return fmt.Errorf("%s", "WARN: unsupported operating system.")
	}
}
