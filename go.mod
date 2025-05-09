module github.com/sys-liqian/csi-driver-webdav

go 1.23.0

toolchain go1.24.1

require (
	github.com/container-storage-interface/spec v1.11.0
	github.com/moby/sys/mountinfo v0.7.2
	golang.org/x/sys v0.33.0
	google.golang.org/grpc v1.72.0
	google.golang.org/protobuf v1.36.6
	k8s.io/klog/v2 v2.130.1
	k8s.io/utils v0.0.0-20250502105355-0f33e8f1c979
	sigs.k8s.io/yaml v1.4.0
)

require (
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/rogpeppe/go-internal v1.10.0 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250505200425-f936aa4a68b2 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)
