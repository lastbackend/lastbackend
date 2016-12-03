package db

type IDB interface {
	Get(string, interface{}) error
	Set(string, interface{}) error
	Clear() error
	Close() error
}
