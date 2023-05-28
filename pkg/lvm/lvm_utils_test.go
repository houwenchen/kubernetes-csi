package lvm

import (
	"fmt"
	"testing"
)

var (
	testLV = &LogicalVolume{
		Path:   "/dev/lvmvg/test",
		Name:   "test",
		VGName: "lvmvg",
		Size:   "5Gi",
	}
)

func TestCreateLogicalVolume(t *testing.T) {
	CreateLogicalVolume(testLV)
	exist, _ := CheckVolumeExists(testLV)
	fmt.Println(exist)
}

func TestRemoveLogicalVolume(t *testing.T) {
	RemoveLogicalVolume(testLV)
	exist, _ := CheckVolumeExists(testLV)
	fmt.Println(exist)
}

func TestCheckVolumeExists(t *testing.T) {
	exist, _ := CheckVolumeExists(testLV)
	fmt.Println(exist)
}
