package hostpath

import (
	"context"

	"github.com/container-storage-interface/spec/lib/go/csi"
)

/*
实现下面的接口
type ControllerServer interface {
    CreateVolume(context.Context, *CreateVolumeRequest) (*CreateVolumeResponse, error)
    DeleteVolume(context.Context, *DeleteVolumeRequest) (*DeleteVolumeResponse, error)
    ControllerPublishVolume(context.Context, *ControllerPublishVolumeRequest) (*ControllerPublishVolumeResponse, error)
    ControllerUnpublishVolume(context.Context, *ControllerUnpublishVolumeRequest) (*ControllerUnpublishVolumeResponse, error)
    ValidateVolumeCapabilities(context.Context, *ValidateVolumeCapabilitiesRequest) (*ValidateVolumeCapabilitiesResponse, error)
    ListVolumes(context.Context, *ListVolumesRequest) (*ListVolumesResponse, error)
    GetCapacity(context.Context, *GetCapacityRequest) (*GetCapacityResponse, error)
    ControllerGetCapabilities(context.Context, *ControllerGetCapabilitiesRequest) (*ControllerGetCapabilitiesResponse, error)
    CreateSnapshot(context.Context, *CreateSnapshotRequest) (*CreateSnapshotResponse, error)
    DeleteSnapshot(context.Context, *DeleteSnapshotRequest) (*DeleteSnapshotResponse, error)
    ListSnapshots(context.Context, *ListSnapshotsRequest) (*ListSnapshotsResponse, error)
    ControllerExpandVolume(context.Context, *ControllerExpandVolumeRequest) (*ControllerExpandVolumeResponse, error)
    ControllerGetVolume(context.Context, *ControllerGetVolumeRequest) (*ControllerGetVolumeResponse, error)
}
*/

var (
	defaultControllerServiceCapability_RPC_Types = []csi.ControllerServiceCapability_RPC_Type{
		csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
		csi.ControllerServiceCapability_RPC_GET_VOLUME,
		csi.ControllerServiceCapability_RPC_GET_CAPACITY,
		csi.ControllerServiceCapability_RPC_LIST_VOLUMES,
		csi.ControllerServiceCapability_RPC_VOLUME_CONDITION,
		csi.ControllerServiceCapability_RPC_SINGLE_NODE_MULTI_WRITER,
		csi.ControllerServiceCapability_RPC_PUBLISH_UNPUBLISH_VOLUME,
	}
)

type CSIControllerServer struct {
	driver       *HostPathDriver
	capabilities []*csi.ControllerServiceCapability
}

func NewDefaultCSIControllerServer(driver *HostPathDriver) *CSIControllerServer {
	capabilities := make([]*csi.ControllerServiceCapability, 0)

	for _, RPCType := range defaultControllerServiceCapability_RPC_Types {
		cap := &csi.ControllerServiceCapability{
			Type: &csi.ControllerServiceCapability_Rpc{
				Rpc: &csi.ControllerServiceCapability_RPC{
					Type: RPCType,
				},
			}}

		capabilities = append(capabilities, cap)
	}

	return NewCSIControllerServerWithOpt(driver, capabilities)
}

func NewCSIControllerServerWithOpt(driver *HostPathDriver, opts ...[]*csi.ControllerServiceCapability) *CSIControllerServer {
	capabilities := make([]*csi.ControllerServiceCapability, 0)

	for _, opt := range opts {
		capabilities = append(capabilities, opt...)
	}

	return &CSIControllerServer{
		driver:       driver,
		capabilities: capabilities,
	}
}

func (ccs *CSIControllerServer) CreateVolume(ctx context.Context, req *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {
	return nil, nil
}

func (ccs *CSIControllerServer) DeleteVolume(ctx context.Context, req *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
	return nil, nil
}

func (ccs *CSIControllerServer) ControllerPublishVolume(ctx context.Context, req *csi.ControllerPublishVolumeRequest) (*csi.ControllerPublishVolumeResponse, error) {
	return nil, nil
}

func (ccs *CSIControllerServer) ControllerUnpublishVolume(ctx context.Context, req *csi.ControllerUnpublishVolumeRequest) (*csi.ControllerUnpublishVolumeResponse, error) {
	return nil, nil
}

func (ccs *CSIControllerServer) ValidateVolumeCapabilities(ctx context.Context, req *csi.ValidateVolumeCapabilitiesRequest) (*csi.ValidateVolumeCapabilitiesResponse, error) {
	return nil, nil
}

func (ccs *CSIControllerServer) ListVolumes(ctx context.Context, req *csi.ListVolumesRequest) (*csi.ListVolumesResponse, error) {
	return nil, nil
}

func (ccs *CSIControllerServer) GetCapacity(ctx context.Context, req *csi.GetCapacityRequest) (*csi.GetCapacityResponse, error) {
	return nil, nil
}

func (ccs *CSIControllerServer) ControllerGetCapabilities(ctx context.Context, req *csi.ControllerGetCapabilitiesRequest) (*csi.ControllerGetCapabilitiesResponse, error) {
	return nil, nil
}

func (ccs *CSIControllerServer) CreateSnapshot(ctx context.Context, req *csi.CreateSnapshotRequest) (*csi.CreateSnapshotResponse, error) {
	return nil, nil
}

func (ccs *CSIControllerServer) DeleteSnapshot(ctx context.Context, req *csi.DeleteSnapshotRequest) (*csi.DeleteSnapshotResponse, error) {
	return nil, nil
}

func (ccs *CSIControllerServer) ListSnapshots(ctx context.Context, req *csi.ListSnapshotsRequest) (*csi.ListSnapshotsResponse, error) {
	return nil, nil
}

func (ccs *CSIControllerServer) ControllerExpandVolume(ctx context.Context, req *csi.ControllerExpandVolumeRequest) (*csi.ControllerExpandVolumeResponse, error) {
	return nil, nil
}

func (ccs *CSIControllerServer) ControllerGetVolume(ctx context.Context, req *csi.ControllerGetVolumeRequest) (*csi.ControllerGetVolumeResponse, error) {
	return nil, nil
}
