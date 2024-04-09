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

package botanist_test

import (
	"context"
	"errors"

	"github.com/go-logr/logr"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"

	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	v1beta1constants "github.com/gardener/gardener/pkg/apis/core/v1beta1/constants"
	v1beta1helper "github.com/gardener/gardener/pkg/apis/core/v1beta1/helper"
	"github.com/gardener/gardener/pkg/client/kubernetes"
	kubernetesfake "github.com/gardener/gardener/pkg/client/kubernetes/fake"
	"github.com/gardener/gardener/pkg/gardenlet/operation"
	. "github.com/gardener/gardener/pkg/gardenlet/operation/botanist"
	shootpkg "github.com/gardener/gardener/pkg/gardenlet/operation/shoot"
)

var _ = Describe("Waiter", func() {
	var (
		botanist *Botanist

		ctx = context.Background()
	)

	BeforeEach(func() {
		shootClient := fakeclient.NewClientBuilder().WithScheme(kubernetes.ShootScheme).Build()
		shootClientSet := kubernetesfake.NewClientSetBuilder().WithClient(shootClient).Build()

		botanist = &Botanist{
			Operation: &operation.Operation{
				Logger:         logr.Discard(),
				ShootClientSet: shootClientSet,
				Shoot:          &shootpkg.Shoot{},
			},
		}
	})

	Describe("#WaitUntilNodesDeleted", func() {
		var (
			node *corev1.Node
		)

		BeforeEach(func() {
			node = &corev1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name: "infinity-node",
				},
			}
		})

		It("should return ok when no node is found", func() {
			Expect(botanist.WaitUntilNodesDeleted(ctx)).To(Succeed())
		})

		It("should return an error when a node is still present", func() {
			ctxCanceled, cancel := context.WithCancel(ctx)
			cancel()

			Expect(botanist.ShootClientSet.Client().Create(ctx, node)).To(Succeed())

			err := botanist.WaitUntilNodesDeleted(ctxCanceled)

			Expect(err).To(MatchError("retry failed with context canceled, last error: not all nodes have been deleted in the shoot cluster"))
		})
	})

	Describe("#WaitUntilNoPodRunning", func() {
		var (
			pod *corev1.Pod
		)

		BeforeEach(func() {
			pod = &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: "infinity-pod",
				},
			}
		})

		It("should return ok when no pod is found", func() {
			Expect(botanist.WaitUntilNoPodRunning(ctx)).To(Succeed())
		})

		It("should return ok when no pod is running", func() {
			ctxCanceled, cancel := context.WithCancel(ctx)
			cancel()

			pod.Status = corev1.PodStatus{Phase: corev1.PodFailed}

			Expect(botanist.ShootClientSet.Client().Create(ctx, pod)).To(Succeed())

			Expect(botanist.WaitUntilNoPodRunning(ctxCanceled)).To(Succeed())
		})

		It("should return an error with when a pod is running in non system namespace", func() {
			ctxCanceled, cancel := context.WithCancel(ctx)
			cancel()

			pod.Status = corev1.PodStatus{Phase: corev1.PodRunning}
			pod.Namespace = "foo"

			Expect(botanist.ShootClientSet.Client().Create(ctx, pod)).To(Succeed())

			err := botanist.WaitUntilNoPodRunning(ctxCanceled)

			var coder v1beta1helper.Coder
			Expect(errors.As(err, &coder)).To(BeTrue())

			Expect(coder.Codes()).To(ContainElement(gardencorev1beta1.ErrorCleanupClusterResources))
		})

		It("should return an error when a pod is running in kube-system namespace", func() {
			ctxCanceled, cancel := context.WithCancel(ctx)
			cancel()

			pod.Status = corev1.PodStatus{Phase: corev1.PodRunning}
			pod.Namespace = metav1.NamespaceSystem

			Expect(botanist.ShootClientSet.Client().Create(ctx, pod)).To(Succeed())

			err := botanist.WaitUntilNoPodRunning(ctxCanceled)

			var coder v1beta1helper.Coder
			Expect(errors.As(err, &coder)).To(BeFalse())

			Expect(err).To(MatchError("retry failed with context canceled, last error: waiting until there are no running Pods in the shoot cluster, there is still at least one running Pod in the shoot cluster: \"kube-system/infinity-pod\""))
		})

		It("should return an error when a pod is running in kubernetes-dashboard namespace", func() {
			ctxCanceled, cancel := context.WithCancel(ctx)
			cancel()

			pod.Status = corev1.PodStatus{Phase: corev1.PodRunning}
			pod.Namespace = v1beta1constants.KubernetesDashboardNamespace

			Expect(botanist.ShootClientSet.Client().Create(ctxCanceled, pod)).To(Succeed())

			err := botanist.WaitUntilNoPodRunning(ctxCanceled)

			var coder v1beta1helper.Coder
			Expect(errors.As(err, &coder)).To(BeFalse())

			Expect(err).To(MatchError("retry failed with context canceled, last error: waiting until there are no running Pods in the shoot cluster, there is still at least one running Pod in the shoot cluster: \"kubernetes-dashboard/infinity-pod\""))
		})
	})

	Describe("#WaitUntilEndpointsDoNotContainPodIPs", func() {
		var (
			shoot    *gardencorev1beta1.Shoot
			endpoint *corev1.Endpoints
		)

		BeforeEach(func() {
			shoot = &gardencorev1beta1.Shoot{
				Spec: gardencorev1beta1.ShootSpec{
					Networking: &gardencorev1beta1.Networking{},
				},
			}
			endpoint = &corev1.Endpoints{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default", Name: "basic-endpoint",
				},
				Subsets: []corev1.EndpointSubset{
					{
						Addresses: []corev1.EndpointAddress{{
							IP:        "10.0.0.1",
							Hostname:  "",
							NodeName:  nil,
							TargetRef: nil,
						}},
						NotReadyAddresses: nil,
						Ports:             nil,
					},
				},
			}
		})

		It("should return an error when shoot's pod network is empty", func() {
			botanist.Shoot.SetInfo(shoot)

			err := botanist.WaitUntilEndpointsDoNotContainPodIPs(ctx)
			Expect(err).To(MatchError("unable to check if there are still Endpoints containing Pod IPs in the shoot cluster. Shoot's Pods network is empty"))
		})

		It("should return an error when shoot's pod network is invalid", func() {
			shoot.Spec.Networking.Pods = ptr.To("abc123")
			botanist.Shoot.SetInfo(shoot)

			err := botanist.WaitUntilEndpointsDoNotContainPodIPs(ctx)

			Expect(err).To(MatchError("unable to check if there are still Endpoints containing Pod IPs in the shoot cluster. Shoots's Pods network could not be parsed: invalid CIDR address: abc123"))
		})

		It("should return an error when an endpoint is still present", func() {
			ctxCanceled, cancel := context.WithCancel(ctx)
			cancel()

			shoot.Spec.Networking.Pods = ptr.To("10.0.0.0/8")
			botanist.Shoot.SetInfo(shoot)

			Expect(botanist.ShootClientSet.Client().Create(ctx, endpoint)).To(Succeed())

			err := botanist.WaitUntilEndpointsDoNotContainPodIPs(ctxCanceled)

			Expect(err).To(MatchError("retry failed with context canceled, last error: waiting until there are no running Pods in the shoot cluster, there is still at least one Endpoint containing pod IPs in the shoot cluster: \"default/basic-endpoint\""))
		})

		It("should succeed when endpoints do not contain pod IPs", func() {
			ctxCanceled, cancel := context.WithCancel(ctx)
			cancel()

			shoot.Spec.Networking.Pods = ptr.To("10.0.0.0/8")
			botanist.Shoot.SetInfo(shoot)

			endpoint.Subsets[0].Addresses[0].IP = "128.0.0.1"
			Expect(botanist.ShootClientSet.Client().Create(ctx, endpoint)).To(Succeed())

			Expect(botanist.WaitUntilEndpointsDoNotContainPodIPs(ctxCanceled)).To(Succeed())
		})

		It("should succeed when no endpoints are present", func() {
			ctxCanceled, cancel := context.WithCancel(ctx)
			cancel()

			shoot.Spec.Networking.Pods = ptr.To("10.0.0.0/8")
			botanist.Shoot.SetInfo(shoot)

			Expect(botanist.WaitUntilEndpointsDoNotContainPodIPs(ctxCanceled)).To(Succeed())
		})
	})
})
