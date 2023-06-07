package lvm

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"syscall"

	"github.com/houwenchen/kubernetes-csi/pkg/helper"
)

type Request struct {
	DeviceName       string
	PodUid           string
	ContainerRuntime string
	IOLimit          *IOMax
}

type ValidRequest struct {
	FilePath     string
	DeviceNumber *DeviceNumber
	IOMax        *IOMax
}

type IOMax struct {
	Riops uint64
	Wiops uint64
	Rbps  uint64
	Wbps  uint64
}

type DeviceNumber struct {
	Major uint64
	Minor uint64
}

const (
	baseCgroupPath = "/sys/fs/cgroup"
)

// SetIOLimits sets iops, bps limits for a pod with uid podUid for accessing a device named deviceName
// provided that the underlying cgroup used for pod namespacing is cgroup2 (cgroup v2)
func SetIOLimits(request *Request) error {
	if !helper.DirExists(baseCgroupPath) {
		return errors.New(baseCgroupPath + " does not exist")
	}
	if err := checkCgroupV2(); err != nil {
		return err
	}
	validRequest, err := validate(request)
	if err != nil {
		return err
	}
	err = setIOLimits(validRequest)
	return err
}

func validate(request *Request) (*ValidRequest, error) {
	if !helper.IsValidUUID(request.PodUid) {
		return nil, errors.New("Expected PodUid in UUID format, Got " + request.PodUid)
	}
	podCGPath, err := getPodCGroupPath(request.PodUid, request.ContainerRuntime)
	if err != nil {
		return nil, err
	}
	ioMaxFile := podCGPath + "/io.max"
	if !helper.FileExists(ioMaxFile) {
		return nil, errors.New("io.max file is not present in pod CGroup")
	}
	deviceNumber, err := getDeviceNumber(request.DeviceName)
	if err != nil {
		return nil, errors.New("device major:minor numbers could not be obtained")
	}
	return &ValidRequest{
		FilePath:     ioMaxFile,
		DeviceNumber: deviceNumber,
		IOMax:        request.IOLimit,
	}, nil
}

func getPodCGroupPath(podUid string, cruntime string) (string, error) {
	switch cruntime {
	case "containerd":
		path, err := getContainerdCGPath(podUid)
		if err != nil {
			return "", err
		}
		return path, nil
	default:
		return "", errors.New(cruntime + " runtime support is not present")
	}

}

func checkCgroupV2() error {
	if !helper.FileExists(baseCgroupPath + "/cgroup.controllers") {
		return errors.New("CGroupV2 not enabled")
	}
	return nil
}

func getContainerdPodCGSuffix(podUid string) string {
	return "pod" + strings.ReplaceAll(podUid, "-", "_")
}

func getContainerdCGPath(podUid string) (string, error) {
	kubepodsCGPath := baseCgroupPath + "/kubepods.slice"
	podSuffix := getContainerdPodCGSuffix(podUid)
	podCGPath := kubepodsCGPath + "/kubepods-" + podSuffix + ".slice"
	if helper.DirExists(podCGPath) {
		return podCGPath, nil
	}
	podCGPath = kubepodsCGPath + "/kubepods-besteffort.slice/kubepods-besteffort-" + podSuffix + ".slice"
	if helper.DirExists(podCGPath) {
		return podCGPath, nil
	}
	podCGPath = kubepodsCGPath + "/kubepods-burstable.slice/kubepods-burstable-" + podSuffix + ".slice"
	if helper.DirExists(podCGPath) {
		return podCGPath, nil
	}
	return "", errors.New("CGroup Path not found for pod with Uid: " + podUid)
}

func getDeviceNumber(deviceName string) (*DeviceNumber, error) {
	stat := syscall.Stat_t{}
	if err := syscall.Stat(deviceName, &stat); err != nil {
		return nil, err
	}
	return &DeviceNumber{
		Major: uint64(stat.Rdev / 256),
		Minor: uint64(stat.Rdev % 256),
	}, nil
}

func getIOLimitsStr(deviceNumber *DeviceNumber, ioMax *IOMax) string {
	line := strconv.FormatUint(deviceNumber.Major, 10) + ":" + strconv.FormatUint(deviceNumber.Minor, 10)
	if ioMax.Riops != 0 {
		line += " riops=" + strconv.FormatUint(ioMax.Riops, 10)
	}
	if ioMax.Wiops != 0 {
		line += " wiops=" + strconv.FormatUint(ioMax.Wiops, 10)
	}
	if ioMax.Rbps != 0 {
		line += " rbps=" + strconv.FormatUint(ioMax.Rbps, 10)
	}
	if ioMax.Wbps != 0 {
		line += " wbps=" + strconv.FormatUint(ioMax.Wbps, 10)
	}
	return line
}

func setIOLimits(request *ValidRequest) error {
	line := getIOLimitsStr(request.DeviceNumber, request.IOMax)
	err := os.WriteFile(request.FilePath, []byte(line), 0600)
	return err
}
