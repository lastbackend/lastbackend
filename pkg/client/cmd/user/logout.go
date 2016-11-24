package user

import (
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/client/context"
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
		return err
	}

	fmt.Println("Logout successfully")

	return nil
}
