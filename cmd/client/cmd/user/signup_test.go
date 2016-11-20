package cmd

import (
	"fmt"
	"github.com/lastbackend/lastbackend/cmd/client/context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSignUpMock(t *testing.T) {
	ctx := context.Mock()

	expected, err := CreateNewUser(ctx)
	if err != nil {
		fmt.Println(err.Error())
	}

	assert.Equal(t, expected, "token")
}
