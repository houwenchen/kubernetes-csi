package lvm

import (
	"errors"
	"sync"

	"k8s.io/klog/v2"
)

type LVsMap struct {
	mu  sync.RWMutex
	lvs map[string]*LogicalVolume
}

var (
	lvsSet = &LVsMap{
		lvs: make(map[string]*LogicalVolume),
	}
)

func (lvsmap *LVsMap) addLV(lv *LogicalVolume) error {
	lvsmap.mu.Lock()
	defer lvsmap.mu.Unlock()

	if lvsmap.checkLVExists(lv) {
		klog.Infof("lv already exists, lvname: %s\n", lv.Name)
		return errors.New("lv already exists")
	}

	lvsmap.lvs[lv.Name] = lv

	return nil
}

func (lvsmap *LVsMap) checkLVExists(lv *LogicalVolume) bool {
	lvsmap.mu.RLock()
	defer lvsmap.mu.RUnlock()

	_, exist := lvsmap.lvs[lv.Name]

	return exist
}

func (lvsmap *LVsMap) deleteLV(lv *LogicalVolume) error {
	lvsmap.mu.Lock()
	defer lvsmap.mu.Unlock()

	if !lvsmap.checkLVExists(lv) {
		klog.Infof("lv doesn't exists, lvname: %s\n", lv.Name)
		return errors.New("lv doesn't exists")
	}

	delete(lvsmap.lvs, lv.Name)

	return nil
}

func (lvsmap *LVsMap) getLVByName(name string) (*LogicalVolume, error) {
	lvsmap.mu.RLock()
	defer lvsmap.mu.RUnlock()

	lv, exists := lvsmap.lvs[name]
	if !exists {
		return nil, errors.New("lvsSet doesn't contain this lv")
	}

	return lv, nil
}
