// Copyright (c) 2021 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
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
	"context"

	"github.com/gardener/gardener/pkg/operation/botanist/component/vpnseedserver"
	"github.com/gardener/gardener/pkg/operation/botanist/component/vpnshoot"
	"github.com/gardener/gardener/pkg/utils/images"
	"github.com/gardener/gardener/pkg/utils/imagevector"
	"k8s.io/utils/pointer"
)

// DefaultVPNShoot returns a deployer for the VPNShoot
func (b *Botanist) DefaultVPNShoot() (vpnshoot.Interface, error) {
	var (
		imageName         = images.ImageNameVpnShoot
		reversedVPNValues = vpnshoot.ReversedVPNValues{
			Enabled: false,
		}
	)

	if b.Shoot.ReversedVPNEnabled {
		imageName = images.ImageNameVpnShootClient

		reversedVPNValues = vpnshoot.ReversedVPNValues{
			Enabled:     true,
			Header:      "outbound|1194||" + vpnseedserver.ServiceName + "." + b.Shoot.SeedNamespace + ".svc.cluster.local",
			Endpoint:    b.outOfClusterAPIServerFQDN(),
			OpenVPNPort: 8132,
		}
	}

	image, err := b.ImageVector.FindImage(imageName, imagevector.RuntimeVersion(b.ShootVersion()), imagevector.TargetVersion(b.ShootVersion()))
	if err != nil {
		return nil, err
	}

	values := vpnshoot.Values{
		Image:       image.String(),
		VPAEnabled:  b.Shoot.WantsVerticalPodAutoscaler,
		ReversedVPN: reversedVPNValues,
		Network: vpnshoot.NetworkValues{
			PodCIDR:     b.Shoot.Networks.Pods[0].String(),
			ServiceCIDR: b.Shoot.Networks.Services[0].String(),
			NodeCIDR:    pointer.StringDeref(b.Shoot.GetInfo().Spec.Networking.Nodes, ""),
		},
		HighAvailabilityEnabled:              b.Shoot.VPNHighAvailabilityEnabled,
		HighAvailabilityNumberOfSeedServers:  b.Shoot.VPNHighAvailabilityNumberOfSeedServers,
		HighAvailabilityNumberOfShootClients: b.Shoot.VPNHighAvailabilityNumberOfShootClients,
		PSPDisabled:                          b.Shoot.PSPDisabled,
		KubernetesVersion:                    b.Shoot.KubernetesVersion,
	}

	return vpnshoot.New(
		b.SeedClientSet.Client(),
		b.Shoot.SeedNamespace,
		b.SecretsManager,
		values,
	), nil
}

// DeployVPNShoot deploys the VPNShoot system component.
func (b *Botanist) DeployVPNShoot(ctx context.Context) error {
	secrets := vpnshoot.Secrets{}

	if !b.Shoot.ReversedVPNEnabled {
		dhSecret := b.getDiffieHellmanSecret()
		secrets.DH = &dhSecret
	}

	b.Shoot.Components.SystemComponents.VPNShoot.SetSecrets(secrets)

	return b.Shoot.Components.SystemComponents.VPNShoot.Deploy(ctx)
}
