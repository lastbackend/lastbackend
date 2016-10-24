package k8s

import (
	account "github.com/lastbackend/lastbackend/libs/adapter/k8s/api/v1"
)

type ComponentAccountsGetter interface {
	Accounts() ComponentAccountInterface
}

// ComponentStatusInterface has methods to work with ComponentStatus resources.
type ComponentAccountInterface interface {
	Create(*account.Account) (*account.Account, error)
}

// componentStatuses implements ComponentStatusInterface
type componentAccounts struct {
	client *LBClient
}

// newComponentStatuses returns a ComponentStatuses
func newComponentAccounts(c *LBClient) *componentAccounts {
	return &componentAccounts{
		client: c,
	}
}

// Create takes the representation of a componentStatus and creates it.  Returns the server's representation of the componentStatus, and an error, if there is any.
func (c *componentAccounts) Create(Account *account.Account) (result *account.Account, err error) {
	result = &account.Account{}
	err = c.client.Post().
		Resource("accounts").
		Body(Account).
		Do().
		Into(result)
	return
}
