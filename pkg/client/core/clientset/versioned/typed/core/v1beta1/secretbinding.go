// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

// Code generated by client-gen. DO NOT EDIT.

package v1beta1

import (
	"context"
	"time"

	v1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	scheme "github.com/gardener/gardener/pkg/client/core/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// SecretBindingsGetter has a method to return a SecretBindingInterface.
// A group's client should implement this interface.
type SecretBindingsGetter interface {
	SecretBindings(namespace string) SecretBindingInterface
}

// SecretBindingInterface has methods to work with SecretBinding resources.
type SecretBindingInterface interface {
	Create(ctx context.Context, secretBinding *v1beta1.SecretBinding, opts v1.CreateOptions) (*v1beta1.SecretBinding, error)
	Update(ctx context.Context, secretBinding *v1beta1.SecretBinding, opts v1.UpdateOptions) (*v1beta1.SecretBinding, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1beta1.SecretBinding, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1beta1.SecretBindingList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.SecretBinding, err error)
	SecretBindingExpansion
}

// secretBindings implements SecretBindingInterface
type secretBindings struct {
	client rest.Interface
	ns     string
}

// newSecretBindings returns a SecretBindings
func newSecretBindings(c *CoreV1beta1Client, namespace string) *secretBindings {
	return &secretBindings{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the secretBinding, and returns the corresponding secretBinding object, and an error if there is any.
func (c *secretBindings) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1beta1.SecretBinding, err error) {
	result = &v1beta1.SecretBinding{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("secretbindings").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of SecretBindings that match those selectors.
func (c *secretBindings) List(ctx context.Context, opts v1.ListOptions) (result *v1beta1.SecretBindingList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1beta1.SecretBindingList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("secretbindings").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested secretBindings.
func (c *secretBindings) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("secretbindings").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a secretBinding and creates it.  Returns the server's representation of the secretBinding, and an error, if there is any.
func (c *secretBindings) Create(ctx context.Context, secretBinding *v1beta1.SecretBinding, opts v1.CreateOptions) (result *v1beta1.SecretBinding, err error) {
	result = &v1beta1.SecretBinding{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("secretbindings").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(secretBinding).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a secretBinding and updates it. Returns the server's representation of the secretBinding, and an error, if there is any.
func (c *secretBindings) Update(ctx context.Context, secretBinding *v1beta1.SecretBinding, opts v1.UpdateOptions) (result *v1beta1.SecretBinding, err error) {
	result = &v1beta1.SecretBinding{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("secretbindings").
		Name(secretBinding.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(secretBinding).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the secretBinding and deletes it. Returns an error if one occurs.
func (c *secretBindings) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("secretbindings").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *secretBindings) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("secretbindings").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched secretBinding.
func (c *secretBindings) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.SecretBinding, err error) {
	result = &v1beta1.SecretBinding{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("secretbindings").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
