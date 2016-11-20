package cmd

import (
	"fmt"
	"github.com/lastbackend/lastbackend/cmd/client/context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuthMock(t *testing.T) {
	ctx := context.Mock()

	ctx.Info.Version = "OK"
	expected, err, _ := Login(ctx)
	if err != nil {
		fmt.Println(err.Error())
	}

	assert.Equal(t, expected, "token")

	ctx.Info.Version = "BAD"
	_, err, httpError := Login(ctx)
	if err != nil {
		fmt.Println(err.Error())
	}

	assert.Equal(t, httpError, "access denied")
}
