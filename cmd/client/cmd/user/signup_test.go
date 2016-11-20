package cmd

import (
	"fmt"
	"github.com/lastbackend/lastbackend/cmd/client/config"
	"github.com/lastbackend/lastbackend/cmd/client/context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSignUpMock(t *testing.T) {
	ctx := context.Mock()
	cfg := config.Get()

	ctx.Info.Version = "OK"
	expected, err, _ := CreateNewUser(ctx, cfg)
	if err != nil {
		fmt.Println(err.Error())
	}

	assert.Equal(t, expected, "token")

	ctx.Info.Version = "BAD_USERNAME"
	_, err, httpError := CreateNewUser(ctx, cfg)
	if err != nil {
		fmt.Println(err.Error())
	}

	assert.Equal(t, httpError, "bad username parameter")

	ctx.Info.Version = "BAD_EMAIL"
	_, err, httpError = CreateNewUser(ctx, cfg)
	if err != nil {
		fmt.Println(err.Error())
	}

	assert.Equal(t, httpError, "bad email parameter")

	ctx.Info.Version = "BAD_PASSWORD"
	_, err, httpError = CreateNewUser(ctx, cfg)
	if err != nil {
		fmt.Println(err.Error())
	}

	assert.Equal(t, httpError, "bad password parameter")
}
