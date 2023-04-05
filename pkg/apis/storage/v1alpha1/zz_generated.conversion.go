//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright The Kubernetes Authors.

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

	v1alpha1 "k8s.io/api/storage/v1alpha1"
	resource "k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
	storage "k8s.io/kubernetes/pkg/apis/storage"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*v1alpha1.CSIStorageCapacity)(nil), (*storage.CSIStorageCapacity)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_CSIStorageCapacity_To_storage_CSIStorageCapacity(a.(*v1alpha1.CSIStorageCapacity), b.(*storage.CSIStorageCapacity), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*storage.CSIStorageCapacity)(nil), (*v1alpha1.CSIStorageCapacity)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_storage_CSIStorageCapacity_To_v1alpha1_CSIStorageCapacity(a.(*storage.CSIStorageCapacity), b.(*v1alpha1.CSIStorageCapacity), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1alpha1.CSIStorageCapacityList)(nil), (*storage.CSIStorageCapacityList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_CSIStorageCapacityList_To_storage_CSIStorageCapacityList(a.(*v1alpha1.CSIStorageCapacityList), b.(*storage.CSIStorageCapacityList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*storage.CSIStorageCapacityList)(nil), (*v1alpha1.CSIStorageCapacityList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_storage_CSIStorageCapacityList_To_v1alpha1_CSIStorageCapacityList(a.(*storage.CSIStorageCapacityList), b.(*v1alpha1.CSIStorageCapacityList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1alpha1.VolumeError)(nil), (*storage.VolumeError)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_VolumeError_To_storage_VolumeError(a.(*v1alpha1.VolumeError), b.(*storage.VolumeError), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*storage.VolumeError)(nil), (*v1alpha1.VolumeError)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_storage_VolumeError_To_v1alpha1_VolumeError(a.(*storage.VolumeError), b.(*v1alpha1.VolumeError), scope)
	}); err != nil {
		return err
	}
	return nil
}

func autoConvert_v1alpha1_CSIStorageCapacity_To_storage_CSIStorageCapacity(in *v1alpha1.CSIStorageCapacity, out *storage.CSIStorageCapacity, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	out.NodeTopology = (*v1.LabelSelector)(unsafe.Pointer(in.NodeTopology))
	out.StorageClassName = in.StorageClassName
	out.Capacity = (*resource.Quantity)(unsafe.Pointer(in.Capacity))
	out.MaximumVolumeSize = (*resource.Quantity)(unsafe.Pointer(in.MaximumVolumeSize))
	return nil
}

// Convert_v1alpha1_CSIStorageCapacity_To_storage_CSIStorageCapacity is an autogenerated conversion function.
func Convert_v1alpha1_CSIStorageCapacity_To_storage_CSIStorageCapacity(in *v1alpha1.CSIStorageCapacity, out *storage.CSIStorageCapacity, s conversion.Scope) error {
	return autoConvert_v1alpha1_CSIStorageCapacity_To_storage_CSIStorageCapacity(in, out, s)
}

func autoConvert_storage_CSIStorageCapacity_To_v1alpha1_CSIStorageCapacity(in *storage.CSIStorageCapacity, out *v1alpha1.CSIStorageCapacity, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	out.NodeTopology = (*v1.LabelSelector)(unsafe.Pointer(in.NodeTopology))
	out.StorageClassName = in.StorageClassName
	out.Capacity = (*resource.Quantity)(unsafe.Pointer(in.Capacity))
	out.MaximumVolumeSize = (*resource.Quantity)(unsafe.Pointer(in.MaximumVolumeSize))
	return nil
}

// Convert_storage_CSIStorageCapacity_To_v1alpha1_CSIStorageCapacity is an autogenerated conversion function.
func Convert_storage_CSIStorageCapacity_To_v1alpha1_CSIStorageCapacity(in *storage.CSIStorageCapacity, out *v1alpha1.CSIStorageCapacity, s conversion.Scope) error {
	return autoConvert_storage_CSIStorageCapacity_To_v1alpha1_CSIStorageCapacity(in, out, s)
}

func autoConvert_v1alpha1_CSIStorageCapacityList_To_storage_CSIStorageCapacityList(in *v1alpha1.CSIStorageCapacityList, out *storage.CSIStorageCapacityList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]storage.CSIStorageCapacity)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_v1alpha1_CSIStorageCapacityList_To_storage_CSIStorageCapacityList is an autogenerated conversion function.
func Convert_v1alpha1_CSIStorageCapacityList_To_storage_CSIStorageCapacityList(in *v1alpha1.CSIStorageCapacityList, out *storage.CSIStorageCapacityList, s conversion.Scope) error {
	return autoConvert_v1alpha1_CSIStorageCapacityList_To_storage_CSIStorageCapacityList(in, out, s)
}

func autoConvert_storage_CSIStorageCapacityList_To_v1alpha1_CSIStorageCapacityList(in *storage.CSIStorageCapacityList, out *v1alpha1.CSIStorageCapacityList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]v1alpha1.CSIStorageCapacity)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_storage_CSIStorageCapacityList_To_v1alpha1_CSIStorageCapacityList is an autogenerated conversion function.
func Convert_storage_CSIStorageCapacityList_To_v1alpha1_CSIStorageCapacityList(in *storage.CSIStorageCapacityList, out *v1alpha1.CSIStorageCapacityList, s conversion.Scope) error {
	return autoConvert_storage_CSIStorageCapacityList_To_v1alpha1_CSIStorageCapacityList(in, out, s)
}

func autoConvert_v1alpha1_VolumeError_To_storage_VolumeError(in *v1alpha1.VolumeError, out *storage.VolumeError, s conversion.Scope) error {
	out.Time = in.Time
	out.Message = in.Message
	return nil
}

// Convert_v1alpha1_VolumeError_To_storage_VolumeError is an autogenerated conversion function.
func Convert_v1alpha1_VolumeError_To_storage_VolumeError(in *v1alpha1.VolumeError, out *storage.VolumeError, s conversion.Scope) error {
	return autoConvert_v1alpha1_VolumeError_To_storage_VolumeError(in, out, s)
}

func autoConvert_storage_VolumeError_To_v1alpha1_VolumeError(in *storage.VolumeError, out *v1alpha1.VolumeError, s conversion.Scope) error {
	out.Time = in.Time
	out.Message = in.Message
	return nil
}

// Convert_storage_VolumeError_To_v1alpha1_VolumeError is an autogenerated conversion function.
func Convert_storage_VolumeError_To_v1alpha1_VolumeError(in *storage.VolumeError, out *v1alpha1.VolumeError, s conversion.Scope) error {
	return autoConvert_storage_VolumeError_To_v1alpha1_VolumeError(in, out, s)
}
