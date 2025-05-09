/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"github.com/sys-liqian/csi-driver-webdav/pkg/webdav"
	"k8s.io/klog/v2"
	"os"
)

var (
	endpoint              = flag.String("endpoint", "unix://tmp/csi.sock", "CSI endpoint")
	nodeID                = flag.String("nodeid", "", "node id")
	mountPermissions      = flag.Uint64("mount-permissions", 0, "mounted folder permissions")
	driverName            = flag.String("drivername", "", "name of the driver")
	workingMountDir       = flag.String("working-mount-dir", "/tmp/csi-storage", "working directory for provisioner to mount davfs shares temporarily")
	defaultOnDeletePolicy = flag.String("default-ondelete-policy", "", "default policy for deleting subdirectory when deleting a volume")
	directMount           = flag.Bool("direct-mount", false, "should the driver mount the share directly instead of creating a new pvc path")
)

func main() {
	klog.InitFlags(nil)
	_ = flag.Set("logtostderr", "true")
	flag.Parse()
	if *nodeID == "" {
		klog.Warning("nodeid is empty")
	}

	driverOptions := webdav.DriverOpt{
		Name:                  *driverName,
		NodeID:                *nodeID,
		Endpoint:              *endpoint,
		MountPermissions:      *mountPermissions,
		WorkingMountDir:       *workingMountDir,
		DefaultOnDeletePolicy: *defaultOnDeletePolicy,
	}
	d := webdav.NewDriver(&driverOptions)
	d.Run()
	os.Exit(0)
}
