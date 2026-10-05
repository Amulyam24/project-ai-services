package amd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/project-ai-services/ai-services/internal/pkg/logger"
)

// IsApplicable returns true when AMD GPU devices are present on the host.
// It checks for /dev/kfd which is created by the amdgpu kernel module and
// is required for any ROCm workload.
func IsApplicable() bool {
	_, err := os.Stat("/dev/kfd")
	return err == nil
}

// ListDevices returns the list of AMD GPU PCI addresses found via lspci.
// AMD GPU PCI vendor ID is 1002. Only display/compute class devices are counted
// (PCI class 0300 VGA Compatible Controller, 0302 3D Controller, 0380 Display
// Controller) — this excludes AMD audio devices (class 0403, printed as
// "Audio device") that share the same vendor ID and would otherwise inflate
// the GPU count.
//
// lspci output format (one device per line):
//
//	<addr> <Class>: <Vendor> <Description>
//	e.g. "0b:00.0 Display controller: Advanced Micro Devices..."
func ListDevices(ctx context.Context) ([]string, error) {
	cmd := exec.Command("lspci", "-d", "1002:")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to list AMD GPU PCI devices: %w, output: %s", err, string(out))
	}

	// gpuClasses are the substrings lspci uses in the class field of its
	// single-line output for GPU compute/display devices.
	gpuClasses := []string{
		"VGA compatible controller",
		"3D controller",
		"Display controller",
	}

	var devices []string
	for _, line := range strings.Split(string(out), "\n") {
		if line == "" {
			continue
		}
		// Each line is "<addr> <class>: …". Extract the class portion between
		// the first space and the colon.
		rest := strings.SplitN(line, " ", 2)
		if len(rest) < 2 {
			continue
		}
		pciAddr := rest[0]
		classAndDesc := rest[1]
		colonIdx := strings.Index(classAndDesc, ":")
		if colonIdx < 0 {
			continue
		}
		lineClass := classAndDesc[:colonIdx]
		for _, gpuClass := range gpuClasses {
			if strings.EqualFold(lineClass, gpuClass) {
				logger.DebugfCtx(ctx, "AMD GPU detected, PCI address: %s\n", pciAddr)
				devices = append(devices, pciAddr)
				break
			}
		}
	}

	logger.DebuglnCtx(ctx, "AMD GPU devices found: "+strings.Join(devices, ", "))

	return devices, nil
}
