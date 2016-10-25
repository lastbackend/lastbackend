package k8s

import (
	user "github.com/lastbackend/lastbackend/libs/adapter/k8s/api/v1"
)

type UserExpansion interface{}

type ComponentUsersGetter interface {
	Users() ComponentUserInterface
}

// ComponentStatusInterface has methods to work with ComponentStatus resources.
type ComponentUserInterface interface {
	Create(*user.User) (*user.User, error)
	Get(name string) (*user.User, error)
	UserExpansion
}

// componentStatuses implements ComponentStatusInterface
type componentUserResources struct {
	client *LBClient
}

// newComponentStatuses returns a ComponentStatuses
func newComponentUsers(c *LBClient) *componentUserResources {
	return &componentUserResources{
		client: c,
	}
}

// Create takes the representation of a componentStatus and creates it.  Returns the server's representation of the componentStatus, and an error, if there is any.
func (c *componentUserResources) Create(User *user.User) (result *user.User, err error) {
	result = &user.User{}
	err = c.client.Post().
		Resource("users").
		Body(User).
		Do().
		Into(result)
	return
}

func (c *componentUserResources) Get(name string) (result *user.User, err error) {
	result = &user.User{}
	err = c.client.Get().
		Resource("users").
		Name(name).
		Do().
		Into(result)
	return
}
