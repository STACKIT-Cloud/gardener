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

package fluentoperator

import (
	"context"
	_ "embed"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/gardener/gardener/pkg/client/kubernetes"
	"github.com/gardener/gardener/pkg/component"
	"github.com/gardener/gardener/pkg/utils/flow"
)

var (
	//go:embed assets/crd-fluentbit.fluent.io_clusterfilters.yaml
	fluentBitClusterFilterCRD string
	//go:embed assets/crd-fluentbit.fluent.io_clusterfluentbitconfigs.yaml
	fluentBitClusterFBConfigCRD string
	//go:embed assets/crd-fluentbit.fluent.io_clusterinputs.yaml
	fluentBitClusterInputCRD string
	//go:embed assets/crd-fluentbit.fluent.io_clusteroutputs.yaml
	fluentBitClusterOutputCRD string
	//go:embed assets/crd-fluentbit.fluent.io_clusterparsers.yaml
	fluentBitClusterParserCRD string
	//go:embed assets/crd-fluentbit.fluent.io_fluentbits.yaml
	fluentBitCRD string
	//go:embed assets/crd-fluentbit.fluent.io_collectors.yaml
	fluentBitCollectorCRD string
	//go:embed assets/crd-fluentbit.fluent.io_fluentbitconfigs.yaml
	fluentBitConfigCRD string
	//go:embed assets/crd-fluentbit.fluent.io_filters.yaml
	fluentBitFilterCRD string
	//go:embed assets/crd-fluentbit.fluent.io_parsers.yaml
	fluentBitParserCRD string
	//go:embed assets/crd-fluentbit.fluent.io_outputs.yaml
	fluentBitOutputCRD string

	resources []string
)

func init() {
	resources = append(resources,
		fluentBitClusterFilterCRD,
		fluentBitClusterFBConfigCRD,
		fluentBitClusterInputCRD,
		fluentBitClusterOutputCRD,
		fluentBitClusterParserCRD,
		fluentBitCRD,
		fluentBitCollectorCRD,
		fluentBitConfigCRD,
		fluentBitFilterCRD,
		fluentBitParserCRD,
		fluentBitOutputCRD,
	)
}

type crds struct {
	applier kubernetes.Applier
}

// NewCRDs can be used to deploy Fluent Operator CRDs.
func NewCRDs(a kubernetes.Applier) component.DeployWaiter {
	return &crds{
		applier: a,
	}
}

// Deploy creates and updates the CRD definitions for the Fluent Operator.
func (c *crds) Deploy(ctx context.Context) error {
	var fns []flow.TaskFn

	for _, resource := range resources {
		r := resource
		fns = append(fns, func(ctx context.Context) error {
			return c.applier.ApplyManifest(ctx, kubernetes.NewManifestReader([]byte(r)), kubernetes.DefaultMergeFuncs)
		})
	}

	return flow.Parallel(fns...)(ctx)
}

// Destroy deletes the CRDs for the Fluent Operator.
func (c *crds) Destroy(ctx context.Context) error {
	var fns []flow.TaskFn

	for _, resource := range resources {
		r := resource
		fns = append(fns, func(ctx context.Context) error {
			return client.IgnoreNotFound(c.applier.DeleteManifest(ctx, kubernetes.NewManifestReader([]byte(r))))
		})
	}

	return flow.Parallel(fns...)(ctx)
}

// Wait does nothing
func (c *crds) Wait(_ context.Context) error {
	return nil
}

// WaitCleanup does nothing
func (c *crds) WaitCleanup(_ context.Context) error {
	return nil
}
