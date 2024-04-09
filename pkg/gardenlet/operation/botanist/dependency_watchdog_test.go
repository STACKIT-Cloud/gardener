// Copyright 2022 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	mockcomponent "github.com/gardener/gardener/pkg/component/mock"
	"github.com/gardener/gardener/pkg/gardenlet/operation"
	. "github.com/gardener/gardener/pkg/gardenlet/operation/botanist"
	"github.com/gardener/gardener/pkg/gardenlet/operation/seed"
	shootpkg "github.com/gardener/gardener/pkg/gardenlet/operation/shoot"
)

var _ = Describe("DependencyWatchdogAccess", func() {
	var (
		ctrl     *gomock.Controller
		botanist *Botanist
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		botanist = &Botanist{
			Operation: &operation.Operation{
				Seed: &seed.Seed{},
			},
		}
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("#DeployDependencyWatchdogAccess", func() {
		var (
			dwdAccess *mockcomponent.MockDeployer

			ctx     = context.TODO()
			fakeErr = errors.New("fake err")
		)

		BeforeEach(func() {
			dwdAccess = mockcomponent.NewMockDeployer(ctrl)

			botanist.Shoot = &shootpkg.Shoot{
				Components: &shootpkg.Components{
					DependencyWatchdogAccess: dwdAccess,
				},
			}
		})

		Context("DWD prober is enabled", func() {
			BeforeEach(func() {
				botanist.Seed.SetInfo(&gardencorev1beta1.Seed{
					Spec: gardencorev1beta1.SeedSpec{
						Settings: &gardencorev1beta1.SeedSettings{
							DependencyWatchdog: &gardencorev1beta1.SeedSettingDependencyWatchdog{
								Prober: &gardencorev1beta1.SeedSettingDependencyWatchdogProber{
									Enabled: true,
								},
							},
						},
					},
				})
			})

			It("should deploy", func() {
				dwdAccess.EXPECT().Deploy(ctx)
				Expect(botanist.DeployDependencyWatchdogAccess(ctx)).To(Succeed())
			})

			It("should fail when the deploy function fails", func() {
				dwdAccess.EXPECT().Deploy(ctx).Return(fakeErr)
				Expect(botanist.DeployDependencyWatchdogAccess(ctx)).To(MatchError(fakeErr))
			})
		})

		Context("DWD prober is disabled", func() {
			BeforeEach(func() {
				botanist.Seed.SetInfo(&gardencorev1beta1.Seed{
					Spec: gardencorev1beta1.SeedSpec{
						Settings: &gardencorev1beta1.SeedSettings{
							DependencyWatchdog: &gardencorev1beta1.SeedSettingDependencyWatchdog{
								Prober: &gardencorev1beta1.SeedSettingDependencyWatchdogProber{
									Enabled: false,
								},
							},
						},
					},
				})
			})

			It("should destroy", func() {
				dwdAccess.EXPECT().Destroy(ctx)
				Expect(botanist.DeployDependencyWatchdogAccess(ctx)).To(Succeed())
			})

			It("should fail when the destroy function fails", func() {
				dwdAccess.EXPECT().Destroy(ctx).Return(fakeErr)
				Expect(botanist.DeployDependencyWatchdogAccess(ctx)).To(MatchError(fakeErr))
			})
		})
	})
})
