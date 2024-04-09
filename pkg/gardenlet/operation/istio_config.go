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

package operation

import (
	"strings"

	"k8s.io/apimachinery/pkg/util/sets"

	resourcesv1alpha1 "github.com/gardener/gardener/pkg/apis/resources/v1alpha1"
	sharedcomponent "github.com/gardener/gardener/pkg/component/shared"
	gardenletconfig "github.com/gardener/gardener/pkg/gardenlet/apis/config"
	"github.com/gardener/gardener/pkg/utils"
	gardenerutils "github.com/gardener/gardener/pkg/utils/gardener"
)

// IstioServiceName is the currently used name of the istio ingress service, which is responsible for the shoot cluster.
func (o *Operation) IstioServiceName() string {
	return *o.sniConfig().Ingress.ServiceName
}

// IstioNamespace is the currently used namespace of the istio ingress gateway, which is responsible for the shoot cluster.
func (o *Operation) IstioNamespace() string {
	return o.addZonePinningIfRequired(o.DefaultIstioNamespace())
}

// DefaultIstioNamespace is the default namespace of the istio ingress gateway disregarding zonal affinities of the shoot cluster.
func (o *Operation) DefaultIstioNamespace() string {
	return *o.sniConfig().Ingress.Namespace
}

// IstioLoadBalancerAnnotations contain the annotation to be used for the istio ingress service load balancer.
func (o *Operation) IstioLoadBalancerAnnotations() map[string]string {
	zone := o.singleZoneIfPinned()
	if exposureClassHandler := o.exposureClassHandler(); exposureClassHandler != nil {
		if zone != nil {
			return utils.MergeStringMaps(exposureClassHandler.LoadBalancerService.Annotations, o.Seed.GetZonalLoadBalancerServiceAnnotations(*zone))
		}
		return utils.MergeStringMaps(o.Seed.GetLoadBalancerServiceAnnotations(), exposureClassHandler.LoadBalancerService.Annotations)
	}
	if zone != nil {
		return o.Seed.GetZonalLoadBalancerServiceAnnotations(*zone)
	}
	return o.Seed.GetLoadBalancerServiceAnnotations()
}

// IstioLabels contain the labels to be used for the istio ingress gateway entities.
func (o *Operation) IstioLabels() map[string]string {
	return o.istioLabels(o.singleZoneIfPinned())
}

// DefaultIstioLabels contain the labels to be used for the default istio ingress gateway entities disregarding zonal affinities.
func (o *Operation) DefaultIstioLabels() map[string]string {
	return o.istioLabels(nil)
}

func (o *Operation) istioLabels(zone *string) map[string]string {
	if exposureClassHandler := o.exposureClassHandler(); exposureClassHandler != nil {
		return sharedcomponent.GetIstioZoneLabels(gardenerutils.GetMandatoryExposureClassHandlerSNILabels(exposureClassHandler.SNI.Ingress.Labels, exposureClassHandler.Name), zone)
	}
	return sharedcomponent.GetIstioZoneLabels(o.sniConfig().Ingress.Labels, zone)
}

func (o *Operation) exposureClassHandler() *gardenletconfig.ExposureClassHandler {
	if exposureClass := o.Shoot.ExposureClass; exposureClass != nil {
		for _, handler := range o.Config.ExposureClassHandlers {
			if exposureClass.Handler == handler.Name {
				return &handler
			}
		}
	}
	return nil
}

func (o *Operation) sniConfig() *gardenletconfig.SNI {
	if exposureClassHandler := o.exposureClassHandler(); exposureClassHandler != nil {
		return exposureClassHandler.SNI
	}
	return o.Config.SNI
}

func (o *Operation) addZonePinningIfRequired(namespace string) string {
	// Only clusters pinned to exactly one zone are exposed via zonal istio ingress gateway.
	// All other clusters, i.e. true HA and legacy/accidental multi-zonal clusters, are exposed
	// via the default multi-zonal istio ingress gateway.
	zone := o.singleZoneIfPinned()
	if zone != nil {
		return sharedcomponent.GetIstioNamespaceForZone(namespace, *zone)
	}
	return namespace
}

func (o *Operation) singleZoneIfPinned() *string {
	// Zone-specific istio ingress gateways are only deployed with more than one zone
	if len(o.Seed.GetInfo().Spec.Provider.Zones) <= 1 {
		return nil
	}
	if v, ok := o.SeedNamespaceObject.Annotations[resourcesv1alpha1.HighAvailabilityConfigZones]; ok {
		zones := sets.List(sets.New(strings.Split(v, ",")...).Delete(""))
		if len(zones) == 1 {
			return &zones[0]
		}
	}
	return nil
}
