package lvm

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/caoyingjunz/pixiulib/exec"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/houwenchen/kubernetes-csi/pkg/config"
	"k8s.io/klog/v2"
)

/*
使用 linux 的 lvm2
-----------------
名词解释：
物理存储介质（The physical media）：LVM存储介质，可以是硬盘分区、整个硬盘、raid阵列或SAN硬盘。设备必须初始化为LVM物理卷，才能与LVM结合使用。
物理卷PV（physical volume）：物理卷就是LVM的基本存储逻辑块，但和基本的物理存储介质比较却包含与LVM相关的管理参数，创建物理卷可以用硬盘分区，也可以用硬盘本身。
卷组VG（Volume Group）：LVM卷组类似于非LVM系统中的物理硬盘，一个卷组VG由一个或多个物理卷PV组成。可以在卷组VG上建立逻辑卷LV。
逻辑卷LV（logical volume）：类似于非LVM系统中的硬盘分区，逻辑卷LV建立在卷组VG之上。在逻辑卷LV之上建立文件系统。
物理块PE（physical Extent）：物理卷PV中可以分配的最小存储单元，PE的大小可以指定，默认为4MB
逻辑块LE（Logical Extent）：逻辑卷LV中可以分配的最小存储单元，在同一卷组VG中LE的大小和PE是相同的，并且一一相对。
------------------
常用命令：
1. 查看主机可用块设备
root@master:~/tmp/openebs# lsblk
NAME                                                  MAJ:MIN RM   SIZE RO TYPE MOUNTPOINTS
loop0                                                   7:0    0  59.1M  1 loop /snap/core20/1883
loop1                                                   7:1    0  46.4M  1 loop /snap/snapd/19127
loop2                                                   7:2    0   334M  1 loop /snap/gnome-3-38-2004/141
loop3                                                   7:3    0   334M  1 loop /snap/gnome-3-38-2004/138
loop4                                                   7:4    0  46.4M  1 loop /snap/snapd/18940
loop5                                                   7:5    0  59.3M  1 loop /snap/core20/1895
loop6                                                   7:6    0 217.7M  1 loop /snap/firefox/2708
loop7                                                   7:7    0 217.7M  1 loop /snap/firefox/2670
loop8                                                   7:8    0     4K  1 loop /snap/bare/5
loop9                                                   7:9    0  91.7M  1 loop /snap/gtk-common-themes/1535
loop10                                                  7:10   0    10G  0 loop
└─lvmvg-pvc--82de61c0--8097--4c62--966e--ae9d472af74c 253:0    0     4G  0 lvm
sda                                                     8:0    0    64G  0 disk
├─sda1                                                  8:1    0     1G  0 part /boot/efi
└─sda2                                                  8:2    0  62.9G  0 part /var/snap/firefox/common/host-hunspell
                                                                                /
sr0                                                    11:0    1  1024M  0 rom

2. 使用 MOUNTPOINTS 为空的设备创建 pv
root@master:~/tmp/openebs# pvcreate /dev/loop10
root@master:~/tmp/openebs# pvs
  PV          VG    Fmt  Attr PSize   PFree
  /dev/loop10 lvmvg lvm2 a--  <10.00g <6.00g

3. 创建 vg
root@master:~/tmp/openebs# vgcreate lvmvg /dev/loop10
root@master:~/tmp/openebs# vgs
  VG    #PV #LV #SN Attr   VSize   VFree
  lvmvg   1   0   0 wz--n- <10.00g <10.00g
*/

// lvm 模块所需命令
const (
	lvCreate string = "lvcreate"
	lvRemove string = "lvremove"
)

type LogicalVolume struct {
	Path     string
	Name     string
	VGName   string
	UUID     string
	LVAccess []string
	LVStatus string
	Size     string
}

// 根据 CreateVolumeRequest 生成 LV
func NewLogicalVolumeForCreate(config *config.Config, req *csi.CreateVolumeRequest) (*LogicalVolume, error) {
	name := req.GetName()
	paras := req.GetParameters()
	vgname, ok := paras["vgname"]
	if !ok {
		klog.Info("create volume request sholud contain para of vgname")
		return nil, errors.New("miss vgname")
	}

	// TODO: 优化 size 的获取方式
	size := req.GetCapacityRange().RequiredBytes

	path := filepath.Join(config.VolumeDir, vgname, name)

	return &LogicalVolume{
		Path:   path,
		Name:   name,
		VGName: vgname,
		Size:   fmt.Sprint(size),
	}, nil
}

// 根据 DeleteVolumeRequest 生成 LV
func NewLogicalVolumeForDelete(req *csi.DeleteVolumeRequest) (*LogicalVolume, error) {
	name := req.GetVolumeId()
	lv, err := lvsSet.getLVByName(name)
	if err != nil {
		return nil, err
	}
	return lv, nil
}

// lvcreate -n test -L 5Gi lvmvg
func CreateLogicalVolume(lv *LogicalVolume) error {
	// 构造 lvcreate 的命令
	var createLVArg []string

	if len(lv.Name) == 0 || len(lv.VGName) == 0 {
		klog.Info("lvname and vgname can't be empty")
		return errors.New("miss lvname or vgname")
	}

	// TODO: 优化 size 的校验方式
	if len(lv.Size) == 0 {
		klog.Info("lvsize can't be empty")
		return errors.New("miss lvsize")
	}

	// lv 是否存在检查
	exist, err := CheckVolumeExists(lv)
	if err != nil {
		return err
	}

	if exist {
		return errors.New("lv already exists")
	}

	createLVArg = append(createLVArg, "-n", lv.Name)
	createLVArg = append(createLVArg, "-L", lv.Size)
	createLVArg = append(createLVArg, lv.VGName)

	exec := exec.New()
	out, err := exec.Command(lvCreate, createLVArg...).CombinedOutput()
	if err != nil {
		klog.Infof("lvcreate failed, lvname: %s, vgname: %s, size: %v\n", lv.Name, lv.VGName, lv.Size)
		return errors.New("lvcreate failed")
	}

	// 数据归档入 lvsSet
	lvsSet.addLV(lv)

	klog.Info(string(out))
	return nil
}

func CheckVolumeExists(lv *LogicalVolume) (bool, error) {
	if _, err := os.Stat(lv.Path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// lvremove /dev/lvmvg/test -f
func RemoveLogicalVolume(lv *LogicalVolume) error {
	// 构造 lvremove 的命令
	var removeLVArg []string

	// 删除前预检查
	// TODO: 检查 lv 是否还是被 mount 的，可能存在误删除的情况，这里进行维护
	if len(lv.Name) == 0 {
		klog.Info("lvname can't be empty")
		return errors.New("miss lvname")
	}

	// lv 是否存在检查
	exist, err := CheckVolumeExists(lv)
	if err != nil {
		return err
	}

	if !exist {
		return errors.New("lv doesn't exists")
	}

	removeLVArg = append(removeLVArg, lv.Path)
	removeLVArg = append(removeLVArg, "-f")

	exec := exec.New()
	out, err := exec.Command(lvRemove, removeLVArg...).CombinedOutput()
	if err != nil {
		klog.Infof("lvremove failed, lvname: %s, vgname: %s, size: %v\n", lv.Name, lv.VGName, lv.Size)
		return errors.New("lvremove failed")
	}

	// 数据归档入 lvsSet
	lvsSet.deleteLV(lv)

	klog.Info(string(out))
	return nil
}
