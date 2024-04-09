// Copyright 2023 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package apiserverexposure

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	v1beta1constants "github.com/gardener/gardener/pkg/apis/core/v1beta1/constants"
	"github.com/gardener/gardener/pkg/component"
	"github.com/gardener/gardener/pkg/controllerutils"
)

type internalNameService struct {
	client    client.Client
	namespace string
}

// NewInternalNameService creates a new instance of Deployer for the service kubernetes.default.svc.cluster.local.
func NewInternalNameService(c client.Client, namespace string) component.Deployer {
	return &internalNameService{
		client:    c,
		namespace: namespace,
	}
}

func (in *internalNameService) Deploy(ctx context.Context) error {
	svc := in.emptyKubernetesDefaultService()

	if _, err := controllerutils.GetAndCreateOrMergePatch(ctx, in.client, svc, func() error {
		metav1.SetMetaDataAnnotation(&svc.ObjectMeta, "networking.istio.io/exportTo", "*")
		return nil
	}); err != nil {
		return err
	}
	return client.IgnoreNotFound(in.client.Delete(ctx, in.emptyService()))
}

func (in *internalNameService) Destroy(ctx context.Context) error {
	return client.IgnoreNotFound(in.client.Delete(ctx, in.emptyService()))
}

func (in *internalNameService) emptyService() *corev1.Service {
	return &corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: v1beta1constants.DeploymentNameKubeAPIServer, Namespace: in.namespace}}
}

func (in *internalNameService) emptyKubernetesDefaultService() *corev1.Service {
	return &corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: "kubernetes", Namespace: metav1.NamespaceDefault}}
}
