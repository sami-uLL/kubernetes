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

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1beta1

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubeletAnonymousAuthentication) DeepCopyInto(out *KubeletAnonymousAuthentication) {
	*out = *in
	if in.Enabled != nil {
		in, out := &in.Enabled, &out.Enabled
		*out = new(bool)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubeletAnonymousAuthentication.
func (in *KubeletAnonymousAuthentication) DeepCopy() *KubeletAnonymousAuthentication {
	if in == nil {
		return nil
	}
	out := new(KubeletAnonymousAuthentication)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubeletAuthentication) DeepCopyInto(out *KubeletAuthentication) {
	*out = *in
	out.X509 = in.X509
	in.Webhook.DeepCopyInto(&out.Webhook)
	in.Anonymous.DeepCopyInto(&out.Anonymous)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubeletAuthentication.
func (in *KubeletAuthentication) DeepCopy() *KubeletAuthentication {
	if in == nil {
		return nil
	}
	out := new(KubeletAuthentication)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubeletAuthorization) DeepCopyInto(out *KubeletAuthorization) {
	*out = *in
	out.Webhook = in.Webhook
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubeletAuthorization.
func (in *KubeletAuthorization) DeepCopy() *KubeletAuthorization {
	if in == nil {
		return nil
	}
	out := new(KubeletAuthorization)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubeletConfiguration) DeepCopyInto(out *KubeletConfiguration) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	if in.EnableServer != nil {
		in, out := &in.EnableServer, &out.EnableServer
		*out = new(bool)
		**out = **in
	}
	out.SyncFrequency = in.SyncFrequency
	out.FileCheckFrequency = in.FileCheckFrequency
	out.HTTPCheckFrequency = in.HTTPCheckFrequency
	if in.StaticPodURLHeader != nil {
		in, out := &in.StaticPodURLHeader, &out.StaticPodURLHeader
		*out = make(map[string][]string, len(*in))
		for key, val := range *in {
			var outVal []string
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = make([]string, len(*in))
				copy(*out, *in)
			}
			(*out)[key] = outVal
		}
	}
	if in.TLSCipherSuites != nil {
		in, out := &in.TLSCipherSuites, &out.TLSCipherSuites
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	in.Authentication.DeepCopyInto(&out.Authentication)
	out.Authorization = in.Authorization
	if in.RegistryPullQPS != nil {
		in, out := &in.RegistryPullQPS, &out.RegistryPullQPS
		*out = new(int32)
		**out = **in
	}
	if in.EventRecordQPS != nil {
		in, out := &in.EventRecordQPS, &out.EventRecordQPS
		*out = new(int32)
		**out = **in
	}
	if in.EnableDebuggingHandlers != nil {
		in, out := &in.EnableDebuggingHandlers, &out.EnableDebuggingHandlers
		*out = new(bool)
		**out = **in
	}
	if in.HealthzPort != nil {
		in, out := &in.HealthzPort, &out.HealthzPort
		*out = new(int32)
		**out = **in
	}
	if in.OOMScoreAdj != nil {
		in, out := &in.OOMScoreAdj, &out.OOMScoreAdj
		*out = new(int32)
		**out = **in
	}
	if in.ClusterDNS != nil {
		in, out := &in.ClusterDNS, &out.ClusterDNS
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	out.StreamingConnectionIdleTimeout = in.StreamingConnectionIdleTimeout
	out.NodeStatusUpdateFrequency = in.NodeStatusUpdateFrequency
	out.NodeStatusReportFrequency = in.NodeStatusReportFrequency
	out.ImageMinimumGCAge = in.ImageMinimumGCAge
	if in.ImageGCHighThresholdPercent != nil {
		in, out := &in.ImageGCHighThresholdPercent, &out.ImageGCHighThresholdPercent
		*out = new(int32)
		**out = **in
	}
	if in.ImageGCLowThresholdPercent != nil {
		in, out := &in.ImageGCLowThresholdPercent, &out.ImageGCLowThresholdPercent
		*out = new(int32)
		**out = **in
	}
	out.VolumeStatsAggPeriod = in.VolumeStatsAggPeriod
	if in.CgroupsPerQOS != nil {
		in, out := &in.CgroupsPerQOS, &out.CgroupsPerQOS
		*out = new(bool)
		**out = **in
	}
	out.CPUManagerReconcilePeriod = in.CPUManagerReconcilePeriod
	if in.QOSReserved != nil {
		in, out := &in.QOSReserved, &out.QOSReserved
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	out.RuntimeRequestTimeout = in.RuntimeRequestTimeout
	if in.PodPidsLimit != nil {
		in, out := &in.PodPidsLimit, &out.PodPidsLimit
		*out = new(int64)
		**out = **in
	}
	if in.CPUCFSQuota != nil {
		in, out := &in.CPUCFSQuota, &out.CPUCFSQuota
		*out = new(bool)
		**out = **in
	}
	if in.CPUCFSQuotaPeriod != nil {
		in, out := &in.CPUCFSQuotaPeriod, &out.CPUCFSQuotaPeriod
		*out = new(v1.Duration)
		**out = **in
	}
	if in.NodeStatusMaxImages != nil {
		in, out := &in.NodeStatusMaxImages, &out.NodeStatusMaxImages
		*out = new(int32)
		**out = **in
	}
	if in.KubeAPIQPS != nil {
		in, out := &in.KubeAPIQPS, &out.KubeAPIQPS
		*out = new(int32)
		**out = **in
	}
	if in.SerializeImagePulls != nil {
		in, out := &in.SerializeImagePulls, &out.SerializeImagePulls
		*out = new(bool)
		**out = **in
	}
	if in.EvictionHard != nil {
		in, out := &in.EvictionHard, &out.EvictionHard
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.EvictionSoft != nil {
		in, out := &in.EvictionSoft, &out.EvictionSoft
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.EvictionSoftGracePeriod != nil {
		in, out := &in.EvictionSoftGracePeriod, &out.EvictionSoftGracePeriod
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	out.EvictionPressureTransitionPeriod = in.EvictionPressureTransitionPeriod
	if in.EvictionMinimumReclaim != nil {
		in, out := &in.EvictionMinimumReclaim, &out.EvictionMinimumReclaim
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.EnableControllerAttachDetach != nil {
		in, out := &in.EnableControllerAttachDetach, &out.EnableControllerAttachDetach
		*out = new(bool)
		**out = **in
	}
	if in.MakeIPTablesUtilChains != nil {
		in, out := &in.MakeIPTablesUtilChains, &out.MakeIPTablesUtilChains
		*out = new(bool)
		**out = **in
	}
	if in.IPTablesMasqueradeBit != nil {
		in, out := &in.IPTablesMasqueradeBit, &out.IPTablesMasqueradeBit
		*out = new(int32)
		**out = **in
	}
	if in.IPTablesDropBit != nil {
		in, out := &in.IPTablesDropBit, &out.IPTablesDropBit
		*out = new(int32)
		**out = **in
	}
	if in.FeatureGates != nil {
		in, out := &in.FeatureGates, &out.FeatureGates
		*out = make(map[string]bool, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.FailSwapOn != nil {
		in, out := &in.FailSwapOn, &out.FailSwapOn
		*out = new(bool)
		**out = **in
	}
	if in.ContainerLogMaxFiles != nil {
		in, out := &in.ContainerLogMaxFiles, &out.ContainerLogMaxFiles
		*out = new(int32)
		**out = **in
	}
	if in.SystemReserved != nil {
		in, out := &in.SystemReserved, &out.SystemReserved
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.KubeReserved != nil {
		in, out := &in.KubeReserved, &out.KubeReserved
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.EnforceNodeAllocatable != nil {
		in, out := &in.EnforceNodeAllocatable, &out.EnforceNodeAllocatable
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.AllowedUnsafeSysctls != nil {
		in, out := &in.AllowedUnsafeSysctls, &out.AllowedUnsafeSysctls
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	out.Logging = in.Logging
	if in.EnableSystemLogHandler != nil {
		in, out := &in.EnableSystemLogHandler, &out.EnableSystemLogHandler
		*out = new(bool)
		**out = **in
	}
	if in.CadvisorMetrics != nil {
		in, out := &in.CadvisorMetrics, &out.CadvisorMetrics
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubeletConfiguration.
func (in *KubeletConfiguration) DeepCopy() *KubeletConfiguration {
	if in == nil {
		return nil
	}
	out := new(KubeletConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *KubeletConfiguration) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubeletWebhookAuthentication) DeepCopyInto(out *KubeletWebhookAuthentication) {
	*out = *in
	if in.Enabled != nil {
		in, out := &in.Enabled, &out.Enabled
		*out = new(bool)
		**out = **in
	}
	out.CacheTTL = in.CacheTTL
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubeletWebhookAuthentication.
func (in *KubeletWebhookAuthentication) DeepCopy() *KubeletWebhookAuthentication {
	if in == nil {
		return nil
	}
	out := new(KubeletWebhookAuthentication)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubeletWebhookAuthorization) DeepCopyInto(out *KubeletWebhookAuthorization) {
	*out = *in
	out.CacheAuthorizedTTL = in.CacheAuthorizedTTL
	out.CacheUnauthorizedTTL = in.CacheUnauthorizedTTL
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubeletWebhookAuthorization.
func (in *KubeletWebhookAuthorization) DeepCopy() *KubeletWebhookAuthorization {
	if in == nil {
		return nil
	}
	out := new(KubeletWebhookAuthorization)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubeletX509Authentication) DeepCopyInto(out *KubeletX509Authentication) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubeletX509Authentication.
func (in *KubeletX509Authentication) DeepCopy() *KubeletX509Authentication {
	if in == nil {
		return nil
	}
	out := new(KubeletX509Authentication)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SerializedNodeConfigSource) DeepCopyInto(out *SerializedNodeConfigSource) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.Source.DeepCopyInto(&out.Source)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SerializedNodeConfigSource.
func (in *SerializedNodeConfigSource) DeepCopy() *SerializedNodeConfigSource {
	if in == nil {
		return nil
	}
	out := new(SerializedNodeConfigSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SerializedNodeConfigSource) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}
