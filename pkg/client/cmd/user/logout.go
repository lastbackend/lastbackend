package user

import (
	"errors"
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

	var (
		err error
		ctx = context.Get()
	)

	err = ctx.Storage.Clear()
	if err != nil {
		return errors.New("Some problems with logout")
	}

	ctx.Log.Info("Logout successfully")

	return nil
}
