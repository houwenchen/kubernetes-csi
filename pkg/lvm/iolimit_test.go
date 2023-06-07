package lvm

import (
	"fmt"
	"testing"
)

func TestTest(t *testing.T) {
	deviceNumber := &DeviceNumber{
		Major: 152,
		Minor: 2,
	}

	ioInfos := &IOMax{
		Riops: 1024 * 1024,
		Rbps:  1024 * 1024,
		Wiops: 1024 * 1024,
		Wbps:  1024 * 1024,
	}

	str := getIOLimitsStr(deviceNumber, ioInfos)
	fmt.Print(str)
}
