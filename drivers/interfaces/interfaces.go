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
	Write(string, string) error
	Read(string) (string, error)
	ListAllFiles() (map[string]string, error)
	Delete(string) error
}

type ILDB interface {
	Read(key string, i interface{}) error
	Write(key string, i interface{}) error
	Remove(key string) error
}

type IContainers interface {
	PullImage(i Image) error
	BuildImage(opts BuildImageOptions) error

	StartContainer(*Container) error
	StopContainer(*Container) error
	RestartContainer(*Container) error
	RemoveContainer(*Container) error

	ListImages() (map[string]Image, error)
	ListContainers() (map[string]Container, error)

	InspectContainers(c *Container) ([]int64, error)
}

type IPrint interface {
	SetDebug(bool)
	Info(...interface{})
	WhiteInfo(...interface{})
	Infof(string, ...interface{})
	Debugf(string, ...interface{})
	Debug(...interface{})
	Warnf(string, ...interface{})
	Warn(...interface{})
	Errorf(string, ...interface{})
	Error(...interface{})
}

var ErrBucketNotFound error = errors.New("BUCKET_NOT_FOUND")
