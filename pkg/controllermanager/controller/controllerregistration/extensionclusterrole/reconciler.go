// Copyright 2024 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
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

package extensionclusterrole

import (
	"context"
	"fmt"
	"sort"
	"strings"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/gardener/gardener/pkg/controllerutils"
	gardenerutils "github.com/gardener/gardener/pkg/utils/gardener"
)

// Reconciler reconciles ClusterRoles for additional extension permissions and creates ClusterRoleBindings for binding
// extension service accounts to such ClusterRoles.
type Reconciler struct {
	Client client.Client
}

// Reconcile reconciles ClusterRoles. It creates related ClusterRoleBindings or updates their subjects to match
// ServiceAccounts selected by the ClusterRoles via 'authorization.gardener.cloud/extensions-serviceaccount-selector'
// annotation.
func (r *Reconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	log := logf.FromContext(ctx)

	ctx, cancel := controllerutils.GetMainReconciliationContext(ctx, controllerutils.DefaultReconciliationTimeout)
	defer cancel()

	clusterRole := &metav1.PartialObjectMetadata{}
	clusterRole.SetGroupVersionKind(rbacv1.SchemeGroupVersion.WithKind("ClusterRole"))
	if err := r.Client.Get(ctx, request.NamespacedName, clusterRole); err != nil {
		if apierrors.IsNotFound(err) {
			log.V(1).Info("Object is gone, stop reconciling")
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, fmt.Errorf("error retrieving object from store: %w", err)
	}

	if clusterRole.DeletionTimestamp != nil {
		return reconcile.Result{}, nil
	}

	subjects, err := r.computeSubjects(ctx, clusterRole)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed computing subjects for ClusterRoleBinding: %w", err)
	}

	clusterRoleBinding := &rbacv1.ClusterRoleBinding{ObjectMeta: metav1.ObjectMeta{Name: clusterRole.Name}}
	_, err = controllerutils.GetAndCreateOrMergePatch(ctx, r.Client, clusterRoleBinding, func() error {
		ownerReference := metav1.NewControllerRef(clusterRole, rbacv1.SchemeGroupVersion.WithKind("ClusterRole"))
		ownerReference.BlockOwnerDeletion = ptr.To(false)
		clusterRoleBinding.OwnerReferences = []metav1.OwnerReference{*ownerReference}

		clusterRoleBinding.RoleRef = rbacv1.RoleRef{
			APIGroup: rbacv1.SchemeGroupVersion.Group,
			Kind:     "ClusterRole",
			Name:     clusterRole.Name,
		}
		clusterRoleBinding.Subjects = subjects
		return nil
	})
	return reconcile.Result{}, err
}

func (r *Reconciler) computeSubjects(ctx context.Context, clusterRole *metav1.PartialObjectMetadata) ([]rbacv1.Subject, error) {
	labelSelector, err := clusterRoleServiceAccountLabelSelectorToSelector(clusterRole.Annotations)
	if err != nil {
		return nil, fmt.Errorf("failed parsing label selector: %w", err)
	}

	serviceAccountList := &metav1.PartialObjectMetadataList{}
	serviceAccountList.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("ServiceAccountList"))
	if err := r.Client.List(ctx, serviceAccountList, client.MatchingLabelsSelector{Selector: labelSelector}); err != nil {
		return nil, fmt.Errorf("failed listing ServiceAccounts with label selector %s: %w", labelSelector.String(), err)
	}

	var subjects []rbacv1.Subject
	for _, serviceAccount := range serviceAccountList.Items {
		if strings.HasPrefix(serviceAccount.GetNamespace(), gardenerutils.SeedNamespaceNamePrefix) {
			subjects = append(subjects, rbacv1.Subject{
				Kind:      rbacv1.ServiceAccountKind,
				Name:      serviceAccount.Name,
				Namespace: serviceAccount.Namespace,
			})
		}
	}

	// ensure stable order of subjects
	sort.Slice(subjects, func(i, j int) bool {
		if subjects[i].Namespace != subjects[j].Namespace {
			return subjects[i].Namespace < subjects[j].Namespace
		}
		return subjects[i].Name < subjects[j].Name
	})

	return subjects, nil
}
