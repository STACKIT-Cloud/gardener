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

// Code generated by conversion-gen. DO NOT EDIT.

package v1alpha1

import (
	unsafe "unsafe"

	authentication "github.com/gardener/gardener/pkg/apis/authentication"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*AdminKubeconfigRequest)(nil), (*authentication.AdminKubeconfigRequest)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_AdminKubeconfigRequest_To_authentication_AdminKubeconfigRequest(a.(*AdminKubeconfigRequest), b.(*authentication.AdminKubeconfigRequest), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authentication.AdminKubeconfigRequest)(nil), (*AdminKubeconfigRequest)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authentication_AdminKubeconfigRequest_To_v1alpha1_AdminKubeconfigRequest(a.(*authentication.AdminKubeconfigRequest), b.(*AdminKubeconfigRequest), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*AdminKubeconfigRequestSpec)(nil), (*authentication.AdminKubeconfigRequestSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_AdminKubeconfigRequestSpec_To_authentication_AdminKubeconfigRequestSpec(a.(*AdminKubeconfigRequestSpec), b.(*authentication.AdminKubeconfigRequestSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authentication.AdminKubeconfigRequestSpec)(nil), (*AdminKubeconfigRequestSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authentication_AdminKubeconfigRequestSpec_To_v1alpha1_AdminKubeconfigRequestSpec(a.(*authentication.AdminKubeconfigRequestSpec), b.(*AdminKubeconfigRequestSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*AdminKubeconfigRequestStatus)(nil), (*authentication.AdminKubeconfigRequestStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_AdminKubeconfigRequestStatus_To_authentication_AdminKubeconfigRequestStatus(a.(*AdminKubeconfigRequestStatus), b.(*authentication.AdminKubeconfigRequestStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authentication.AdminKubeconfigRequestStatus)(nil), (*AdminKubeconfigRequestStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authentication_AdminKubeconfigRequestStatus_To_v1alpha1_AdminKubeconfigRequestStatus(a.(*authentication.AdminKubeconfigRequestStatus), b.(*AdminKubeconfigRequestStatus), scope)
	}); err != nil {
		return err
	}
	return nil
}

func autoConvert_v1alpha1_AdminKubeconfigRequest_To_authentication_AdminKubeconfigRequest(in *AdminKubeconfigRequest, out *authentication.AdminKubeconfigRequest, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1alpha1_AdminKubeconfigRequestSpec_To_authentication_AdminKubeconfigRequestSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_v1alpha1_AdminKubeconfigRequestStatus_To_authentication_AdminKubeconfigRequestStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha1_AdminKubeconfigRequest_To_authentication_AdminKubeconfigRequest is an autogenerated conversion function.
func Convert_v1alpha1_AdminKubeconfigRequest_To_authentication_AdminKubeconfigRequest(in *AdminKubeconfigRequest, out *authentication.AdminKubeconfigRequest, s conversion.Scope) error {
	return autoConvert_v1alpha1_AdminKubeconfigRequest_To_authentication_AdminKubeconfigRequest(in, out, s)
}

func autoConvert_authentication_AdminKubeconfigRequest_To_v1alpha1_AdminKubeconfigRequest(in *authentication.AdminKubeconfigRequest, out *AdminKubeconfigRequest, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_authentication_AdminKubeconfigRequestSpec_To_v1alpha1_AdminKubeconfigRequestSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_authentication_AdminKubeconfigRequestStatus_To_v1alpha1_AdminKubeconfigRequestStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_authentication_AdminKubeconfigRequest_To_v1alpha1_AdminKubeconfigRequest is an autogenerated conversion function.
func Convert_authentication_AdminKubeconfigRequest_To_v1alpha1_AdminKubeconfigRequest(in *authentication.AdminKubeconfigRequest, out *AdminKubeconfigRequest, s conversion.Scope) error {
	return autoConvert_authentication_AdminKubeconfigRequest_To_v1alpha1_AdminKubeconfigRequest(in, out, s)
}

func autoConvert_v1alpha1_AdminKubeconfigRequestSpec_To_authentication_AdminKubeconfigRequestSpec(in *AdminKubeconfigRequestSpec, out *authentication.AdminKubeconfigRequestSpec, s conversion.Scope) error {
	if err := v1.Convert_Pointer_int64_To_int64(&in.ExpirationSeconds, &out.ExpirationSeconds, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha1_AdminKubeconfigRequestSpec_To_authentication_AdminKubeconfigRequestSpec is an autogenerated conversion function.
func Convert_v1alpha1_AdminKubeconfigRequestSpec_To_authentication_AdminKubeconfigRequestSpec(in *AdminKubeconfigRequestSpec, out *authentication.AdminKubeconfigRequestSpec, s conversion.Scope) error {
	return autoConvert_v1alpha1_AdminKubeconfigRequestSpec_To_authentication_AdminKubeconfigRequestSpec(in, out, s)
}

func autoConvert_authentication_AdminKubeconfigRequestSpec_To_v1alpha1_AdminKubeconfigRequestSpec(in *authentication.AdminKubeconfigRequestSpec, out *AdminKubeconfigRequestSpec, s conversion.Scope) error {
	if err := v1.Convert_int64_To_Pointer_int64(&in.ExpirationSeconds, &out.ExpirationSeconds, s); err != nil {
		return err
	}
	return nil
}

// Convert_authentication_AdminKubeconfigRequestSpec_To_v1alpha1_AdminKubeconfigRequestSpec is an autogenerated conversion function.
func Convert_authentication_AdminKubeconfigRequestSpec_To_v1alpha1_AdminKubeconfigRequestSpec(in *authentication.AdminKubeconfigRequestSpec, out *AdminKubeconfigRequestSpec, s conversion.Scope) error {
	return autoConvert_authentication_AdminKubeconfigRequestSpec_To_v1alpha1_AdminKubeconfigRequestSpec(in, out, s)
}

func autoConvert_v1alpha1_AdminKubeconfigRequestStatus_To_authentication_AdminKubeconfigRequestStatus(in *AdminKubeconfigRequestStatus, out *authentication.AdminKubeconfigRequestStatus, s conversion.Scope) error {
	out.Kubeconfig = *(*[]byte)(unsafe.Pointer(&in.Kubeconfig))
	out.ExpirationTimestamp = in.ExpirationTimestamp
	return nil
}

// Convert_v1alpha1_AdminKubeconfigRequestStatus_To_authentication_AdminKubeconfigRequestStatus is an autogenerated conversion function.
func Convert_v1alpha1_AdminKubeconfigRequestStatus_To_authentication_AdminKubeconfigRequestStatus(in *AdminKubeconfigRequestStatus, out *authentication.AdminKubeconfigRequestStatus, s conversion.Scope) error {
	return autoConvert_v1alpha1_AdminKubeconfigRequestStatus_To_authentication_AdminKubeconfigRequestStatus(in, out, s)
}

func autoConvert_authentication_AdminKubeconfigRequestStatus_To_v1alpha1_AdminKubeconfigRequestStatus(in *authentication.AdminKubeconfigRequestStatus, out *AdminKubeconfigRequestStatus, s conversion.Scope) error {
	out.Kubeconfig = *(*[]byte)(unsafe.Pointer(&in.Kubeconfig))
	out.ExpirationTimestamp = in.ExpirationTimestamp
	return nil
}

// Convert_authentication_AdminKubeconfigRequestStatus_To_v1alpha1_AdminKubeconfigRequestStatus is an autogenerated conversion function.
func Convert_authentication_AdminKubeconfigRequestStatus_To_v1alpha1_AdminKubeconfigRequestStatus(in *authentication.AdminKubeconfigRequestStatus, out *AdminKubeconfigRequestStatus, s conversion.Scope) error {
	return autoConvert_authentication_AdminKubeconfigRequestStatus_To_v1alpha1_AdminKubeconfigRequestStatus(in, out, s)
}
