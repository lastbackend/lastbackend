package k8s

import (
	build "github.com/lastbackend/lastbackend/libs/adapter/k8s/api/v1"
	"k8s.io/client-go/1.5/pkg/api"
	"k8s.io/client-go/1.5/pkg/api/v1"
	"k8s.io/client-go/1.5/pkg/watch"
)

type BuildsGetter interface {
	Builds() BuildInterface
}

// Build has methods to work with Builds resources.
type BuildInterface interface {
	List(string, v1.ListOptions) (*build.BuildList, error)
	Get(string, string) (*build.Build, error)
	Create(string, *build.Build) (*build.Build, error)
	Update(string, *build.Build) (*build.Build, error)
	Delete(string, string) error
	Watch(string, v1.ListOptions) (watch.Interface, error)
}

// Build implements Build Interface
type builds struct {
	c *LBClient
}

// newBuild returns a new build interface
func newBuilds(c *LBClient) *builds {
	return &builds{
		c: c,
	}
}

// List returns a list of builds that match query
func (c *builds) List(namespace string, opts v1.ListOptions) (result *build.BuildList, err error) {
	result = &build.BuildList{}
	err = c.c.Get().
		Namespace(namespace).
		Resource("builds").
		VersionedParams(&opts, api.ParameterCodec).
		Do().
		Into(result)
	return
}

// Get returns information about a particular build
func (c *builds) Get(namespace, name string) (result *build.Build, err error) {
	result = new(build.Build)
	err = c.c.Get().Namespace(namespace).Resource("builds").Name(name).Do().Into(result)
	return
}

// Create creates new build. Returns create build
func (c *builds) Create(namespace string, spec *build.Build) (result *build.Build, err error) {
	result = new(build.Build)
	err = c.c.Post().Namespace(namespace).Resource("builds").Body(spec).Do().Into(result)
	return
}

// Update updates the build on server. Returns create build
func (c *builds) Update(namespace string, spec *build.Build) (result *build.Build, err error) {
	result = new(build.Build)
	err = c.c.Put().Namespace(namespace).Resource("builds").Name(spec.Name).Body(spec).Do().Into(result)
	return
}

// Delete deletes a build
func (c *builds) Delete(namespace, name string) (err error) {
	err = c.c.Delete().Namespace(namespace).Resource("builds").Name(name).Do().Error()
	return
}

// Watch returns a watch.Interface that watches the requested builds changes
func (c *builds) Watch(namespace string, opts v1.ListOptions) (watch.Interface, error) {
	return c.c.Get().
		Prefix("watch").
		Namespace(namespace).
		Resource("builds").
		VersionedParams(&opts, api.ParameterCodec).
		Watch()
}
