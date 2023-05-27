package lvm

import (
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
}

func TestRemoveLogicalVolume(t *testing.T) {
	removeLogicalVolume(testLV)
}
