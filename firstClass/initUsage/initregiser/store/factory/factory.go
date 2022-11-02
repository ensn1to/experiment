package factory

import (
	"fmt"
	"sync"

	"github.com/ensn1to/experiment/tree/master/firstClass/init/initregiser/store"
)

var (
	providerMux sync.Mutex
	// 支持存储方式列表
	providers = make(map[string]store.Store)
)

// 注册存储方式
func Register(name string, s store.Store) {
	providerMux.Lock()
	defer providerMux.Unlock()
	if s == nil {
		panic("store: Register provider is nil")
	}

	if _, ok := providers[name]; ok {
		panic("store: Register provider is exists")
	}

	providers[name] = s
}

// 获取使用
func New(providerName string) (store.Store, error) {
	providerMux.Lock()
	defer providerMux.Unlock()

	if provider, ok := providers[providerName]; ok {
		return provider, nil
	}

	return nil, fmt.Errorf("store: provider %s is not registered", providerName)
}
