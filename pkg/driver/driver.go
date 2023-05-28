package driver

import (
	"errors"
	"sync"

	"google.golang.org/grpc/balancer/grpclb/state"
	"k8s.io/klog/v2"
)

var (
	DefaultDriverName string = "csidriver.whou.io"
)

type CSIDriver struct {
	config Config

	// gRPC calls involving any of the fields below must be serialized
	// by locking this mutex before starting. Internal helper
	// functions assume that the mutex has been locked.
	mutex sync.Mutex
	state state.State
}

type Config struct {
	DriverName    string //必须要有的，GetPluginInfo 会用到
	EndPoint      string
	NodeID        string
	VendorVersion string //必须要有的，GetPluginInfo 会用到

	VolumeDir string

	EnableLVM bool
}

func NewCSIDriver(cfg *Config) (*CSIDriver, error) {
	if cfg.DriverName == "" {
		klog.Infof("cofig doesn't have filed DriverName, use default DriverName: %s\n", DefaultDriverName)
		cfg.DriverName = DefaultDriverName
	}

	if len(cfg.NodeID) == 0 {
		klog.Fatalf("config doesn't have filed NodeID")
		return nil, errors.New("miss filed NodeID")
	}

	if cfg.EndPoint == "" {
		klog.Fatalf("config doesn't have filed EndPoint")
		return nil, errors.New("miss driver endpoint")
	}

	if err := makeVolumeDir(cfg.VolumeDir); err != nil {
		return nil, err
	}

	return nil, nil
}

func (d *CSIDriver) Run() error {
	s := NewNonBlockingGRPCServer()

	ids := NewDefaultCSIIdentityServer(d)
	cs := NewDefaultCSIControllerServer(d)
	ns := NewDefaultCSINodeServer(d)

	s.Start(d.config.EndPoint, ids, cs, ns)
	s.Wait()

	return nil
}
