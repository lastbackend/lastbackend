package k8s

import (
	user "github.com/lastbackend/lastbackend/libs/adapter/k8s/api/v1"
)

type UsersGetter interface {
	Users() UserInterface
}

// UsersInterface has methods to work with Users resources.
type UserInterface interface {
	Create(*user.User) (*user.User, error)
	Get(name string) (*user.User, error)
}

// Users implements UsersInterface
type users struct {
	c *LBClient
}

// newUsers returns a new user interface
func newUsers(c *LBClient) *users {
	return &users{
		c: c,
	}
}

// Create takes the representation of a users and creates it.
// Returns the server's representation of the users, and an error, if there is any.
func (c *users) Create(User *user.User) (result *user.User, err error) {
	result = &user.User{}
	err = c.c.Post().
		Resource("users").
		Body(User).
		Do().
		Into(result)
	return
}

// Returns the server's representation of the users, and an error, if there is any.
func (c *users) Get(name string) (result *user.User, err error) {
	result = &user.User{}
	err = c.c.Get().
		Resource("users").
		Name(name).
		Do().
		Into(result)
	return
}

// Update updates the user on server. Returns updated user
func (c *users) Update(namespace string, spec *user.User) (result *user.User, err error) {
	result = new(user.User)
	err = c.c.Put().
		Namespace(namespace).
		Resource("users").
		Name(spec.Name).
		Body(spec).
		Do().
		Into(result)
	return
}
