package browser

import "testing"

func TestOpen(t *testing.T) {
	tmp := CommandWrapper
	CommandWrapper = func(name string, parameters ...string) error {
		return nil
	}

	err := Open("http://dummmy")

	if err != nil {
		t.Error("Unexpected error")
	}

	CommandWrapper = tmp
}

func TestOpenFail(t *testing.T) {
	tmp := Os
	Os = "Dummy"
	err := Open("http://dummmy")

	if err == nil {
		t.Error("Unexpected successfully url call")
	}

	Os = tmp
}
