package k8s

import (
	account "github.com/lastbackend/lastbackend/libs/adapter/k8s/api/v1"
)

type AccountsGetter interface {
	Accounts() AccountInterface
}

// ComponentStatusInterface has methods to work with ComponentStatus resources.
type AccountInterface interface {
	Create(*account.Account) (*account.Account, error)
}

// componentStatuses implements ComponentStatusInterface
type Accounts struct {
	c *LBClient
}

// newComponentStatuses returns a ComponentStatuses
func newAccounts(c *LBClient) *Accounts {
	return &Accounts{
		c: c,
	}
}

// Create takes the representation of a componentStatus and creates it.  Returns the server's representation of the componentStatus, and an error, if there is any.
func (c *Accounts) Create(Account *account.Account) (result *account.Account, err error) {
	result = &account.Account{}
	err = c.c.Post().
		Resource("accounts").
		Body(Account).
		Do().
		Into(result)
	return
}
