package driver

import (
	"errors"

	"github.com/houwenchen/kubernetes-csi/pkg/config"
	"k8s.io/klog/v2"
)

type CSIDriver struct {
	config *config.Config
}

func NewCSIDriver(cfg *config.Config) (*CSIDriver, error) {
	if cfg.DriverName == "" {
		klog.Infof("cofig doesn't have filed DriverName")
		return nil, errors.New("miss field DriverName")
	}

	if len(cfg.NodeID) == 0 {
		klog.Fatalf("config doesn't have filed NodeID")
		return nil, errors.New("miss filed NodeID")
	}

	if cfg.EndPoint == "" {
		klog.Fatalf("config doesn't have filed EndPoint")
		return nil, errors.New("miss driver endpoint")
	}

	return &CSIDriver{
		config: cfg,
	}, nil
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
