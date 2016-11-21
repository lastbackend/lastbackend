package http



type IHTTP interface {
	Post(string, []byte, string, string) ([]byte, int)
	Get(string, []byte, string, string) ([]byte, int)
	Delete(string, string, string) int
}
