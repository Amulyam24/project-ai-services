package configure

// ArgParam keys common to both Podman and OpenShift deployments.
const (
	ArgParamAdminPasswordHash = "backend.adminPasswordHash"
	ArgParamDBPassword        = "db.password"
	ArgParamWorkerGatewayPort = "backend.workerGatewayPort"
)

// ArgParam keys used only by the Podman deployment.
const (
	ArgParamRuntime               = "backend.runtime"
	ArgParamPodmanAuthFileContent = "backend.podman.authFileContent"
	ArgParamPodmanURI             = "backend.podman.uri"
	ArgParamCaddyHTTPSPort        = "caddy.httpsPort"
	ArgParamLocalWorker           = "backend.localWorker"
	ArgParamCaddyFileContent      = "caddy.caddyFileContent"
	ArgParamSSLCertFileContent    = "caddy.sslCertContent"
	ArgParamSSLKeyFileContent     = "caddy.sslKeyContent"
	// ArgParamSpyreEnabled controls whether the catalog backend requests
	// /dev/vfio devices and mounts /sys/kernel/iommu_groups. Set to "false"
	// on AMD GPU hosts where VFIO is absent.
	ArgParamSpyreEnabled = "backend.spyreEnabled"
)
