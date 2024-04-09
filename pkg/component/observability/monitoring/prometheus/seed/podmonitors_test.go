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

package seed_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gardener/gardener/pkg/component/observability/monitoring/prometheus/seed"
)

var _ = Describe("PodMonitors", func() {
	Describe("#CentralPodMonitors", func() {
		It("should return the expected objects", func() {
			Expect(seed.CentralPodMonitors()).To(HaveExactElements(
				&monitoringv1.PodMonitor{
					ObjectMeta: metav1.ObjectMeta{
						Name: "extensions",
					},
					Spec: monitoringv1.PodMonitorSpec{
						NamespaceSelector: monitoringv1.NamespaceSelector{Any: true},
						PodMetricsEndpoints: []monitoringv1.PodMetricsEndpoint{{
							RelabelConfigs: []*monitoringv1.RelabelConfig{
								{
									SourceLabels: []monitoringv1.LabelName{
										"__meta_kubernetes_namespace",
										"__meta_kubernetes_pod_annotation_prometheus_io_scrape",
										"__meta_kubernetes_pod_annotation_prometheus_io_port",
									},
									Regex:  `extension-(.+);true;(.+)`,
									Action: "keep",
								},
								{
									SourceLabels: []monitoringv1.LabelName{"__meta_kubernetes_pod_annotation_prometheus_io_name"},
									Regex:        `(.+)`,
									Action:       "replace",
									TargetLabel:  "job",
								},
								{
									SourceLabels: []monitoringv1.LabelName{"__address__", "__meta_kubernetes_pod_annotation_prometheus_io_port"},
									Regex:        `([^:]+)(?::\d+)?;(\d+)`,
									Action:       "replace",
									Replacement:  `$1:$2`,
									TargetLabel:  "__address__",
								},
								{
									Action: "labelmap",
									Regex:  `__meta_kubernetes_pod_label_(.+)`,
								},
							},
						}},
					},
				},
				&monitoringv1.PodMonitor{
					ObjectMeta: metav1.ObjectMeta{
						Name: "garden",
					},
					Spec: monitoringv1.PodMonitorSpec{
						PodMetricsEndpoints: []monitoringv1.PodMetricsEndpoint{{
							Scheme:    "https",
							TLSConfig: &monitoringv1.SafeTLSConfig{InsecureSkipVerify: true},
							RelabelConfigs: []*monitoringv1.RelabelConfig{
								{
									SourceLabels: []monitoringv1.LabelName{
										"__meta_kubernetes_pod_annotation_prometheus_io_scrape",
										"__meta_kubernetes_pod_annotation_prometheus_io_port",
									},
									Regex:  `true;(.+)`,
									Action: "keep",
								},
								{
									SourceLabels: []monitoringv1.LabelName{"__meta_kubernetes_pod_annotation_prometheus_io_name"},
									Regex:        `(.+)`,
									Action:       "replace",
									TargetLabel:  "job",
								},
								{
									SourceLabels: []monitoringv1.LabelName{"__meta_kubernetes_pod_annotation_prometheus_io_scheme"},
									Regex:        `(https?)`,
									Action:       "replace",
									TargetLabel:  "__scheme__",
								},
								{
									SourceLabels: []monitoringv1.LabelName{"__address__", "__meta_kubernetes_pod_annotation_prometheus_io_port"},
									Regex:        `([^:]+)(?::\d+)?;(\d+)`,
									Replacement:  `$1:$2`,
									Action:       "replace",
									TargetLabel:  "__address__",
								},
								{
									Action: "labelmap",
									Regex:  `__meta_kubernetes_pod_label_(.+)`,
								},
							},
						}},
					},
				},
			))
		})
	})
})
