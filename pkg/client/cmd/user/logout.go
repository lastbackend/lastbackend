package user

import (
	"errors"
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

	fmt.Println("LOGOUT.GO CTX", ctx)

	fmt.Println("LOGOUT.GO CTX.STORAGE", ctx.Storage)

	err := ctx.Storage.Clear()

	if err != nil {
		return errors.New("Some problems with logout")
	}

	fmt.Println("Logout successfully")

	return nil
}
