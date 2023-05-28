package lvm

import (
	"fmt"
	"testing"
)

var (
	testLV = &LogicalVolume{
		Path:   "/dev/lvmvg/test",
		Name:   "test",
		vgName: "lvmvg",
		Size:   "5Gi",
	}
)

func TestCreateLogicalVolume(t *testing.T) {
	createLogicalVolume(testLV)
	exist, _ := checkVolumeExists(testLV)
	fmt.Println(exist)
}

func TestRemoveLogicalVolume(t *testing.T) {
	removeLogicalVolume(testLV)
	exist, _ := checkVolumeExists(testLV)
	fmt.Println(exist)
}

func TestCheckVolumeExists(t *testing.T) {
	exist, _ := checkVolumeExists(testLV)
	fmt.Println(exist)
}
