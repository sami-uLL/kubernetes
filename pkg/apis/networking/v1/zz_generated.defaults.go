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

// Code generated by defaulter-gen. DO NOT EDIT.

package v1

import (
	"encoding/json"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/networking/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// RegisterDefaults adds defaulters functions to the given scheme.
// Public to allow building arbitrary schemes.
// All generated defaulters are covering - they call all nested defaulters.
func RegisterDefaults(scheme *runtime.Scheme) error {
	scheme.AddTypeDefaultingFunc(&v1.IngressClass{}, func(obj interface{}) { SetObjectDefaults_IngressClass(obj.(*v1.IngressClass)) })
	scheme.AddTypeDefaultingFunc(&v1.IngressClassList{}, func(obj interface{}) { SetObjectDefaults_IngressClassList(obj.(*v1.IngressClassList)) })
	scheme.AddTypeDefaultingFunc(&v1.NetworkPolicy{}, func(obj interface{}) { SetObjectDefaults_NetworkPolicy(obj.(*v1.NetworkPolicy)) })
	scheme.AddTypeDefaultingFunc(&v1.NetworkPolicyList{}, func(obj interface{}) { SetObjectDefaults_NetworkPolicyList(obj.(*v1.NetworkPolicyList)) })
	return nil
}

func SetObjectDefaults_IngressClass(in *v1.IngressClass) {
	SetDefaults_IngressClass(in)
	if in.Spec.Parameters != nil {
		if in.Spec.Parameters.Scope == nil {
			ptrVar1 := string(v1.IngressClassParametersReferenceScopeCluster)
			in.Spec.Parameters.Scope = &ptrVar1
		}
	}
}

func SetObjectDefaults_IngressClassList(in *v1.IngressClassList) {
	for i := range in.Items {
		a := &in.Items[i]
		SetObjectDefaults_IngressClass(a)
	}
}

func SetObjectDefaults_NetworkPolicy(in *v1.NetworkPolicy) {
	SetDefaults_NetworkPolicy(in)
	for i := range in.Spec.Ingress {
		a := &in.Spec.Ingress[i]
		for j := range a.Ports {
			b := &a.Ports[j]
			SetDefaults_NetworkPolicyPort(b)
			if b.Protocol == nil {
				ptrVar1 := corev1.Protocol(corev1.ProtocolTCP)
				b.Protocol = &ptrVar1
			}
		}
	}
	for i := range in.Spec.Egress {
		a := &in.Spec.Egress[i]
		for j := range a.Ports {
			b := &a.Ports[j]
			SetDefaults_NetworkPolicyPort(b)
			if b.Protocol == nil {
				ptrVar1 := corev1.Protocol(corev1.ProtocolTCP)
				b.Protocol = &ptrVar1
			}
		}
	}
	if in.Spec.PolicyTypes == nil {
		if err := json.Unmarshal([]byte(`v1.DefaultPolicyTypes()`), &in.Spec.PolicyTypes); err != nil {
			panic(err)
		}
	}
}

func SetObjectDefaults_NetworkPolicyList(in *v1.NetworkPolicyList) {
	for i := range in.Items {
		a := &in.Items[i]
		SetObjectDefaults_NetworkPolicy(a)
	}
}
