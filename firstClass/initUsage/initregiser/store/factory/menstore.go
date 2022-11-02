package factory

import (
	"sync"

	"github.com/ensn1to/experiment/tree/master/firstClass/init/initregiser/store"
)

/*
	内存方式存储
*/

// 注册到store.factory
func init() {
	Register("men", &MemStore{})
}

type MemStore struct {
	sync.RWMutex
}

// Create creates a new Book in the store.
func (ms *MemStore) Create(book *store.Book) error {
	return nil
}

// Update updates the existed Book in the store.
func (ms *MemStore) Update(book *store.Book) error {
	ms.Lock()
	defer ms.Unlock()

	return nil
}

// Get retrieves a book from the store, by id. If no such id exists. an
// error is returned.
func (ms *MemStore) Get(id string) (store.Book, error) {
	ms.RLock()
	defer ms.RUnlock()

	return store.Book{}, nil
}

// Delete deletes the book with the given id. If no such id exist. an error
// is returned.
func (ms *MemStore) Delete(id string) error {
	ms.Lock()
	defer ms.Unlock()

	return nil
}

// GetAll returns all the books in the store, in arbitrary order.
func (ms *MemStore) GetAll() ([]store.Book, error) {
	ms.RLock()
	defer ms.RUnlock()

	return nil, nil
}
