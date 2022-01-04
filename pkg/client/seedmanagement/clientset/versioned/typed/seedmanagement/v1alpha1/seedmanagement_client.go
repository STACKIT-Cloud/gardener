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

package v1alpha1

import (
	v1alpha1 "github.com/gardener/gardener/pkg/apis/seedmanagement/v1alpha1"
	"github.com/gardener/gardener/pkg/client/seedmanagement/clientset/versioned/scheme"
	rest "k8s.io/client-go/rest"
)

type SeedmanagementV1alpha1Interface interface {
	RESTClient() rest.Interface
	ManagedSeedsGetter
	ManagedSeedSetsGetter
}

// SeedmanagementV1alpha1Client is used to interact with features provided by the seedmanagement.gardener.cloud group.
type SeedmanagementV1alpha1Client struct {
	restClient rest.Interface
}

func (c *SeedmanagementV1alpha1Client) ManagedSeeds(namespace string) ManagedSeedInterface {
	return newManagedSeeds(c, namespace)
}

func (c *SeedmanagementV1alpha1Client) ManagedSeedSets(namespace string) ManagedSeedSetInterface {
	return newManagedSeedSets(c, namespace)
}

// NewForConfig creates a new SeedmanagementV1alpha1Client for the given config.
func NewForConfig(c *rest.Config) (*SeedmanagementV1alpha1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &SeedmanagementV1alpha1Client{client}, nil
}

// NewForConfigOrDie creates a new SeedmanagementV1alpha1Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *SeedmanagementV1alpha1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new SeedmanagementV1alpha1Client for the given RESTClient.
func New(c rest.Interface) *SeedmanagementV1alpha1Client {
	return &SeedmanagementV1alpha1Client{c}
}

func setConfigDefaults(config *rest.Config) error {
	gv := v1alpha1.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *SeedmanagementV1alpha1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
