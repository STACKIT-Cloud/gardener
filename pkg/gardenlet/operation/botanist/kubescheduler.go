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

package botanist

import (
	"github.com/gardener/gardener/imagevector"
	kubescheduler "github.com/gardener/gardener/pkg/component/kubernetes/scheduler"
	imagevectorutils "github.com/gardener/gardener/pkg/utils/imagevector"
)

// DefaultKubeScheduler returns a deployer for the kube-scheduler.
func (b *Botanist) DefaultKubeScheduler() (kubescheduler.Interface, error) {
	image, err := imagevector.ImageVector().FindImage(imagevector.ImageNameKubeScheduler, imagevectorutils.RuntimeVersion(b.SeedVersion()), imagevectorutils.TargetVersion(b.ShootVersion()))
	if err != nil {
		return nil, err
	}

	return kubescheduler.New(
		b.SeedClientSet.Client(),
		b.Shoot.SeedNamespace,
		b.SecretsManager,
		b.Seed.KubernetesVersion,
		b.Shoot.KubernetesVersion,
		image.String(),
		b.Shoot.GetReplicas(1),
		b.Shoot.GetInfo().Spec.Kubernetes.KubeScheduler,
	), nil
}
