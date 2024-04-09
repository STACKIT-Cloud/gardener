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

package validation_test

import (
	"time"

	"github.com/Masterminds/semver/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"

	. "github.com/gardener/gardener/pkg/nodeagent/apis/config"
	. "github.com/gardener/gardener/pkg/nodeagent/apis/config/validation"
)

var _ = Describe("#ValidateNodeAgentConfiguration", func() {
	var config *NodeAgentConfiguration

	BeforeEach(func() {
		config = &NodeAgentConfiguration{
			Controllers: ControllerConfiguration{
				OperatingSystemConfig: OperatingSystemConfigControllerConfig{
					SecretName:        "osc-secret",
					SyncPeriod:        &metav1.Duration{Duration: time.Minute},
					KubernetesVersion: semver.MustParse("v1.27.0"),
				},
				Token: TokenControllerConfig{
					SyncConfigs: []TokenSecretSyncConfig{{
						SecretName: "gardener-node-agent",
						Path:       "/var/lib/gardener-node-agent/credentials/token",
					}},
					SyncPeriod: &metav1.Duration{Duration: time.Hour},
				},
			},
		}
	})

	It("should pass because all necessary fields is specified", func() {
		Expect(ValidateNodeAgentConfiguration(config)).To(BeEmpty())
	})

	Context("Operating System Config Controller", func() {
		It("should fail because kubernetes version is empty", func() {
			config.Controllers.OperatingSystemConfig.KubernetesVersion = nil

			Expect(ValidateNodeAgentConfiguration(config)).To(ConsistOf(
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":  Equal(field.ErrorTypeRequired),
					"Field": Equal("controllers.operatingSystemConfig.kubernetesVersion"),
				})),
			))
		})

		It("should fail because operating system config secret name is not specified", func() {
			config.Controllers.OperatingSystemConfig.SecretName = ""

			Expect(ValidateNodeAgentConfiguration(config)).To(ConsistOf(
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":  Equal(field.ErrorTypeRequired),
					"Field": Equal("controllers.operatingSystemConfig.secretName"),
				})),
			))
		})

		It("should fail because sync period is too small", func() {
			config.Controllers.OperatingSystemConfig.SyncPeriod.Duration = 10 * time.Second

			Expect(ValidateNodeAgentConfiguration(config)).To(ConsistOf(
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":  Equal(field.ErrorTypeInvalid),
					"Field": Equal("controllers.operatingSystemConfig.syncPeriod"),
				})),
			))
		})
	})

	Context("Token Controller", func() {
		It("should fail because access token secret name is not specified", func() {
			config.Controllers.Token.SyncConfigs = append(config.Controllers.Token.SyncConfigs, TokenSecretSyncConfig{
				Path: "/some/path",
			})

			Expect(ValidateNodeAgentConfiguration(config)).To(ConsistOf(
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":  Equal(field.ErrorTypeRequired),
					"Field": Equal("controllers.token.syncConfigs[1].secretName"),
				})),
			))
		})

		It("should fail because path is not specified", func() {
			config.Controllers.Token.SyncConfigs = append(config.Controllers.Token.SyncConfigs, TokenSecretSyncConfig{
				SecretName: "/some/secret",
			})

			Expect(ValidateNodeAgentConfiguration(config)).To(ConsistOf(
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":  Equal(field.ErrorTypeRequired),
					"Field": Equal("controllers.token.syncConfigs[1].path"),
				})),
			))
		})

		It("should fail because path is duplicated", func() {
			config.Controllers.Token.SyncConfigs = append(config.Controllers.Token.SyncConfigs, TokenSecretSyncConfig{
				SecretName: "other-secret",
				Path:       "/var/lib/gardener-node-agent/credentials/token",
			})

			Expect(ValidateNodeAgentConfiguration(config)).To(ConsistOf(
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":  Equal(field.ErrorTypeDuplicate),
					"Field": Equal("controllers.token.syncConfigs[1].path"),
				})),
			))
		})

		It("should fail because gardener-node-agent access token config is missing", func() {
			config.Controllers.Token.SyncConfigs = nil

			Expect(ValidateNodeAgentConfiguration(config)).To(ConsistOf(
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":  Equal(field.ErrorTypeRequired),
					"Field": Equal("controllers.token.syncConfigs"),
				})),
			))
		})

		It("should fail because sync period is too small", func() {
			config.Controllers.Token.SyncPeriod.Duration = 10 * time.Second

			Expect(ValidateNodeAgentConfiguration(config)).To(ConsistOf(
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":  Equal(field.ErrorTypeInvalid),
					"Field": Equal("controllers.token.syncPeriod"),
				})),
			))
		})
	})
})
