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

	v1alpha1 "github.com/gardener/gardener/pkg/apis/core/v1alpha1"
	scheme "github.com/gardener/gardener/pkg/client/core/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// PlantsGetter has a method to return a PlantInterface.
// A group's client should implement this interface.
type PlantsGetter interface {
	Plants(namespace string) PlantInterface
}

// PlantInterface has methods to work with Plant resources.
type PlantInterface interface {
	Create(*v1alpha1.Plant) (*v1alpha1.Plant, error)
	Update(*v1alpha1.Plant) (*v1alpha1.Plant, error)
	UpdateStatus(*v1alpha1.Plant) (*v1alpha1.Plant, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.Plant, error)
	List(opts v1.ListOptions) (*v1alpha1.PlantList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Plant, err error)
	PlantExpansion
}

// plants implements PlantInterface
type plants struct {
	client rest.Interface
	ns     string
}

// newPlants returns a Plants
func newPlants(c *CoreV1alpha1Client, namespace string) *plants {
	return &plants{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the plant, and returns the corresponding plant object, and an error if there is any.
func (c *plants) Get(name string, options v1.GetOptions) (result *v1alpha1.Plant, err error) {
	result = &v1alpha1.Plant{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("plants").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Plants that match those selectors.
func (c *plants) List(opts v1.ListOptions) (result *v1alpha1.PlantList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.PlantList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("plants").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested plants.
func (c *plants) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("plants").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a plant and creates it.  Returns the server's representation of the plant, and an error, if there is any.
func (c *plants) Create(plant *v1alpha1.Plant) (result *v1alpha1.Plant, err error) {
	result = &v1alpha1.Plant{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("plants").
		Body(plant).
		Do().
		Into(result)
	return
}

// Update takes the representation of a plant and updates it. Returns the server's representation of the plant, and an error, if there is any.
func (c *plants) Update(plant *v1alpha1.Plant) (result *v1alpha1.Plant, err error) {
	result = &v1alpha1.Plant{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("plants").
		Name(plant.Name).
		Body(plant).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *plants) UpdateStatus(plant *v1alpha1.Plant) (result *v1alpha1.Plant, err error) {
	result = &v1alpha1.Plant{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("plants").
		Name(plant.Name).
		SubResource("status").
		Body(plant).
		Do().
		Into(result)
	return
}

// Delete takes name of the plant and deletes it. Returns an error if one occurs.
func (c *plants) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("plants").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *plants) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("plants").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched plant.
func (c *plants) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Plant, err error) {
	result = &v1alpha1.Plant{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("plants").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
