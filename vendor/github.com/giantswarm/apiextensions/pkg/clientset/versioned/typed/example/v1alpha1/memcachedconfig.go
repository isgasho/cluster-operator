/*
Copyright 2019 Giant Swarm GmbH.

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

package v1alpha1

import (
	v1alpha1 "github.com/giantswarm/apiextensions/pkg/apis/example/v1alpha1"
	scheme "github.com/giantswarm/apiextensions/pkg/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// MemcachedConfigsGetter has a method to return a MemcachedConfigInterface.
// A group's client should implement this interface.
type MemcachedConfigsGetter interface {
	MemcachedConfigs(namespace string) MemcachedConfigInterface
}

// MemcachedConfigInterface has methods to work with MemcachedConfig resources.
type MemcachedConfigInterface interface {
	Create(*v1alpha1.MemcachedConfig) (*v1alpha1.MemcachedConfig, error)
	Update(*v1alpha1.MemcachedConfig) (*v1alpha1.MemcachedConfig, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.MemcachedConfig, error)
	List(opts v1.ListOptions) (*v1alpha1.MemcachedConfigList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.MemcachedConfig, err error)
	MemcachedConfigExpansion
}

// memcachedConfigs implements MemcachedConfigInterface
type memcachedConfigs struct {
	client rest.Interface
	ns     string
}

// newMemcachedConfigs returns a MemcachedConfigs
func newMemcachedConfigs(c *ExampleV1alpha1Client, namespace string) *memcachedConfigs {
	return &memcachedConfigs{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the memcachedConfig, and returns the corresponding memcachedConfig object, and an error if there is any.
func (c *memcachedConfigs) Get(name string, options v1.GetOptions) (result *v1alpha1.MemcachedConfig, err error) {
	result = &v1alpha1.MemcachedConfig{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("memcachedconfigs").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of MemcachedConfigs that match those selectors.
func (c *memcachedConfigs) List(opts v1.ListOptions) (result *v1alpha1.MemcachedConfigList, err error) {
	result = &v1alpha1.MemcachedConfigList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("memcachedconfigs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested memcachedConfigs.
func (c *memcachedConfigs) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("memcachedconfigs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a memcachedConfig and creates it.  Returns the server's representation of the memcachedConfig, and an error, if there is any.
func (c *memcachedConfigs) Create(memcachedConfig *v1alpha1.MemcachedConfig) (result *v1alpha1.MemcachedConfig, err error) {
	result = &v1alpha1.MemcachedConfig{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("memcachedconfigs").
		Body(memcachedConfig).
		Do().
		Into(result)
	return
}

// Update takes the representation of a memcachedConfig and updates it. Returns the server's representation of the memcachedConfig, and an error, if there is any.
func (c *memcachedConfigs) Update(memcachedConfig *v1alpha1.MemcachedConfig) (result *v1alpha1.MemcachedConfig, err error) {
	result = &v1alpha1.MemcachedConfig{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("memcachedconfigs").
		Name(memcachedConfig.Name).
		Body(memcachedConfig).
		Do().
		Into(result)
	return
}

// Delete takes name of the memcachedConfig and deletes it. Returns an error if one occurs.
func (c *memcachedConfigs) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("memcachedconfigs").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *memcachedConfigs) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("memcachedconfigs").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched memcachedConfig.
func (c *memcachedConfigs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.MemcachedConfig, err error) {
	result = &v1alpha1.MemcachedConfig{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("memcachedconfigs").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
