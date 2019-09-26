package utils

import (
	"github.com/antlinker/go-dirtyfilter"
	"github.com/antlinker/go-dirtyfilter/store"
)

var (
	filterManage *filter.DirtyManager
)

func InitDirtyManager(wds []string) (err error) {
	memStore, err := store.NewMemoryStore(store.MemoryConfig{
		DataSource: wds,
	})
	if err != nil {
		return
	}
	filterManage = filter.NewDirtyManager(memStore)

	return
}

func IsForbidden(content string) bool {
	if filterManage == nil {
		return false
	}
	d, err := filterManage.Filter().Filter(content, ' ')
	if err != nil {
		return false
	}
	return len(d) > 0
}
