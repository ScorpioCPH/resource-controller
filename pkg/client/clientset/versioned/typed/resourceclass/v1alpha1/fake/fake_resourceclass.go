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

package fake

import (
	v1alpha1 "github.com/caicloud/resource-controller/pkg/apis/resourceclass/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeResourceClasses implements ResourceClassInterface
type FakeResourceClasses struct {
	Fake *FakeResourceclassV1alpha1
	ns   string
}

var resourceclassesResource = schema.GroupVersionResource{Group: "resourceclass.controller.k8s.io", Version: "v1alpha1", Resource: "resourceclasses"}

var resourceclassesKind = schema.GroupVersionKind{Group: "resourceclass.controller.k8s.io", Version: "v1alpha1", Kind: "ResourceClass"}

// Get takes name of the resourceClass, and returns the corresponding resourceClass object, and an error if there is any.
func (c *FakeResourceClasses) Get(name string, options v1.GetOptions) (result *v1alpha1.ResourceClass, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(resourceclassesResource, c.ns, name), &v1alpha1.ResourceClass{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ResourceClass), err
}

// List takes label and field selectors, and returns the list of ResourceClasses that match those selectors.
func (c *FakeResourceClasses) List(opts v1.ListOptions) (result *v1alpha1.ResourceClassList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(resourceclassesResource, resourceclassesKind, c.ns, opts), &v1alpha1.ResourceClassList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.ResourceClassList{}
	for _, item := range obj.(*v1alpha1.ResourceClassList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested resourceClasses.
func (c *FakeResourceClasses) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(resourceclassesResource, c.ns, opts))

}

// Create takes the representation of a resourceClass and creates it.  Returns the server's representation of the resourceClass, and an error, if there is any.
func (c *FakeResourceClasses) Create(resourceClass *v1alpha1.ResourceClass) (result *v1alpha1.ResourceClass, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(resourceclassesResource, c.ns, resourceClass), &v1alpha1.ResourceClass{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ResourceClass), err
}

// Update takes the representation of a resourceClass and updates it. Returns the server's representation of the resourceClass, and an error, if there is any.
func (c *FakeResourceClasses) Update(resourceClass *v1alpha1.ResourceClass) (result *v1alpha1.ResourceClass, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(resourceclassesResource, c.ns, resourceClass), &v1alpha1.ResourceClass{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ResourceClass), err
}

// Delete takes name of the resourceClass and deletes it. Returns an error if one occurs.
func (c *FakeResourceClasses) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(resourceclassesResource, c.ns, name), &v1alpha1.ResourceClass{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeResourceClasses) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(resourceclassesResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.ResourceClassList{})
	return err
}

// Patch applies the patch and returns the patched resourceClass.
func (c *FakeResourceClasses) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.ResourceClass, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(resourceclassesResource, c.ns, name, data, subresources...), &v1alpha1.ResourceClass{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ResourceClass), err
}
