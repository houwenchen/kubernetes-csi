package driver

import (
	"context"
	"errors"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"k8s.io/klog/v2"
)

/*
实现下面的接口
type IdentityServer interface {
    GetPluginInfo(context.Context, *GetPluginInfoRequest) (*GetPluginInfoResponse, error)
    GetPluginCapabilities(context.Context, *GetPluginCapabilitiesRequest) (*GetPluginCapabilitiesResponse, error)
    Probe(context.Context, *ProbeRequest) (*ProbeResponse, error)
}
*/

var (
	defaultPluginCapability_Service_Types = []csi.PluginCapability_Service_Type{
		csi.PluginCapability_Service_CONTROLLER_SERVICE,
	}
)

type CSIIdentityServer struct {
	driver       *CSIDriver
	capabilities []*csi.PluginCapability
}

func NewDefaultCSIIdentityServer(driver *CSIDriver) *CSIIdentityServer {
	capabilities := make([]*csi.PluginCapability, 0)

	for _, svcType := range defaultPluginCapability_Service_Types {
		cap := &csi.PluginCapability{
			Type: &csi.PluginCapability_Service_{
				Service: &csi.PluginCapability_Service{
					Type: svcType,
				},
			}}

		capabilities = append(capabilities, cap)
	}

	return NewCSIIdentityServerWithOpt(driver, capabilities)
}

func NewCSIIdentityServerWithOpt(driver *CSIDriver, opts ...[]*csi.PluginCapability) *CSIIdentityServer {
	capabilities := make([]*csi.PluginCapability, 0)

	for _, opt := range opts {
		capabilities = append(capabilities, opt...)
	}

	return &CSIIdentityServer{
		driver:       driver,
		capabilities: capabilities,
	}
}

func (cis *CSIIdentityServer) GetPluginInfo(ctx context.Context, req *csi.GetPluginInfoRequest) (*csi.GetPluginInfoResponse, error) {
	klog.Info("start GetPluginInfo function")

	if cis.driver.config.DriverName == "" {
		klog.Fatal("miss driver name")
		return nil, errors.New("miss driver name")
	}

	if cis.driver.config.VendorVersion == "" {
		klog.Fatal("miss vendor version")
		return nil, errors.New("miss vendor version")
	}

	return &csi.GetPluginInfoResponse{
		Name:          cis.driver.config.DriverName,
		VendorVersion: cis.driver.config.VendorVersion,
	}, nil
}

func (cis *CSIIdentityServer) GetPluginCapabilities(ctx context.Context, req *csi.GetPluginCapabilitiesRequest) (*csi.GetPluginCapabilitiesResponse, error) {
	klog.Info("start GetPluginCapabilities function")

	return &csi.GetPluginCapabilitiesResponse{
		Capabilities: cis.capabilities,
	}, nil
}

// 这个慎用，一般情况下不用具体实现
func (cis *CSIIdentityServer) Probe(ctx context.Context, req *csi.ProbeRequest) (*csi.ProbeResponse, error) {
	klog.Info("start Probe function")

	return &csi.ProbeResponse{}, nil
}
