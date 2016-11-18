package cmd

import (
	"github.com/lastbackend/lastbackend/cmd/client/context"
	"testing"
)

func TestAuthMock(t *testing.T) {
	ctx := context.Mock()

	Auth(ctx)
}
