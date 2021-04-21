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

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	time "time"

	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	versioned "github.com/gardener/gardener/pkg/client/extensions/clientset/versioned"
	internalinterfaces "github.com/gardener/gardener/pkg/client/extensions/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/gardener/gardener/pkg/client/extensions/listers/extensions/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// ContainerRuntimeInformer provides access to a shared informer and lister for
// ContainerRuntimes.
type ContainerRuntimeInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.ContainerRuntimeLister
}

type containerRuntimeInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewContainerRuntimeInformer constructs a new informer for ContainerRuntime type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewContainerRuntimeInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredContainerRuntimeInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredContainerRuntimeInformer constructs a new informer for ContainerRuntime type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredContainerRuntimeInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ExtensionsV1alpha1().ContainerRuntimes(namespace).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ExtensionsV1alpha1().ContainerRuntimes(namespace).Watch(options)
			},
		},
		&extensionsv1alpha1.ContainerRuntime{},
		resyncPeriod,
		indexers,
	)
}

func (f *containerRuntimeInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredContainerRuntimeInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *containerRuntimeInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&extensionsv1alpha1.ContainerRuntime{}, f.defaultInformer)
}

func (f *containerRuntimeInformer) Lister() v1alpha1.ContainerRuntimeLister {
	return v1alpha1.NewContainerRuntimeLister(f.Informer().GetIndexer())
}
