package localdb

type ILocalStorage interface {
	Get(string, interface{}) error
	Set(string, interface{}) error
	Clear() error
	Init() error
	Close() error
}
