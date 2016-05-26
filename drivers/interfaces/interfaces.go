package interfaces

import "errors"

type Log interface {
	Debug(...interface{})
	Error(...interface{})
	Fatal(...interface{})
	SetDebugLevel()
}

type DB interface {
	Write(Log, string, string) error
	Read(Log, string) (string, error)
	ListAllFiles(Log) ([]string, error)
	Delete(Log, string) error
}

type Env struct {
	Log      Log
	Database DB
}

var ErrBucketNotFound error = errors.New("BUCKET_NOT_FOUND")
