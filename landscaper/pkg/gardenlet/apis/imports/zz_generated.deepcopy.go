//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright (c) 2022 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by deepcopy-gen. DO NOT EDIT.

package imports

import (
	json "encoding/json"

	seedmanagement "github.com/gardener/gardener/pkg/apis/seedmanagement"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Imports) DeepCopyInto(out *Imports) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.SeedCluster.DeepCopyInto(&out.SeedCluster)
	in.GardenCluster.DeepCopyInto(&out.GardenCluster)
	if in.SeedBackupCredentials != nil {
		in, out := &in.SeedBackupCredentials, &out.SeedBackupCredentials
		*out = new(json.RawMessage)
		if **in != nil {
			in, out := *in, *out
			*out = make([]byte, len(*in))
			copy(*out, *in)
		}
	}
	if in.DeploymentConfiguration != nil {
		in, out := &in.DeploymentConfiguration, &out.DeploymentConfiguration
		*out = new(seedmanagement.GardenletDeployment)
		(*in).DeepCopyInto(*out)
	}
	if in.ComponentConfiguration != nil {
		out.ComponentConfiguration = in.ComponentConfiguration.DeepCopyObject()
	}
	if in.ImageVectorOverwrite != nil {
		in, out := &in.ImageVectorOverwrite, &out.ImageVectorOverwrite
		*out = new(json.RawMessage)
		if **in != nil {
			in, out := *in, *out
			*out = make([]byte, len(*in))
			copy(*out, *in)
		}
	}
	if in.ComponentImageVectorOverwrites != nil {
		in, out := &in.ComponentImageVectorOverwrites, &out.ComponentImageVectorOverwrites
		*out = new(json.RawMessage)
		if **in != nil {
			in, out := *in, *out
			*out = make([]byte, len(*in))
			copy(*out, *in)
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Imports.
func (in *Imports) DeepCopy() *Imports {
	if in == nil {
		return nil
	}
	out := new(Imports)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Imports) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}
