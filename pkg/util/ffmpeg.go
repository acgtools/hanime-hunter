package util

import (
	"bytes"
	"fmt"
	"os/exec"
)

func MergeToMP4(fileListPath, outputPath string) error {
	cmd := exec.Command("ffmpeg", "-f", "concat", "-safe", "0", "-i", fileListPath, "-c", "copy", outputPath, "-y")
	var errOutput bytes.Buffer
	cmd.Stderr = &errOutput

	cmd.Stdout = nil

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("%w\n%s", err, errOutput.String())
	}

	return nil
}
