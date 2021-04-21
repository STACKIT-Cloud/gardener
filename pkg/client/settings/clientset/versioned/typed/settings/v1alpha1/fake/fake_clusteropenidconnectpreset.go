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

package fake

import (
	v1alpha1 "github.com/gardener/gardener/pkg/apis/settings/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeClusterOpenIDConnectPresets implements ClusterOpenIDConnectPresetInterface
type FakeClusterOpenIDConnectPresets struct {
	Fake *FakeSettingsV1alpha1
}

var clusteropenidconnectpresetsResource = schema.GroupVersionResource{Group: "settings.gardener.cloud", Version: "v1alpha1", Resource: "clusteropenidconnectpresets"}

var clusteropenidconnectpresetsKind = schema.GroupVersionKind{Group: "settings.gardener.cloud", Version: "v1alpha1", Kind: "ClusterOpenIDConnectPreset"}

// Get takes name of the clusterOpenIDConnectPreset, and returns the corresponding clusterOpenIDConnectPreset object, and an error if there is any.
func (c *FakeClusterOpenIDConnectPresets) Get(name string, options v1.GetOptions) (result *v1alpha1.ClusterOpenIDConnectPreset, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(clusteropenidconnectpresetsResource, name), &v1alpha1.ClusterOpenIDConnectPreset{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ClusterOpenIDConnectPreset), err
}

// List takes label and field selectors, and returns the list of ClusterOpenIDConnectPresets that match those selectors.
func (c *FakeClusterOpenIDConnectPresets) List(opts v1.ListOptions) (result *v1alpha1.ClusterOpenIDConnectPresetList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(clusteropenidconnectpresetsResource, clusteropenidconnectpresetsKind, opts), &v1alpha1.ClusterOpenIDConnectPresetList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.ClusterOpenIDConnectPresetList{ListMeta: obj.(*v1alpha1.ClusterOpenIDConnectPresetList).ListMeta}
	for _, item := range obj.(*v1alpha1.ClusterOpenIDConnectPresetList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested clusterOpenIDConnectPresets.
func (c *FakeClusterOpenIDConnectPresets) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(clusteropenidconnectpresetsResource, opts))
}

// Create takes the representation of a clusterOpenIDConnectPreset and creates it.  Returns the server's representation of the clusterOpenIDConnectPreset, and an error, if there is any.
func (c *FakeClusterOpenIDConnectPresets) Create(clusterOpenIDConnectPreset *v1alpha1.ClusterOpenIDConnectPreset) (result *v1alpha1.ClusterOpenIDConnectPreset, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(clusteropenidconnectpresetsResource, clusterOpenIDConnectPreset), &v1alpha1.ClusterOpenIDConnectPreset{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ClusterOpenIDConnectPreset), err
}

// Update takes the representation of a clusterOpenIDConnectPreset and updates it. Returns the server's representation of the clusterOpenIDConnectPreset, and an error, if there is any.
func (c *FakeClusterOpenIDConnectPresets) Update(clusterOpenIDConnectPreset *v1alpha1.ClusterOpenIDConnectPreset) (result *v1alpha1.ClusterOpenIDConnectPreset, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(clusteropenidconnectpresetsResource, clusterOpenIDConnectPreset), &v1alpha1.ClusterOpenIDConnectPreset{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ClusterOpenIDConnectPreset), err
}

// Delete takes name of the clusterOpenIDConnectPreset and deletes it. Returns an error if one occurs.
func (c *FakeClusterOpenIDConnectPresets) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(clusteropenidconnectpresetsResource, name), &v1alpha1.ClusterOpenIDConnectPreset{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeClusterOpenIDConnectPresets) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(clusteropenidconnectpresetsResource, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.ClusterOpenIDConnectPresetList{})
	return err
}

// Patch applies the patch and returns the patched clusterOpenIDConnectPreset.
func (c *FakeClusterOpenIDConnectPresets) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.ClusterOpenIDConnectPreset, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(clusteropenidconnectpresetsResource, name, pt, data, subresources...), &v1alpha1.ClusterOpenIDConnectPreset{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ClusterOpenIDConnectPreset), err
}
