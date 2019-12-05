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

package v1

import (
	unsafe "unsafe"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
	config "k8s.io/apiserver/pkg/apis/config"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*AESConfiguration)(nil), (*config.AESConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_AESConfiguration_To_config_AESConfiguration(a.(*AESConfiguration), b.(*config.AESConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.AESConfiguration)(nil), (*AESConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_AESConfiguration_To_v1_AESConfiguration(a.(*config.AESConfiguration), b.(*AESConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*EncryptionConfiguration)(nil), (*config.EncryptionConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_EncryptionConfiguration_To_config_EncryptionConfiguration(a.(*EncryptionConfiguration), b.(*config.EncryptionConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.EncryptionConfiguration)(nil), (*EncryptionConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_EncryptionConfiguration_To_v1_EncryptionConfiguration(a.(*config.EncryptionConfiguration), b.(*EncryptionConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*IdentityConfiguration)(nil), (*config.IdentityConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_IdentityConfiguration_To_config_IdentityConfiguration(a.(*IdentityConfiguration), b.(*config.IdentityConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.IdentityConfiguration)(nil), (*IdentityConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_IdentityConfiguration_To_v1_IdentityConfiguration(a.(*config.IdentityConfiguration), b.(*IdentityConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*KMSConfiguration)(nil), (*config.KMSConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_KMSConfiguration_To_config_KMSConfiguration(a.(*KMSConfiguration), b.(*config.KMSConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.KMSConfiguration)(nil), (*KMSConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_KMSConfiguration_To_v1_KMSConfiguration(a.(*config.KMSConfiguration), b.(*KMSConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*Key)(nil), (*config.Key)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_Key_To_config_Key(a.(*Key), b.(*config.Key), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.Key)(nil), (*Key)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_Key_To_v1_Key(a.(*config.Key), b.(*Key), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ProviderConfiguration)(nil), (*config.ProviderConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ProviderConfiguration_To_config_ProviderConfiguration(a.(*ProviderConfiguration), b.(*config.ProviderConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.ProviderConfiguration)(nil), (*ProviderConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_ProviderConfiguration_To_v1_ProviderConfiguration(a.(*config.ProviderConfiguration), b.(*ProviderConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ResourceConfiguration)(nil), (*config.ResourceConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ResourceConfiguration_To_config_ResourceConfiguration(a.(*ResourceConfiguration), b.(*config.ResourceConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.ResourceConfiguration)(nil), (*ResourceConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_ResourceConfiguration_To_v1_ResourceConfiguration(a.(*config.ResourceConfiguration), b.(*ResourceConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*SecretboxConfiguration)(nil), (*config.SecretboxConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_SecretboxConfiguration_To_config_SecretboxConfiguration(a.(*SecretboxConfiguration), b.(*config.SecretboxConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.SecretboxConfiguration)(nil), (*SecretboxConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_SecretboxConfiguration_To_v1_SecretboxConfiguration(a.(*config.SecretboxConfiguration), b.(*SecretboxConfiguration), scope)
	}); err != nil {
		return err
	}
	return nil
}

func autoConvert_v1_AESConfiguration_To_config_AESConfiguration(in *AESConfiguration, out *config.AESConfiguration, s conversion.Scope) error {
	out.Keys = *(*[]config.Key)(unsafe.Pointer(&in.Keys))
	return nil
}

// Convert_v1_AESConfiguration_To_config_AESConfiguration is an autogenerated conversion function.
func Convert_v1_AESConfiguration_To_config_AESConfiguration(in *AESConfiguration, out *config.AESConfiguration, s conversion.Scope) error {
	return autoConvert_v1_AESConfiguration_To_config_AESConfiguration(in, out, s)
}

func autoConvert_config_AESConfiguration_To_v1_AESConfiguration(in *config.AESConfiguration, out *AESConfiguration, s conversion.Scope) error {
	out.Keys = *(*[]Key)(unsafe.Pointer(&in.Keys))
	return nil
}

// Convert_config_AESConfiguration_To_v1_AESConfiguration is an autogenerated conversion function.
func Convert_config_AESConfiguration_To_v1_AESConfiguration(in *config.AESConfiguration, out *AESConfiguration, s conversion.Scope) error {
	return autoConvert_config_AESConfiguration_To_v1_AESConfiguration(in, out, s)
}

func autoConvert_v1_EncryptionConfiguration_To_config_EncryptionConfiguration(in *EncryptionConfiguration, out *config.EncryptionConfiguration, s conversion.Scope) error {
	out.Resources = *(*[]config.ResourceConfiguration)(unsafe.Pointer(&in.Resources))
	return nil
}

// Convert_v1_EncryptionConfiguration_To_config_EncryptionConfiguration is an autogenerated conversion function.
func Convert_v1_EncryptionConfiguration_To_config_EncryptionConfiguration(in *EncryptionConfiguration, out *config.EncryptionConfiguration, s conversion.Scope) error {
	return autoConvert_v1_EncryptionConfiguration_To_config_EncryptionConfiguration(in, out, s)
}

func autoConvert_config_EncryptionConfiguration_To_v1_EncryptionConfiguration(in *config.EncryptionConfiguration, out *EncryptionConfiguration, s conversion.Scope) error {
	out.Resources = *(*[]ResourceConfiguration)(unsafe.Pointer(&in.Resources))
	return nil
}

// Convert_config_EncryptionConfiguration_To_v1_EncryptionConfiguration is an autogenerated conversion function.
func Convert_config_EncryptionConfiguration_To_v1_EncryptionConfiguration(in *config.EncryptionConfiguration, out *EncryptionConfiguration, s conversion.Scope) error {
	return autoConvert_config_EncryptionConfiguration_To_v1_EncryptionConfiguration(in, out, s)
}

func autoConvert_v1_IdentityConfiguration_To_config_IdentityConfiguration(in *IdentityConfiguration, out *config.IdentityConfiguration, s conversion.Scope) error {
	return nil
}

// Convert_v1_IdentityConfiguration_To_config_IdentityConfiguration is an autogenerated conversion function.
func Convert_v1_IdentityConfiguration_To_config_IdentityConfiguration(in *IdentityConfiguration, out *config.IdentityConfiguration, s conversion.Scope) error {
	return autoConvert_v1_IdentityConfiguration_To_config_IdentityConfiguration(in, out, s)
}

func autoConvert_config_IdentityConfiguration_To_v1_IdentityConfiguration(in *config.IdentityConfiguration, out *IdentityConfiguration, s conversion.Scope) error {
	return nil
}

// Convert_config_IdentityConfiguration_To_v1_IdentityConfiguration is an autogenerated conversion function.
func Convert_config_IdentityConfiguration_To_v1_IdentityConfiguration(in *config.IdentityConfiguration, out *IdentityConfiguration, s conversion.Scope) error {
	return autoConvert_config_IdentityConfiguration_To_v1_IdentityConfiguration(in, out, s)
}

func autoConvert_v1_KMSConfiguration_To_config_KMSConfiguration(in *KMSConfiguration, out *config.KMSConfiguration, s conversion.Scope) error {
	out.Name = in.Name
	out.CacheSize = (*int32)(unsafe.Pointer(in.CacheSize))
	out.Endpoint = in.Endpoint
	out.Timeout = (*metav1.Duration)(unsafe.Pointer(in.Timeout))
	out.DataEncryptionAlgorithm = config.EnvelopeDataEncryptionTransformer(in.DataEncryptionAlgorithm)
	return nil
}

// Convert_v1_KMSConfiguration_To_config_KMSConfiguration is an autogenerated conversion function.
func Convert_v1_KMSConfiguration_To_config_KMSConfiguration(in *KMSConfiguration, out *config.KMSConfiguration, s conversion.Scope) error {
	return autoConvert_v1_KMSConfiguration_To_config_KMSConfiguration(in, out, s)
}

func autoConvert_config_KMSConfiguration_To_v1_KMSConfiguration(in *config.KMSConfiguration, out *KMSConfiguration, s conversion.Scope) error {
	out.Name = in.Name
	out.CacheSize = (*int32)(unsafe.Pointer(in.CacheSize))
	out.Endpoint = in.Endpoint
	out.Timeout = (*metav1.Duration)(unsafe.Pointer(in.Timeout))
	out.DataEncryptionAlgorithm = EnvelopeDataEncryptionTransformer(in.DataEncryptionAlgorithm)
	return nil
}

// Convert_config_KMSConfiguration_To_v1_KMSConfiguration is an autogenerated conversion function.
func Convert_config_KMSConfiguration_To_v1_KMSConfiguration(in *config.KMSConfiguration, out *KMSConfiguration, s conversion.Scope) error {
	return autoConvert_config_KMSConfiguration_To_v1_KMSConfiguration(in, out, s)
}

func autoConvert_v1_Key_To_config_Key(in *Key, out *config.Key, s conversion.Scope) error {
	out.Name = in.Name
	out.Secret = in.Secret
	return nil
}

// Convert_v1_Key_To_config_Key is an autogenerated conversion function.
func Convert_v1_Key_To_config_Key(in *Key, out *config.Key, s conversion.Scope) error {
	return autoConvert_v1_Key_To_config_Key(in, out, s)
}

func autoConvert_config_Key_To_v1_Key(in *config.Key, out *Key, s conversion.Scope) error {
	out.Name = in.Name
	out.Secret = in.Secret
	return nil
}

// Convert_config_Key_To_v1_Key is an autogenerated conversion function.
func Convert_config_Key_To_v1_Key(in *config.Key, out *Key, s conversion.Scope) error {
	return autoConvert_config_Key_To_v1_Key(in, out, s)
}

func autoConvert_v1_ProviderConfiguration_To_config_ProviderConfiguration(in *ProviderConfiguration, out *config.ProviderConfiguration, s conversion.Scope) error {
	out.AESGCM = (*config.AESConfiguration)(unsafe.Pointer(in.AESGCM))
	out.AESCBC = (*config.AESConfiguration)(unsafe.Pointer(in.AESCBC))
	out.Secretbox = (*config.SecretboxConfiguration)(unsafe.Pointer(in.Secretbox))
	out.Identity = (*config.IdentityConfiguration)(unsafe.Pointer(in.Identity))
	out.KMS = (*config.KMSConfiguration)(unsafe.Pointer(in.KMS))
	return nil
}

// Convert_v1_ProviderConfiguration_To_config_ProviderConfiguration is an autogenerated conversion function.
func Convert_v1_ProviderConfiguration_To_config_ProviderConfiguration(in *ProviderConfiguration, out *config.ProviderConfiguration, s conversion.Scope) error {
	return autoConvert_v1_ProviderConfiguration_To_config_ProviderConfiguration(in, out, s)
}

func autoConvert_config_ProviderConfiguration_To_v1_ProviderConfiguration(in *config.ProviderConfiguration, out *ProviderConfiguration, s conversion.Scope) error {
	out.AESGCM = (*AESConfiguration)(unsafe.Pointer(in.AESGCM))
	out.AESCBC = (*AESConfiguration)(unsafe.Pointer(in.AESCBC))
	out.Secretbox = (*SecretboxConfiguration)(unsafe.Pointer(in.Secretbox))
	out.Identity = (*IdentityConfiguration)(unsafe.Pointer(in.Identity))
	out.KMS = (*KMSConfiguration)(unsafe.Pointer(in.KMS))
	return nil
}

// Convert_config_ProviderConfiguration_To_v1_ProviderConfiguration is an autogenerated conversion function.
func Convert_config_ProviderConfiguration_To_v1_ProviderConfiguration(in *config.ProviderConfiguration, out *ProviderConfiguration, s conversion.Scope) error {
	return autoConvert_config_ProviderConfiguration_To_v1_ProviderConfiguration(in, out, s)
}

func autoConvert_v1_ResourceConfiguration_To_config_ResourceConfiguration(in *ResourceConfiguration, out *config.ResourceConfiguration, s conversion.Scope) error {
	out.Resources = *(*[]string)(unsafe.Pointer(&in.Resources))
	out.Providers = *(*[]config.ProviderConfiguration)(unsafe.Pointer(&in.Providers))
	return nil
}

// Convert_v1_ResourceConfiguration_To_config_ResourceConfiguration is an autogenerated conversion function.
func Convert_v1_ResourceConfiguration_To_config_ResourceConfiguration(in *ResourceConfiguration, out *config.ResourceConfiguration, s conversion.Scope) error {
	return autoConvert_v1_ResourceConfiguration_To_config_ResourceConfiguration(in, out, s)
}

func autoConvert_config_ResourceConfiguration_To_v1_ResourceConfiguration(in *config.ResourceConfiguration, out *ResourceConfiguration, s conversion.Scope) error {
	out.Resources = *(*[]string)(unsafe.Pointer(&in.Resources))
	out.Providers = *(*[]ProviderConfiguration)(unsafe.Pointer(&in.Providers))
	return nil
}

// Convert_config_ResourceConfiguration_To_v1_ResourceConfiguration is an autogenerated conversion function.
func Convert_config_ResourceConfiguration_To_v1_ResourceConfiguration(in *config.ResourceConfiguration, out *ResourceConfiguration, s conversion.Scope) error {
	return autoConvert_config_ResourceConfiguration_To_v1_ResourceConfiguration(in, out, s)
}

func autoConvert_v1_SecretboxConfiguration_To_config_SecretboxConfiguration(in *SecretboxConfiguration, out *config.SecretboxConfiguration, s conversion.Scope) error {
	out.Keys = *(*[]config.Key)(unsafe.Pointer(&in.Keys))
	return nil
}

// Convert_v1_SecretboxConfiguration_To_config_SecretboxConfiguration is an autogenerated conversion function.
func Convert_v1_SecretboxConfiguration_To_config_SecretboxConfiguration(in *SecretboxConfiguration, out *config.SecretboxConfiguration, s conversion.Scope) error {
	return autoConvert_v1_SecretboxConfiguration_To_config_SecretboxConfiguration(in, out, s)
}

func autoConvert_config_SecretboxConfiguration_To_v1_SecretboxConfiguration(in *config.SecretboxConfiguration, out *SecretboxConfiguration, s conversion.Scope) error {
	out.Keys = *(*[]Key)(unsafe.Pointer(&in.Keys))
	return nil
}

// Convert_config_SecretboxConfiguration_To_v1_SecretboxConfiguration is an autogenerated conversion function.
func Convert_config_SecretboxConfiguration_To_v1_SecretboxConfiguration(in *config.SecretboxConfiguration, out *SecretboxConfiguration, s conversion.Scope) error {
	return autoConvert_config_SecretboxConfiguration_To_v1_SecretboxConfiguration(in, out, s)
}
