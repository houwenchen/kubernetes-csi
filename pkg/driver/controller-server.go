package hostpath

import (
	"context"
	"errors"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/google/uuid"
	"k8s.io/klog/v2"
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
	// CSIControllerServer 的默认能力集
	defaultControllerServiceCapability_RPC_Types = []csi.ControllerServiceCapability_RPC_Type{
		csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
		csi.ControllerServiceCapability_RPC_GET_VOLUME,
		csi.ControllerServiceCapability_RPC_GET_CAPACITY,
		csi.ControllerServiceCapability_RPC_LIST_VOLUMES,
		csi.ControllerServiceCapability_RPC_VOLUME_CONDITION,
		csi.ControllerServiceCapability_RPC_SINGLE_NODE_MULTI_WRITER,
		csi.ControllerServiceCapability_RPC_PUBLISH_UNPUBLISH_VOLUME,
	}

	// CSIControllerServer volume 的能力集
	defaultVolumeCaps = []csi.VolumeCapability_AccessMode{
		{
			Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
		},
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
	klog.Info("start create volume function")

	// 首先对创建 volume 的请求进行判断
	err := ccs.validateCreateVolumeRequest(req)
	if err != nil {
		return nil, err
	}

	// 生成唯一标识 volume 的 VolumeId
	volumeId := uuid.New().String()

	// 生成 VolumeContext
	volumeContext := make(map[string]string)
	volumeContext["driver-name"] = ccs.driver.config.DriverName
	volumeContext["volume-name"] = req.GetName()

	// TODO: 增加 create volume 逻辑

	return &csi.CreateVolumeResponse{
		Volume: &csi.Volume{
			CapacityBytes: req.GetCapacityRange().RequiredBytes,
			VolumeId:      volumeId,
			VolumeContext: volumeContext,
			ContentSource: req.GetVolumeContentSource(),
		},
	}, nil
}

// 对 CreateVolumeRequest 的必选字段进行校验
func (ccs *CSIControllerServer) validateCreateVolumeRequest(req *csi.CreateVolumeRequest) error {
	klog.Info("start validate create volume request")

	// CreateVolumeRequest----Name 字段检查
	// 1. 保证幂等性
	// 2. 特殊需要时可以用这个字段来作为标识字段
	if len(req.Name) == 0 {
		return errors.New("volume's Name is required")
	}

	// CreateVolumeRequest----VolumeCapabilities 字段检查
	reqCaps := req.GetVolumeCapabilities()
	if len(reqCaps) == 0 {
		return errors.New("volume's capability is requored")
	}
	if !ccs.validateVolumeCapabilitiesOfReq(reqCaps) {
		return errors.New("unsupport VolumeCapability")
	}

	return nil
}

func (ccs *CSIControllerServer) validateVolumeCapabilitiesOfReq(caps []*csi.VolumeCapability) bool {
	for _, c := range caps {
		found := false
		for _, dc := range defaultVolumeCaps {
			if dc.GetMode() == c.AccessMode.GetMode() {
				found = true
				continue
			}
		}
		if !found {
			return false
		}
	}

	return true
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
