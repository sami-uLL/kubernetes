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

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1

// RBDPersistentVolumeSourceApplyConfiguration represents an declarative configuration of the RBDPersistentVolumeSource type for use
// with apply.
type RBDPersistentVolumeSourceApplyConfiguration struct {
	Monitors  []string                           `json:"monitors,omitempty"`
	Image     *string                            `json:"image,omitempty"`
	FSType    *string                            `json:"fsType,omitempty"`
	RBDPool   *string                            `json:"pool,omitempty"`
	RadosUser *string                            `json:"user,omitempty"`
	Keyring   *string                            `json:"keyring,omitempty"`
	SecretRef *SecretReferenceApplyConfiguration `json:"secretRef,omitempty"`
	ReadOnly  *bool                              `json:"readOnly,omitempty"`
}

// RBDPersistentVolumeSourceApplyConfiguration constructs an declarative configuration of the RBDPersistentVolumeSource type for use with
// apply.
func RBDPersistentVolumeSource() *RBDPersistentVolumeSourceApplyConfiguration {
	return &RBDPersistentVolumeSourceApplyConfiguration{}
}

// WithMonitors adds the given value to the Monitors field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Monitors field.
func (b *RBDPersistentVolumeSourceApplyConfiguration) WithMonitors(values ...string) *RBDPersistentVolumeSourceApplyConfiguration {
	for i := range values {
		b.Monitors = append(b.Monitors, values[i])
	}
	return b
}

// WithImage sets the Image field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Image field is set to the value of the last call.
func (b *RBDPersistentVolumeSourceApplyConfiguration) WithImage(value string) *RBDPersistentVolumeSourceApplyConfiguration {
	b.Image = &value
	return b
}

// WithFSType sets the FSType field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the FSType field is set to the value of the last call.
func (b *RBDPersistentVolumeSourceApplyConfiguration) WithFSType(value string) *RBDPersistentVolumeSourceApplyConfiguration {
	b.FSType = &value
	return b
}

// WithRBDPool sets the RBDPool field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the RBDPool field is set to the value of the last call.
func (b *RBDPersistentVolumeSourceApplyConfiguration) WithRBDPool(value string) *RBDPersistentVolumeSourceApplyConfiguration {
	b.RBDPool = &value
	return b
}

// WithRadosUser sets the RadosUser field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the RadosUser field is set to the value of the last call.
func (b *RBDPersistentVolumeSourceApplyConfiguration) WithRadosUser(value string) *RBDPersistentVolumeSourceApplyConfiguration {
	b.RadosUser = &value
	return b
}

// WithKeyring sets the Keyring field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Keyring field is set to the value of the last call.
func (b *RBDPersistentVolumeSourceApplyConfiguration) WithKeyring(value string) *RBDPersistentVolumeSourceApplyConfiguration {
	b.Keyring = &value
	return b
}

// WithSecretRef sets the SecretRef field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the SecretRef field is set to the value of the last call.
func (b *RBDPersistentVolumeSourceApplyConfiguration) WithSecretRef(value *SecretReferenceApplyConfiguration) *RBDPersistentVolumeSourceApplyConfiguration {
	b.SecretRef = value
	return b
}

// WithReadOnly sets the ReadOnly field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ReadOnly field is set to the value of the last call.
func (b *RBDPersistentVolumeSourceApplyConfiguration) WithReadOnly(value bool) *RBDPersistentVolumeSourceApplyConfiguration {
	b.ReadOnly = &value
	return b
}
