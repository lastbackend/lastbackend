package interfaces

import "errors"

type ILog interface {
	Debug(...interface{})
	Debugf(string, ...interface{})
	Info(...interface{})
	Infof(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
	Fatal(...interface{})
	Fatalf(string, ...interface{})
	SetDebugLevel()
}

type IStorage interface {
	Write(ILog, string, string) error
	Read(ILog, string) (string, error)
	ListAllFiles(ILog) (map[string]string, error)
	Delete(ILog, string) error
}

type IDB interface {
	Exec(query string, args ...interface{}) (Result, error)
	Query(query string, args ...interface{}) (Rows, error)
	QueryRow(query string, args ...interface{}) Row
}

type Rows interface {
	Scan(dest ...interface{}) error
	Columns() ([]string, error)
	Close() error
	Err() error
	Next() bool
}

type Row interface {
	Scan(dest ...interface{}) error
}

type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}

type IContainers interface {
	GetContainer(string) (Container, error)

	StartContainer(*Container) error
	StopContainer(*Container) error
	RestartContainer(*Container) error
	RemoveContainer(*Container) error

	PullImage(i Image) error
	BuildImage(opts BuildImageOptions) error

	ListImages() (map[string]Image, error)
	ListContainers() (map[string]Container, error)

	System() (*Node, error)
}

var ErrBucketNotFound error = errors.New("BUCKET_NOT_FOUND")
