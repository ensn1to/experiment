package store

import (
	"github.com/ensn1to/experiment/tree/master/firstClass/init/initregiser/store"
	"github.com/ensn1to/experiment/tree/master/firstClass/init/initregiser/store/factory"
)

// 常用用法：interface作为参数
func GetBookInfo(id string, s store.Store) (store.Book, error) {
	return s.Get(id)
}

func main() {
	s, _ := factory.New("mem")
	GetBookInfo("aaa", s)
}
