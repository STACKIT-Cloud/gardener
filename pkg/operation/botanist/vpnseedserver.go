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

	"strings"

	v1beta1constants "github.com/gardener/gardener/pkg/apis/core/v1beta1/constants"
	"github.com/gardener/gardener/pkg/operation/botanist/component"
	"github.com/gardener/gardener/pkg/operation/botanist/component/vpnseedserver"
	"github.com/gardener/gardener/pkg/utils"
	"github.com/gardener/gardener/pkg/utils/images"
	"github.com/gardener/gardener/pkg/utils/imagevector"

	"k8s.io/utils/pointer"
)

var (
	diffieHellmanKeyData = map[string][]byte{
		"dh2048.pem": []byte(`-----BEGIN DH PARAMETERS-----
MIIBCAKCAQEA7cBXxG9an6KRz/sB5uiSOTf7Eg+uWVkhXO4peKDTARzMYa8b7WR8
B/Aw+AyUXtB3tXtrzeC5M3IHnuhFwMo3K4oSOkFJxatLlYKeY15r+Kt5vnOOT3BW
eN5OnWlR5Wi7GZBWbaQgXVR79N4yst43sVhJus6By0lN6Olc9xD/ys9GH/ykJVIh
Z/NLrxAC5lxjwCqJMd8hrryChuDlz597vg6gYFuRV60U/YU4DK71F4H7mI07aGJ9
l+SK8TbkKWF5ITI7kYWbc4zmtfXSXaGjMhM9omQUaTH9csB96hzFJdeZ4XjxybRf
Vc3t7XP5q7afeaKmM3FhSXdeHKCTqQzQuwIBAg==
-----END DH PARAMETERS-----
`,
		)}
	diffieHellmanKeyChecksum string
)

// init calculates the checksum of the default diffie hellman key
func init() {
	diffieHellmanKeyChecksum = utils.ComputeChecksum(diffieHellmanKeyData)
}

func (b *Botanist) getDiffieHellmanSecret() component.Secret {
	data, checksum := diffieHellmanKeyData, diffieHellmanKeyChecksum
	if secret := b.LoadSecret(v1beta1constants.GardenRoleOpenVPNDiffieHellman); secret != nil {
		data, checksum = secret.Data, utils.ComputeSecretChecksum(secret.Data)
	}

	return component.Secret{
		Name:     v1beta1constants.GardenRoleOpenVPNDiffieHellman,
		Data:     data,
		Checksum: checksum,
	}
}

// DefaultVPNSeedServer returns a deployer for the vpn-seed-server.
func (b *Botanist) DefaultVPNSeedServer() (vpnseedserver.Interface, error) {
	imageAPIServerProxy, err := b.ImageVector.FindImage(images.ImageNameApiserverProxy, imagevector.RuntimeVersion(b.SeedVersion()), imagevector.TargetVersion(b.ShootVersion()))
	if err != nil {
		return nil, err
	}

	imageVPNSeedServer, err := b.ImageVector.FindImage(images.ImageNameVpnSeedServer, imagevector.RuntimeVersion(b.SeedVersion()), imagevector.TargetVersion(b.ShootVersion()))
	if err != nil {
		return nil, err
	}

	var kubeAPIServerHost *string
	if b.APIServerSNIEnabled() {
		kubeAPIServerHost = pointer.String(b.outOfClusterAPIServerFQDN())
	}

	var netPods []string
	for _, IPNetPod := range b.Shoot.Networks.Pods {
		netPods = append(netPods, IPNetPod.String())
	}

	var netServices []string
	for _, IPNetService := range b.Shoot.Networks.Services {
		netServices = append(netServices, IPNetService.String())
	}

	var netNodes []string
	var nodeNetworks = b.Shoot.GetInfo().Spec.Networking.Nodes
	for _, IPNetNode := range strings.Split(*nodeNetworks, ",") {
		netNodes = append(netNodes, IPNetNode)
	}

	return vpnseedserver.New(
		b.K8sSeedClient.Client(),
		b.Shoot.SeedNamespace,
		b.SecretsManager,
		imageAPIServerProxy.String(),
		imageVPNSeedServer.String(),
		kubeAPIServerHost,
		netServices,
		netPods,
		netNodes,
		b.Shoot.GetReplicas(1),
		vpnseedserver.IstioIngressGateway{
			Namespace: *b.Config.SNI.Ingress.Namespace,
			Labels:    b.Config.SNI.Ingress.Labels,
		},
	), nil
}

// DeployVPNServer deploys the vpn-seed-server.
func (b *Botanist) DeployVPNServer(ctx context.Context) error {
	if !b.Shoot.ReversedVPNEnabled {
		return b.Shoot.Components.ControlPlane.VPNSeedServer.Destroy(ctx)
	}

	b.Shoot.Components.ControlPlane.VPNSeedServer.SetSecrets(vpnseedserver.Secrets{DiffieHellmanKey: b.getDiffieHellmanSecret()})
	b.Shoot.Components.ControlPlane.VPNSeedServer.SetSeedNamespaceObjectUID(b.SeedNamespaceObject.UID)
	b.Shoot.Components.ControlPlane.VPNSeedServer.SetSNIConfig(b.Config.SNI)

	if b.ExposureClassHandler != nil {
		b.Shoot.Components.ControlPlane.VPNSeedServer.SetExposureClassHandlerName(b.ExposureClassHandler.Name)
		b.Shoot.Components.ControlPlane.VPNSeedServer.SetSNIConfig(b.ExposureClassHandler.SNI)
	}

	return b.Shoot.Components.ControlPlane.VPNSeedServer.Deploy(ctx)
}
