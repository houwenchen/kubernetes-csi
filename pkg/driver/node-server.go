package hostpath

import (
	"context"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"k8s.io/klog/v2"
)

/*
实现下面的接口
type NodeServer interface {
    NodeStageVolume(context.Context, *NodeStageVolumeRequest) (*NodeStageVolumeResponse, error)
    NodeUnstageVolume(context.Context, *NodeUnstageVolumeRequest) (*NodeUnstageVolumeResponse, error)
    NodePublishVolume(context.Context, *NodePublishVolumeRequest) (*NodePublishVolumeResponse, error)
    NodeUnpublishVolume(context.Context, *NodeUnpublishVolumeRequest) (*NodeUnpublishVolumeResponse, error)
    NodeGetVolumeStats(context.Context, *NodeGetVolumeStatsRequest) (*NodeGetVolumeStatsResponse, error)
    NodeExpandVolume(context.Context, *NodeExpandVolumeRequest) (*NodeExpandVolumeResponse, error)
    NodeGetCapabilities(context.Context, *NodeGetCapabilitiesRequest) (*NodeGetCapabilitiesResponse, error)
    NodeGetInfo(context.Context, *NodeGetInfoRequest) (*NodeGetInfoResponse, error)
}
*/

var (
	defaultNodeServiceCapability_RPC_Types = []csi.NodeServiceCapability_RPC_Type{
		csi.NodeServiceCapability_RPC_UNKNOWN,
	}
)

type CSINodeServer struct {
	driver       *HostPathDriver
	capabilities []*csi.NodeServiceCapability
}

func NewDefaultCSINodeServer(driver *HostPathDriver) *CSINodeServer {
	capabilities := make([]*csi.NodeServiceCapability, 0)

	for _, RPCType := range defaultNodeServiceCapability_RPC_Types {
		cap := &csi.NodeServiceCapability{
			Type: &csi.NodeServiceCapability_Rpc{
				Rpc: &csi.NodeServiceCapability_RPC{
					Type: RPCType,
				},
			}}

		capabilities = append(capabilities, cap)
	}

	return NewCSINodeServerWithOpt(driver, capabilities)
}

func NewCSINodeServerWithOpt(driver *HostPathDriver, opts ...[]*csi.NodeServiceCapability) *CSINodeServer {
	capabilities := make([]*csi.NodeServiceCapability, 0)

	for _, opt := range opts {
		capabilities = append(capabilities, opt...)
	}

	return &CSINodeServer{
		driver:       driver,
		capabilities: capabilities,
	}
}

// capabilities 中有 NodeServiceCapability_RPC_STAGE_UNSTAGE_VOLUME 时才需要实现此方法
func (cns *CSINodeServer) NodeStageVolume(ctx context.Context, req *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	klog.Info("start NodeStageVolume function")

	return &csi.NodeStageVolumeResponse{}, nil
}

// capabilities 中有 NodeServiceCapability_RPC_STAGE_UNSTAGE_VOLUME 时才需要实现此方法
func (cns *CSINodeServer) NodeUnstageVolume(ctx context.Context, req *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	klog.Info("start NodeUnstageVolume function")

	return &csi.NodeUnstageVolumeResponse{}, nil
}

func (cns *CSINodeServer) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	klog.Info("start NodePublishVolume function")

	return &csi.NodePublishVolumeResponse{}, nil
}

func (cns *CSINodeServer) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	klog.Info("start NodeUnpublishVolume function")

	return &csi.NodeUnpublishVolumeResponse{}, nil
}

// capabilities 中有 NodeServiceCapability_RPC_GET_VOLUME_STATS 时才需要实现此方法
func (cns *CSINodeServer) NodeGetVolumeStats(ctx context.Context, req *csi.NodeGetVolumeStatsRequest) (*csi.NodeGetVolumeStatsResponse, error) {
	klog.Info("start NodeGetVolumeStats function")

	return &csi.NodeGetVolumeStatsResponse{}, nil
}

// capabilities 中有 NodeServiceCapability_RPC_EXPAND_VOLUME 时才需要实现此方法
func (cns *CSINodeServer) NodeExpandVolume(ctx context.Context, req *csi.NodeExpandVolumeRequest) (*csi.NodeExpandVolumeResponse, error) {
	klog.Info("start NodeExpandVolume function")

	return nil, nil
}

func (cns *CSINodeServer) NodeGetCapabilities(ctx context.Context, req *csi.NodeGetCapabilitiesRequest) (*csi.NodeGetCapabilitiesResponse, error) {
	klog.Info("start NodeGetCapabilities function")

	return &csi.NodeGetCapabilitiesResponse{
		Capabilities: cns.capabilities,
	}, nil
}

func (cns *CSINodeServer) NodeGetInfo(ctx context.Context, req *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
	klog.Info("start NodeGetInfo function")

	return &csi.NodeGetInfoResponse{
		NodeId: cns.driver.config.NodeID,
	}, nil
}
