package store

type Book struct {
	Id string `json:"id"`
}

// 对存储的简单操作
type Store interface {
	Create(*Book) error
	Update(*Book) error
	Get(string) (Book, error)
	GetAll() ([]Book, error)
	Delete(string) error
}
