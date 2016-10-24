package k8s

import (
	build "github.com/lastbackend/lastbackend/libs/adapter/k8s/api/v1"
)

type ComponentBuildsGetter interface {
	Builds() ComponentBuildInterface
}

// ComponentStatusInterface has methods to work with ComponentStatus resources.
type ComponentBuildInterface interface {
	Create(*build.Build) (*build.Build, error)
}

// componentStatuses implements ComponentStatusInterface
type componentBuilds struct {
	client *LBClient
}

// newComponentStatuses returns a ComponentStatuses
func newComponentBuilds(c *LBClient) *componentBuilds {
	return &componentBuilds{
		client: c,
	}
}

// Create takes the representation of a componentStatus and creates it.  Returns the server's representation of the componentStatus, and an error, if there is any.
func (c *componentBuilds) Create(Build *build.Build) (result *build.Build, err error) {
	result = &build.Build{}
	err = c.client.Post().
		Resource("builds").
		Body(Build).
		Do().
		Into(result)
	return
}
