package cmd

import (
	"fmt"
	"github.com/lastbackend/lastbackend/cmd/client/context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWhoamiMock(t *testing.T) {
	ctx := context.Mock()

	actual, err := WhoamiLogic(ctx)
	if err != nil {
		fmt.Println(err.Error())
	}

	assert.Equal(t, "some_id", actual.Id)
	assert.Equal(t, "some_username", actual.Username)
	assert.Equal(t, "some_email", actual.Email)
	assert.Equal(t, "some_gravatar", actual.Gravatar)
	assert.Equal(t, float32(10), actual.Balance)
	assert.Equal(t, false, actual.Organization)
	assert.Equal(t, "some_first_name", actual.Profile.FirstName)
	assert.Equal(t, "some_last_name", actual.Profile.LastName)
	assert.Equal(t, "some_company", actual.Profile.Company)
	assert.Equal(t, "2014-01-16T07:38:28.45Z", actual.Created)
	assert.Equal(t, "2014-01-16T07:38:28.45Z", actual.Updated)
}
