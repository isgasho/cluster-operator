/*
Copyright 2020 Giant Swarm GmbH.

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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "github.com/giantswarm/apiextensions/pkg/apis/release/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeReleaseCycles implements ReleaseCycleInterface
type FakeReleaseCycles struct {
	Fake *FakeReleaseV1alpha1
}

var releasecyclesResource = schema.GroupVersionResource{Group: "release.giantswarm.io", Version: "v1alpha1", Resource: "releasecycles"}

var releasecyclesKind = schema.GroupVersionKind{Group: "release.giantswarm.io", Version: "v1alpha1", Kind: "ReleaseCycle"}

// Get takes name of the releaseCycle, and returns the corresponding releaseCycle object, and an error if there is any.
func (c *FakeReleaseCycles) Get(name string, options v1.GetOptions) (result *v1alpha1.ReleaseCycle, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(releasecyclesResource, name), &v1alpha1.ReleaseCycle{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ReleaseCycle), err
}

// List takes label and field selectors, and returns the list of ReleaseCycles that match those selectors.
func (c *FakeReleaseCycles) List(opts v1.ListOptions) (result *v1alpha1.ReleaseCycleList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(releasecyclesResource, releasecyclesKind, opts), &v1alpha1.ReleaseCycleList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.ReleaseCycleList{ListMeta: obj.(*v1alpha1.ReleaseCycleList).ListMeta}
	for _, item := range obj.(*v1alpha1.ReleaseCycleList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested releaseCycles.
func (c *FakeReleaseCycles) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(releasecyclesResource, opts))
}

// Create takes the representation of a releaseCycle and creates it.  Returns the server's representation of the releaseCycle, and an error, if there is any.
func (c *FakeReleaseCycles) Create(releaseCycle *v1alpha1.ReleaseCycle) (result *v1alpha1.ReleaseCycle, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(releasecyclesResource, releaseCycle), &v1alpha1.ReleaseCycle{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ReleaseCycle), err
}

// Update takes the representation of a releaseCycle and updates it. Returns the server's representation of the releaseCycle, and an error, if there is any.
func (c *FakeReleaseCycles) Update(releaseCycle *v1alpha1.ReleaseCycle) (result *v1alpha1.ReleaseCycle, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(releasecyclesResource, releaseCycle), &v1alpha1.ReleaseCycle{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ReleaseCycle), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeReleaseCycles) UpdateStatus(releaseCycle *v1alpha1.ReleaseCycle) (*v1alpha1.ReleaseCycle, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(releasecyclesResource, "status", releaseCycle), &v1alpha1.ReleaseCycle{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ReleaseCycle), err
}

// Delete takes name of the releaseCycle and deletes it. Returns an error if one occurs.
func (c *FakeReleaseCycles) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(releasecyclesResource, name), &v1alpha1.ReleaseCycle{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeReleaseCycles) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(releasecyclesResource, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.ReleaseCycleList{})
	return err
}

// Patch applies the patch and returns the patched releaseCycle.
func (c *FakeReleaseCycles) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.ReleaseCycle, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(releasecyclesResource, name, pt, data, subresources...), &v1alpha1.ReleaseCycle{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ReleaseCycle), err
}