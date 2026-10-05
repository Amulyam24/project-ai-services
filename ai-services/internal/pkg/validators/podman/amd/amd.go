package amd

import (
	"context"
	"fmt"
	"os"

	amdAccelerator "github.com/project-ai-services/ai-services/internal/pkg/accelerator/amd"
	"github.com/project-ai-services/ai-services/internal/pkg/constants"
	"github.com/project-ai-services/ai-services/internal/pkg/logger"
)

// AMDRule validates that AMD GPU devices are present and accessible.
type AMDRule struct{}

func NewAMDRule() *AMDRule {
	return &AMDRule{}
}

func (r *AMDRule) Name() string {
	return "amd-gpu"
}

func (r *AMDRule) Description() string {
	return "Validates that AMD GPU devices (/dev/kfd, /dev/dri) are present and accessible."
}

func (r *AMDRule) Verify(ctx context.Context) error {
	logger.Debugln("Validating AMD GPU availability...")

	if !amdAccelerator.IsApplicable() {
		return fmt.Errorf("AMD GPU not detected: /dev/kfd not found (is the amdgpu kernel module loaded?)")
	}

	// /dev/dri must also exist for ROCm rendering/compute contexts
	if _, err := os.Stat("/dev/dri"); err != nil {
		return fmt.Errorf("AMD GPU device directory /dev/dri not found: %w", err)
	}

	devices, err := amdAccelerator.ListDevices(ctx)
	if err != nil {
		return fmt.Errorf("failed to enumerate AMD GPU devices: %w", err)
	}

	if len(devices) == 0 {
		return fmt.Errorf("no AMD GPU devices found via lspci")
	}

	logger.Debugf("✓ AMD GPU validation passed: %d device(s) detected\n", len(devices))

	return nil
}

func (r *AMDRule) Message() string {
	return "AMD GPU devices are present and accessible"
}

func (r *AMDRule) Level() constants.ValidationLevel {
	return constants.ValidationLevelError
}

func (r *AMDRule) Hint() string {
	return "Ensure the amdgpu kernel module is loaded ('modprobe amdgpu') and the user has access to /dev/kfd and /dev/dri."
}
