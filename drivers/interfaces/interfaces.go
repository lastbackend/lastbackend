package interfaces

type Log interface {
	Debug(...interface{})
	Error(...interface{})
	SetDebugLevel()
}

type DB interface {
	Write([]byte, []byte) error
	Read([]byte) error
}
