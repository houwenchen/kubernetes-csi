package main

import (
	"flag"

	"github.com/houwenchen/kubernetes-csi/pkg/driver"
	"k8s.io/klog/v2"
)

var (
	endpoint   = flag.String("endpoint", "unix://csi/csi.sock", "CSI Socket")
	driverName = flag.String("drivername", driver.DefaultDriverName, "name of driver")
	nodeID     = flag.String("nodeid", "", "node id")
	enableLVM  = flag.Bool("enablelvm", true, "choose the way to create volume")
)

var (
	version = "v0.0.1"
)

func init() {
	flag.Set("logtostderr", "true")
}

func main() {
	flag.Parse()

	cfg := &driver.Config{
		DriverName:    *driverName,
		EndPoint:      *endpoint,
		NodeID:        *nodeID,
		VendorVersion: version,
		EnableLVM:     *enableLVM,
	}

	csidriver, err := driver.NewCSIDriver(cfg)
	if err != nil {
		klog.Fatalf("create csi driver by config failed, err: %v\n", err)
	}

	if err = csidriver.Run(); err != nil {
		klog.Fatalf("csi driver run failed, err: %v\n", err)
	}
}
