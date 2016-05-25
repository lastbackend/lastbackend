package interfaces

type Log interface {
	Debug(...interface{})
	Error(...interface{})
	Fatal(...interface{})
	SetDebugLevel()
}

type DB interface {
	Write(Log, []byte, []byte) error
	Read(Log, []byte) (string, error)
	ListAllFiles(Log) ([]string, error)
	Delete(Log, []byte) error
}

type Env struct {
	Log      Log
	Database DB
}
