package errors

import "fmt"

func Process(err int) bool {
	switch err {
	case 200 : fmt.Println("OK"); return false
	case 400 : fmt.Println("incorrect json")
	case	401 : fmt.Println("access denied")
	case	406 : fmt.Println("bad parameter")
	case	500 : fmt.Println("server error")
	}
	return true
}
