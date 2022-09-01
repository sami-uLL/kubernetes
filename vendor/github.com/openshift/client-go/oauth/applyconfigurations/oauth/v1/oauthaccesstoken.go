// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1

import (
	oauthv1 "github.com/openshift/api/oauth/v1"
	internal "github.com/openshift/client-go/oauth/applyconfigurations/internal"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	managedfields "k8s.io/apimachinery/pkg/util/managedfields"
	v1 "k8s.io/client-go/applyconfigurations/meta/v1"
)

// OAuthAccessTokenApplyConfiguration represents an declarative configuration of the OAuthAccessToken type for use
// with apply.
type OAuthAccessTokenApplyConfiguration struct {
	v1.TypeMetaApplyConfiguration    `json:",inline"`
	*v1.ObjectMetaApplyConfiguration `json:"metadata,omitempty"`
	ClientName                       *string  `json:"clientName,omitempty"`
	ExpiresIn                        *int64   `json:"expiresIn,omitempty"`
	Scopes                           []string `json:"scopes,omitempty"`
	RedirectURI                      *string  `json:"redirectURI,omitempty"`
	UserName                         *string  `json:"userName,omitempty"`
	UserUID                          *string  `json:"userUID,omitempty"`
	AuthorizeToken                   *string  `json:"authorizeToken,omitempty"`
	RefreshToken                     *string  `json:"refreshToken,omitempty"`
	InactivityTimeoutSeconds         *int32   `json:"inactivityTimeoutSeconds,omitempty"`
}

// OAuthAccessToken constructs an declarative configuration of the OAuthAccessToken type for use with
// apply.
func OAuthAccessToken(name string) *OAuthAccessTokenApplyConfiguration {
	b := &OAuthAccessTokenApplyConfiguration{}
	b.WithName(name)
	b.WithKind("OAuthAccessToken")
	b.WithAPIVersion("oauth.openshift.io/v1")
	return b
}

// ExtractOAuthAccessToken extracts the applied configuration owned by fieldManager from
// oAuthAccessToken. If no managedFields are found in oAuthAccessToken for fieldManager, a
// OAuthAccessTokenApplyConfiguration is returned with only the Name, Namespace (if applicable),
// APIVersion and Kind populated. It is possible that no managed fields were found for because other
// field managers have taken ownership of all the fields previously owned by fieldManager, or because
// the fieldManager never owned fields any fields.
// oAuthAccessToken must be a unmodified OAuthAccessToken API object that was retrieved from the Kubernetes API.
// ExtractOAuthAccessToken provides a way to perform a extract/modify-in-place/apply workflow.
// Note that an extracted apply configuration will contain fewer fields than what the fieldManager previously
// applied if another fieldManager has updated or force applied any of the previously applied fields.
// Experimental!
func ExtractOAuthAccessToken(oAuthAccessToken *oauthv1.OAuthAccessToken, fieldManager string) (*OAuthAccessTokenApplyConfiguration, error) {
	return extractOAuthAccessToken(oAuthAccessToken, fieldManager, "")
}

// ExtractOAuthAccessTokenStatus is the same as ExtractOAuthAccessToken except
// that it extracts the status subresource applied configuration.
// Experimental!
func ExtractOAuthAccessTokenStatus(oAuthAccessToken *oauthv1.OAuthAccessToken, fieldManager string) (*OAuthAccessTokenApplyConfiguration, error) {
	return extractOAuthAccessToken(oAuthAccessToken, fieldManager, "status")
}

func extractOAuthAccessToken(oAuthAccessToken *oauthv1.OAuthAccessToken, fieldManager string, subresource string) (*OAuthAccessTokenApplyConfiguration, error) {
	b := &OAuthAccessTokenApplyConfiguration{}
	err := managedfields.ExtractInto(oAuthAccessToken, internal.Parser().Type("com.github.openshift.api.oauth.v1.OAuthAccessToken"), fieldManager, b, subresource)
	if err != nil {
		return nil, err
	}
	b.WithName(oAuthAccessToken.Name)

	b.WithKind("OAuthAccessToken")
	b.WithAPIVersion("oauth.openshift.io/v1")
	return b, nil
}

// WithKind sets the Kind field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Kind field is set to the value of the last call.
func (b *OAuthAccessTokenApplyConfiguration) WithKind(value string) *OAuthAccessTokenApplyConfiguration {
	b.Kind = &value
	return b
}

// WithAPIVersion sets the APIVersion field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the APIVersion field is set to the value of the last call.
func (b *OAuthAccessTokenApplyConfiguration) WithAPIVersion(value string) *OAuthAccessTokenApplyConfiguration {
	b.APIVersion = &value
	return b
}

// WithName sets the Name field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Name field is set to the value of the last call.
func (b *OAuthAccessTokenApplyConfiguration) WithName(value string) *OAuthAccessTokenApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.Name = &value
	return b
}

// WithGenerateName sets the GenerateName field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the GenerateName field is set to the value of the last call.
func (b *OAuthAccessTokenApplyConfiguration) WithGenerateName(value string) *OAuthAccessTokenApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.GenerateName = &value
	return b
}

// WithNamespace sets the Namespace field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Namespace field is set to the value of the last call.
func (b *OAuthAccessTokenApplyConfiguration) WithNamespace(value string) *OAuthAccessTokenApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.Namespace = &value
	return b
}

// WithUID sets the UID field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the UID field is set to the value of the last call.
func (b *OAuthAccessTokenApplyConfiguration) WithUID(value types.UID) *OAuthAccessTokenApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.UID = &value
	return b
}

// WithResourceVersion sets the ResourceVersion field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ResourceVersion field is set to the value of the last call.
func (b *OAuthAccessTokenApplyConfiguration) WithResourceVersion(value string) *OAuthAccessTokenApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.ResourceVersion = &value
	return b
}

// WithGeneration sets the Generation field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Generation field is set to the value of the last call.
func (b *OAuthAccessTokenApplyConfiguration) WithGeneration(value int64) *OAuthAccessTokenApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.Generation = &value
	return b
}

// WithCreationTimestamp sets the CreationTimestamp field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the CreationTimestamp field is set to the value of the last call.
func (b *OAuthAccessTokenApplyConfiguration) WithCreationTimestamp(value metav1.Time) *OAuthAccessTokenApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.CreationTimestamp = &value
	return b
}

// WithDeletionTimestamp sets the DeletionTimestamp field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the DeletionTimestamp field is set to the value of the last call.
func (b *OAuthAccessTokenApplyConfiguration) WithDeletionTimestamp(value metav1.Time) *OAuthAccessTokenApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.DeletionTimestamp = &value
	return b
}

// WithDeletionGracePeriodSeconds sets the DeletionGracePeriodSeconds field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the DeletionGracePeriodSeconds field is set to the value of the last call.
func (b *OAuthAccessTokenApplyConfiguration) WithDeletionGracePeriodSeconds(value int64) *OAuthAccessTokenApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.DeletionGracePeriodSeconds = &value
	return b
}

// WithLabels puts the entries into the Labels field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, the entries provided by each call will be put on the Labels field,
// overwriting an existing map entries in Labels field with the same key.
func (b *OAuthAccessTokenApplyConfiguration) WithLabels(entries map[string]string) *OAuthAccessTokenApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	if b.Labels == nil && len(entries) > 0 {
		b.Labels = make(map[string]string, len(entries))
	}
	for k, v := range entries {
		b.Labels[k] = v
	}
	return b
}

// WithAnnotations puts the entries into the Annotations field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, the entries provided by each call will be put on the Annotations field,
// overwriting an existing map entries in Annotations field with the same key.
func (b *OAuthAccessTokenApplyConfiguration) WithAnnotations(entries map[string]string) *OAuthAccessTokenApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	if b.Annotations == nil && len(entries) > 0 {
		b.Annotations = make(map[string]string, len(entries))
	}
	for k, v := range entries {
		b.Annotations[k] = v
	}
	return b
}

// WithOwnerReferences adds the given value to the OwnerReferences field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the OwnerReferences field.
func (b *OAuthAccessTokenApplyConfiguration) WithOwnerReferences(values ...*v1.OwnerReferenceApplyConfiguration) *OAuthAccessTokenApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithOwnerReferences")
		}
		b.OwnerReferences = append(b.OwnerReferences, *values[i])
	}
	return b
}

// WithFinalizers adds the given value to the Finalizers field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Finalizers field.
func (b *OAuthAccessTokenApplyConfiguration) WithFinalizers(values ...string) *OAuthAccessTokenApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	for i := range values {
		b.Finalizers = append(b.Finalizers, values[i])
	}
	return b
}

func (b *OAuthAccessTokenApplyConfiguration) ensureObjectMetaApplyConfigurationExists() {
	if b.ObjectMetaApplyConfiguration == nil {
		b.ObjectMetaApplyConfiguration = &v1.ObjectMetaApplyConfiguration{}
	}
}

// WithClientName sets the ClientName field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ClientName field is set to the value of the last call.
func (b *OAuthAccessTokenApplyConfiguration) WithClientName(value string) *OAuthAccessTokenApplyConfiguration {
	b.ClientName = &value
	return b
}

// WithExpiresIn sets the ExpiresIn field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ExpiresIn field is set to the value of the last call.
func (b *OAuthAccessTokenApplyConfiguration) WithExpiresIn(value int64) *OAuthAccessTokenApplyConfiguration {
	b.ExpiresIn = &value
	return b
}

// WithScopes adds the given value to the Scopes field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Scopes field.
func (b *OAuthAccessTokenApplyConfiguration) WithScopes(values ...string) *OAuthAccessTokenApplyConfiguration {
	for i := range values {
		b.Scopes = append(b.Scopes, values[i])
	}
	return b
}

// WithRedirectURI sets the RedirectURI field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the RedirectURI field is set to the value of the last call.
func (b *OAuthAccessTokenApplyConfiguration) WithRedirectURI(value string) *OAuthAccessTokenApplyConfiguration {
	b.RedirectURI = &value
	return b
}

// WithUserName sets the UserName field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the UserName field is set to the value of the last call.
func (b *OAuthAccessTokenApplyConfiguration) WithUserName(value string) *OAuthAccessTokenApplyConfiguration {
	b.UserName = &value
	return b
}

// WithUserUID sets the UserUID field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the UserUID field is set to the value of the last call.
func (b *OAuthAccessTokenApplyConfiguration) WithUserUID(value string) *OAuthAccessTokenApplyConfiguration {
	b.UserUID = &value
	return b
}

// WithAuthorizeToken sets the AuthorizeToken field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the AuthorizeToken field is set to the value of the last call.
func (b *OAuthAccessTokenApplyConfiguration) WithAuthorizeToken(value string) *OAuthAccessTokenApplyConfiguration {
	b.AuthorizeToken = &value
	return b
}

// WithRefreshToken sets the RefreshToken field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the RefreshToken field is set to the value of the last call.
func (b *OAuthAccessTokenApplyConfiguration) WithRefreshToken(value string) *OAuthAccessTokenApplyConfiguration {
	b.RefreshToken = &value
	return b
}

// WithInactivityTimeoutSeconds sets the InactivityTimeoutSeconds field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the InactivityTimeoutSeconds field is set to the value of the last call.
func (b *OAuthAccessTokenApplyConfiguration) WithInactivityTimeoutSeconds(value int32) *OAuthAccessTokenApplyConfiguration {
	b.InactivityTimeoutSeconds = &value
	return b
}
