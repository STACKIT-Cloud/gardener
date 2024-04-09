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

package istio

import (
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	istioapinetworkingv1beta1 "istio.io/api/networking/v1beta1"
	istionetworkingv1beta1 "istio.io/client-go/pkg/apis/networking/v1beta1"
	corev1 "k8s.io/api/core/v1"
)

// DestinationRuleWithLocalityPreference returns a function setting the given attributes to a destination rule object.
func DestinationRuleWithLocalityPreference(destinationRule *istionetworkingv1beta1.DestinationRule, labels map[string]string, destinationHost string) func() error {
	return DestinationRuleWithLocalityPreferenceAndTLS(destinationRule, labels, destinationHost, istioapinetworkingv1beta1.ClientTLSSettings_DISABLE)
}

// DestinationRuleWithLocalityPreferenceAndTLS returns a function setting the given attributes to a destination rule object.
func DestinationRuleWithLocalityPreferenceAndTLS(destinationRule *istionetworkingv1beta1.DestinationRule, labels map[string]string, destinationHost string, tlsMode istioapinetworkingv1beta1.ClientTLSSettings_TLSmode) func() error {
	return func() error {
		destinationRule.Labels = labels
		destinationRule.Spec = istioapinetworkingv1beta1.DestinationRule{
			ExportTo: []string{"*"},
			Host:     destinationHost,
			TrafficPolicy: &istioapinetworkingv1beta1.TrafficPolicy{
				ConnectionPool: &istioapinetworkingv1beta1.ConnectionPoolSettings{
					Tcp: &istioapinetworkingv1beta1.ConnectionPoolSettings_TCPSettings{
						MaxConnections: 5000,
						TcpKeepalive: &istioapinetworkingv1beta1.ConnectionPoolSettings_TCPSettings_TcpKeepalive{
							Time:     &durationpb.Duration{Seconds: 7200},
							Interval: &durationpb.Duration{Seconds: 75},
						},
					},
				},
				LoadBalancer: &istioapinetworkingv1beta1.LoadBalancerSettings{
					LocalityLbSetting: &istioapinetworkingv1beta1.LocalityLoadBalancerSetting{
						Enabled:          &wrapperspb.BoolValue{Value: true},
						FailoverPriority: []string{corev1.LabelTopologyZone},
					},
				},
				// OutlierDetection is required for locality settings to take effect
				OutlierDetection: &istioapinetworkingv1beta1.OutlierDetection{
					MinHealthPercent: 0,
				},
				Tls: &istioapinetworkingv1beta1.ClientTLSSettings{
					Mode: tlsMode,
				},
			},
		}
		return nil
	}
}
