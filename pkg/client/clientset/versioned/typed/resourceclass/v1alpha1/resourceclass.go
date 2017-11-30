/*
Copyright 2017 The Caicloud sample-controller Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	v1alpha1 "github.com/caicloud/resource-controller/pkg/apis/resourceclass/v1alpha1"
	scheme "github.com/caicloud/resource-controller/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// ResourceClassesGetter has a method to return a ResourceClassInterface.
// A group's client should implement this interface.
type ResourceClassesGetter interface {
	ResourceClasses(namespace string) ResourceClassInterface
}

// ResourceClassInterface has methods to work with ResourceClass resources.
type ResourceClassInterface interface {
	Create(*v1alpha1.ResourceClass) (*v1alpha1.ResourceClass, error)
	Update(*v1alpha1.ResourceClass) (*v1alpha1.ResourceClass, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.ResourceClass, error)
	List(opts v1.ListOptions) (*v1alpha1.ResourceClassList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.ResourceClass, err error)
	ResourceClassExpansion
}

// resourceClasses implements ResourceClassInterface
type resourceClasses struct {
	client rest.Interface
	ns     string
}

// newResourceClasses returns a ResourceClasses
func newResourceClasses(c *ResourceclassV1alpha1Client, namespace string) *resourceClasses {
	return &resourceClasses{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the resourceClass, and returns the corresponding resourceClass object, and an error if there is any.
func (c *resourceClasses) Get(name string, options v1.GetOptions) (result *v1alpha1.ResourceClass, err error) {
	result = &v1alpha1.ResourceClass{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("resourceclasses").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ResourceClasses that match those selectors.
func (c *resourceClasses) List(opts v1.ListOptions) (result *v1alpha1.ResourceClassList, err error) {
	result = &v1alpha1.ResourceClassList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("resourceclasses").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested resourceClasses.
func (c *resourceClasses) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("resourceclasses").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a resourceClass and creates it.  Returns the server's representation of the resourceClass, and an error, if there is any.
func (c *resourceClasses) Create(resourceClass *v1alpha1.ResourceClass) (result *v1alpha1.ResourceClass, err error) {
	result = &v1alpha1.ResourceClass{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("resourceclasses").
		Body(resourceClass).
		Do().
		Into(result)
	return
}

// Update takes the representation of a resourceClass and updates it. Returns the server's representation of the resourceClass, and an error, if there is any.
func (c *resourceClasses) Update(resourceClass *v1alpha1.ResourceClass) (result *v1alpha1.ResourceClass, err error) {
	result = &v1alpha1.ResourceClass{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("resourceclasses").
		Name(resourceClass.Name).
		Body(resourceClass).
		Do().
		Into(result)
	return
}

// Delete takes name of the resourceClass and deletes it. Returns an error if one occurs.
func (c *resourceClasses) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("resourceclasses").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *resourceClasses) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("resourceclasses").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched resourceClass.
func (c *resourceClasses) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.ResourceClass, err error) {
	result = &v1alpha1.ResourceClass{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("resourceclasses").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
