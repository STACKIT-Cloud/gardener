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

package service

import (
	"context"
	"fmt"
	"math/rand"
	"strings"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var (
	keyIstioIngressGateway              = client.ObjectKey{Namespace: "istio-ingress", Name: "istio-ingressgateway"}
	keyIstioIngressGatewayZone0         = client.ObjectKey{Namespace: "istio-ingress--0", Name: "istio-ingressgateway"}
	keyIstioIngressGatewayZone1         = client.ObjectKey{Namespace: "istio-ingress--1", Name: "istio-ingressgateway"}
	keyIstioIngressGatewayZone2         = client.ObjectKey{Namespace: "istio-ingress--2", Name: "istio-ingressgateway"}
	keyVirtualGardenIstioIngressGateway = client.ObjectKey{Namespace: "virtual-garden-istio-ingress", Name: "istio-ingressgateway"}
)

const (
	nodePortIstioIngressGateway      int32 = 30443
	nodePortIstioIngressGatewayZone0 int32 = 30444
	nodePortIstioIngressGatewayZone1 int32 = 30445
	nodePortIstioIngressGatewayZone2 int32 = 30446
)

// Reconciler is a reconciler for Service resources.
type Reconciler struct {
	Client      client.Client
	HostIP      string
	Zone0IP     string
	Zone1IP     string
	Zone2IP     string
	IsMultiZone bool
}

// Reconcile reconciles Service resources.
func (r *Reconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	var (
		log     = logf.FromContext(ctx)
		key     = req.NamespacedName
		service = &corev1.Service{}
	)

	if err := r.Client.Get(ctx, key, service); err != nil {
		if apierrors.IsNotFound(err) {
			log.V(1).Info("Object is gone, stop reconciling")
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, fmt.Errorf("error retrieving object from store: %w", err)
	}

	if service.Spec.Type != corev1.ServiceTypeLoadBalancer {
		return reconcile.Result{}, nil
	}

	log.Info("Reconciling service")

	var (
		ip    string
		patch = client.MergeFrom(service.DeepCopy())
	)

	for i, servicePort := range service.Spec.Ports {
		switch {
		case (key == keyIstioIngressGateway || key == keyVirtualGardenIstioIngressGateway) && servicePort.Name == "tcp":
			service.Spec.Ports[i].NodePort = nodePortIstioIngressGateway
			if r.IsMultiZone {
				// Docker desktop for mac v4.23 breaks traffic going through a port mapping to a different docker container.
				// Setting external traffic policy to local mitigates the issue for multi-node setups, e.g. for gardener-operator.
				service.Spec.ExternalTrafficPolicy = corev1.ServiceExternalTrafficPolicyLocal
			}
			ip = r.HostIP
		case key == keyIstioIngressGatewayZone0 && servicePort.Name == "tcp":
			service.Spec.Ports[i].NodePort = nodePortIstioIngressGatewayZone0
			ip = r.Zone0IP
		case key == keyIstioIngressGatewayZone1 && servicePort.Name == "tcp":
			service.Spec.Ports[i].NodePort = nodePortIstioIngressGatewayZone1
			ip = r.Zone1IP
		case key == keyIstioIngressGatewayZone2 && servicePort.Name == "tcp":
			service.Spec.Ports[i].NodePort = nodePortIstioIngressGatewayZone2
			ip = r.Zone2IP
		}
	}

	if err := r.Client.Patch(ctx, service, patch); err != nil {
		if apierrors.IsInvalid(err) && strings.Contains(err.Error(), "port is already allocated") {
			log.Info("Patching nodePort failed because it is already allocated, enabling auto-remediation")
			return reconcile.Result{Requeue: true}, r.remediateAllocatedNodePorts(ctx, log)
		}
		return reconcile.Result{}, err
	}

	patch = client.MergeFrom(service.DeepCopy())
	service.Status.LoadBalancer.Ingress = []corev1.LoadBalancerIngress{}

	nodes := &corev1.NodeList{}
	if err := r.Client.List(ctx, nodes); err != nil {
		return reconcile.Result{}, err
	}

	for _, node := range nodes.Items {
		for _, address := range node.Status.Addresses {
			if address.Type == corev1.NodeInternalIP {
				service.Status.LoadBalancer.Ingress = append(service.Status.LoadBalancer.Ingress, corev1.LoadBalancerIngress{IP: address.Address})
			}
		}
	}

	// We append this IP behind the nodeIPs because the gardenlet puts the last ingress IP into the DNSRecords
	// and we need a change in the DNSRecords if the control-plane moved from one zone to another to update the coredns-custom configMap.
	service.Status.LoadBalancer.Ingress = append(service.Status.LoadBalancer.Ingress, corev1.LoadBalancerIngress{IP: ip})

	return reconcile.Result{}, r.Client.Status().Patch(ctx, service, patch)
}

func (r *Reconciler) remediateAllocatedNodePorts(ctx context.Context, log logr.Logger) error {
	serviceList := &corev1.ServiceList{}
	if err := r.Client.List(ctx, serviceList); err != nil {
		return err
	}

	for _, service := range serviceList.Items {
		var (
			mustUpdate bool
			patch      = client.StrategicMergeFrom(service.DeepCopy())
		)

		for i, port := range service.Spec.Ports {
			if port.NodePort == nodePortIstioIngressGateway ||
				port.NodePort == nodePortIstioIngressGatewayZone0 ||
				port.NodePort == nodePortIstioIngressGatewayZone1 ||
				port.NodePort == nodePortIstioIngressGatewayZone2 {
				var (
					min, max    = 30000, 32767
					newNodePort = int32(rand.Intn(max-min) + min)
				)

				log.Info("Assigning new nodePort to service which already allocates the nodePort",
					"service", client.ObjectKeyFromObject(&service),
					"newNodePort", newNodePort,
				)

				service.Spec.Ports[i].NodePort = newNodePort
				mustUpdate = true
			}

			if mustUpdate {
				if err := r.Client.Patch(ctx, &service, patch); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
