/*
Copyright (c) 2021 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file

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
	"time"

	v1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	scheme "github.com/gardener/gardener/pkg/client/extensions/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// OperatingSystemConfigsGetter has a method to return a OperatingSystemConfigInterface.
// A group's client should implement this interface.
type OperatingSystemConfigsGetter interface {
	OperatingSystemConfigs(namespace string) OperatingSystemConfigInterface
}

// OperatingSystemConfigInterface has methods to work with OperatingSystemConfig resources.
type OperatingSystemConfigInterface interface {
	Create(*v1alpha1.OperatingSystemConfig) (*v1alpha1.OperatingSystemConfig, error)
	Update(*v1alpha1.OperatingSystemConfig) (*v1alpha1.OperatingSystemConfig, error)
	UpdateStatus(*v1alpha1.OperatingSystemConfig) (*v1alpha1.OperatingSystemConfig, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.OperatingSystemConfig, error)
	List(opts v1.ListOptions) (*v1alpha1.OperatingSystemConfigList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.OperatingSystemConfig, err error)
	OperatingSystemConfigExpansion
}

// operatingSystemConfigs implements OperatingSystemConfigInterface
type operatingSystemConfigs struct {
	client rest.Interface
	ns     string
}

// newOperatingSystemConfigs returns a OperatingSystemConfigs
func newOperatingSystemConfigs(c *ExtensionsV1alpha1Client, namespace string) *operatingSystemConfigs {
	return &operatingSystemConfigs{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the operatingSystemConfig, and returns the corresponding operatingSystemConfig object, and an error if there is any.
func (c *operatingSystemConfigs) Get(name string, options v1.GetOptions) (result *v1alpha1.OperatingSystemConfig, err error) {
	result = &v1alpha1.OperatingSystemConfig{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("operatingsystemconfigs").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of OperatingSystemConfigs that match those selectors.
func (c *operatingSystemConfigs) List(opts v1.ListOptions) (result *v1alpha1.OperatingSystemConfigList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.OperatingSystemConfigList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("operatingsystemconfigs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested operatingSystemConfigs.
func (c *operatingSystemConfigs) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("operatingsystemconfigs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a operatingSystemConfig and creates it.  Returns the server's representation of the operatingSystemConfig, and an error, if there is any.
func (c *operatingSystemConfigs) Create(operatingSystemConfig *v1alpha1.OperatingSystemConfig) (result *v1alpha1.OperatingSystemConfig, err error) {
	result = &v1alpha1.OperatingSystemConfig{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("operatingsystemconfigs").
		Body(operatingSystemConfig).
		Do().
		Into(result)
	return
}

// Update takes the representation of a operatingSystemConfig and updates it. Returns the server's representation of the operatingSystemConfig, and an error, if there is any.
func (c *operatingSystemConfigs) Update(operatingSystemConfig *v1alpha1.OperatingSystemConfig) (result *v1alpha1.OperatingSystemConfig, err error) {
	result = &v1alpha1.OperatingSystemConfig{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("operatingsystemconfigs").
		Name(operatingSystemConfig.Name).
		Body(operatingSystemConfig).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *operatingSystemConfigs) UpdateStatus(operatingSystemConfig *v1alpha1.OperatingSystemConfig) (result *v1alpha1.OperatingSystemConfig, err error) {
	result = &v1alpha1.OperatingSystemConfig{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("operatingsystemconfigs").
		Name(operatingSystemConfig.Name).
		SubResource("status").
		Body(operatingSystemConfig).
		Do().
		Into(result)
	return
}

// Delete takes name of the operatingSystemConfig and deletes it. Returns an error if one occurs.
func (c *operatingSystemConfigs) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("operatingsystemconfigs").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *operatingSystemConfigs) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("operatingsystemconfigs").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched operatingSystemConfig.
func (c *operatingSystemConfigs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.OperatingSystemConfig, err error) {
	result = &v1alpha1.OperatingSystemConfig{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("operatingsystemconfigs").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
