package podman

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/project-ai-services/ai-services/internal/pkg/accelerator/amdgpu"
	"github.com/project-ai-services/ai-services/internal/pkg/logger"
)

const (
	podmanSocketWaitDuration = 2 * time.Second
	contextTimeout           = 30 * time.Second
)

// Configure performs the complete configuration of the Podman environment.
func (p *PodmanBootstrap) Configure(ctx context.Context) error {
	euid := os.Geteuid()
	if euid != 0 {
		return fmt.Errorf("podman bootstrap requires root privileges, either run as root or use sudo")
	}

	// Cache once — avoids two os.Stat("/dev/kfd") syscalls further down.
	hasAMDGPU := amdgpu.IsApplicable()

	// 1. Install and configure Podman if not done
	if err := ensurePodmanInstalled(ctx); err != nil {
		return err
	}

	if err := configurePodman(ctx); err != nil {
		return err
	}

	// 2. Configure user groups (sentient group)
	if err := ensureUsergroupConfigured(ctx); err != nil {
		return err
	}

	// 3. Accelerator setup — Spyre or ROCm (mutually exclusive paths).
	// For AMD GPU hosts there is no repair work to do at configure time;
	// ROCm pre-flight checks are handled entirely by the Validate phase
	// (ROCmRegistry rules) which correctly honours --skip-validation flags.
	if hasAMDGPU {
		logger.Infoln("AMD GPU detected, skipping Spyre configuration")
	} else {
		if err := ensureSpyreConfigured(ctx); err != nil {
			return err
		}
	}

	// 4. Configure ulimits (memlock and nofile)
	if err := ensureUlimitsConfigured(ctx); err != nil {
		return err
	}

	// 5. Configure systemd user slice limits for rootless podman
	if err := ensureSystemdSliceLimitsConfigured(ctx); err != nil {
		return err
	}

	// 6. Configure SMT level to 2 and persist via systemd.
	// Skip for ROCm systems — SMT tuning for AMD GPUs is under investigation.
	if !hasAMDGPU {
		if err := ensureSMTConfigured(ctx); err != nil {
			return err
		}
	}

	// 7. Configure SELinux policy for Podman socket access
	if err := ensureSELinuxPolicyConfigured(ctx); err != nil {
		return err
	}

	logger.Infoln("LPAR configured successfully")

	return nil
}

// Made with Bob
