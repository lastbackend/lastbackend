package cmd

import (
	"fmt"
	"github.com/lastbackend/lastbackend/cmd/client/context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuthMock(t *testing.T) {
	ctx := context.Mock()

	expected, err := Login(ctx)
	if err != nil {
		fmt.Println(err.Error())
	}

	assert.Equal(t, expected, "token")
}
