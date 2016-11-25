package user

import (
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/client/context"
	"errors"
)

func LogoutCmd() {

	ctx := context.Get()

	err := Logout()
	if err != nil {
		ctx.Log.Error(err)
		return
	}
}

func Logout() error {

	ctx := context.Get()

	err := ctx.Session.Clear()
	if err != nil {
		return errors.New("Some problems with logout")
	}

	fmt.Println("Logout successfully")

	return nil
}
