package zl_mumble

import (
	"errors"
	"fmt"
	"os/exec"
)

func PlayWavLocal(filepath string) error {
	cmd := exec.Command("/usr/bin/aplay", filepath)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(fmt.Sprintf("alert: cmd.Run() for aplay failed with %s\n", err))
	}
	return nil
}