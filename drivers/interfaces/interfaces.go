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
}
