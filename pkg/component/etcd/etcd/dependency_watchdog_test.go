// Copyright 2020 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
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

package etcd_test

import (
	weederapi "github.com/gardener/dependency-watchdog/api/weeder"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1beta1constants "github.com/gardener/gardener/pkg/apis/core/v1beta1/constants"
	"github.com/gardener/gardener/pkg/component/etcd/etcd"
)

var _ = Describe("DependencyWatchdog", func() {
	Describe("#NewDependencyWatchdogWeederConfiguration", func() {
		It("should compute the correct configuration", func() {
			config, err := etcd.NewDependencyWatchdogWeederConfiguration(testRole)
			Expect(config).To(Equal(map[string]weederapi.DependantSelectors{
				"etcd-" + testRole + "-client": {
					PodSelectors: []*metav1.LabelSelector{
						{
							MatchExpressions: []metav1.LabelSelectorRequirement{
								{
									Key:      v1beta1constants.GardenRole,
									Operator: "In",
									Values:   []string{v1beta1constants.GardenRoleControlPlane},
								},
								{
									Key:      v1beta1constants.LabelRole,
									Operator: "In",
									Values:   []string{v1beta1constants.LabelAPIServer},
								},
							},
						},
					},
				},
			}))
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
