/*
Copyright (c) 2022 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file

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

package internalversion

import (
	"context"
	"time"

	core "github.com/gardener/gardener/pkg/apis/core"
	scheme "github.com/gardener/gardener/pkg/client/core/clientset/internalversion/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// ControllerInstallationsGetter has a method to return a ControllerInstallationInterface.
// A group's client should implement this interface.
type ControllerInstallationsGetter interface {
	ControllerInstallations() ControllerInstallationInterface
}

// ControllerInstallationInterface has methods to work with ControllerInstallation resources.
type ControllerInstallationInterface interface {
	Create(ctx context.Context, controllerInstallation *core.ControllerInstallation, opts v1.CreateOptions) (*core.ControllerInstallation, error)
	Update(ctx context.Context, controllerInstallation *core.ControllerInstallation, opts v1.UpdateOptions) (*core.ControllerInstallation, error)
	UpdateStatus(ctx context.Context, controllerInstallation *core.ControllerInstallation, opts v1.UpdateOptions) (*core.ControllerInstallation, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*core.ControllerInstallation, error)
	List(ctx context.Context, opts v1.ListOptions) (*core.ControllerInstallationList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *core.ControllerInstallation, err error)
	ControllerInstallationExpansion
}

// controllerInstallations implements ControllerInstallationInterface
type controllerInstallations struct {
	client rest.Interface
}

// newControllerInstallations returns a ControllerInstallations
func newControllerInstallations(c *CoreClient) *controllerInstallations {
	return &controllerInstallations{
		client: c.RESTClient(),
	}
}

// Get takes name of the controllerInstallation, and returns the corresponding controllerInstallation object, and an error if there is any.
func (c *controllerInstallations) Get(ctx context.Context, name string, options v1.GetOptions) (result *core.ControllerInstallation, err error) {
	result = &core.ControllerInstallation{}
	err = c.client.Get().
		Resource("controllerinstallations").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ControllerInstallations that match those selectors.
func (c *controllerInstallations) List(ctx context.Context, opts v1.ListOptions) (result *core.ControllerInstallationList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &core.ControllerInstallationList{}
	err = c.client.Get().
		Resource("controllerinstallations").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested controllerInstallations.
func (c *controllerInstallations) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("controllerinstallations").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a controllerInstallation and creates it.  Returns the server's representation of the controllerInstallation, and an error, if there is any.
func (c *controllerInstallations) Create(ctx context.Context, controllerInstallation *core.ControllerInstallation, opts v1.CreateOptions) (result *core.ControllerInstallation, err error) {
	result = &core.ControllerInstallation{}
	err = c.client.Post().
		Resource("controllerinstallations").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(controllerInstallation).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a controllerInstallation and updates it. Returns the server's representation of the controllerInstallation, and an error, if there is any.
func (c *controllerInstallations) Update(ctx context.Context, controllerInstallation *core.ControllerInstallation, opts v1.UpdateOptions) (result *core.ControllerInstallation, err error) {
	result = &core.ControllerInstallation{}
	err = c.client.Put().
		Resource("controllerinstallations").
		Name(controllerInstallation.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(controllerInstallation).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *controllerInstallations) UpdateStatus(ctx context.Context, controllerInstallation *core.ControllerInstallation, opts v1.UpdateOptions) (result *core.ControllerInstallation, err error) {
	result = &core.ControllerInstallation{}
	err = c.client.Put().
		Resource("controllerinstallations").
		Name(controllerInstallation.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(controllerInstallation).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the controllerInstallation and deletes it. Returns an error if one occurs.
func (c *controllerInstallations) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("controllerinstallations").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *controllerInstallations) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("controllerinstallations").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched controllerInstallation.
func (c *controllerInstallations) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *core.ControllerInstallation, err error) {
	result = &core.ControllerInstallation{}
	err = c.client.Patch(pt).
		Resource("controllerinstallations").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
