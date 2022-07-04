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

// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	"context"
	json "encoding/json"
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	corev1 "k8s.io/client-go/applyconfigurations/core/v1"
	scheme "k8s.io/client-go/kubernetes/scheme"
	rest "k8s.io/client-go/rest"
)

// ResourceClassesGetter has a method to return a ResourceClassInterface.
// A group's client should implement this interface.
type ResourceClassesGetter interface {
	ResourceClasses() ResourceClassInterface
}

// ResourceClassInterface has methods to work with ResourceClass resources.
type ResourceClassInterface interface {
	Create(ctx context.Context, resourceClass *v1.ResourceClass, opts metav1.CreateOptions) (*v1.ResourceClass, error)
	Update(ctx context.Context, resourceClass *v1.ResourceClass, opts metav1.UpdateOptions) (*v1.ResourceClass, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.ResourceClass, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.ResourceClassList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.ResourceClass, err error)
	Apply(ctx context.Context, resourceClass *corev1.ResourceClassApplyConfiguration, opts metav1.ApplyOptions) (result *v1.ResourceClass, err error)
	ResourceClassExpansion
}

// resourceClasses implements ResourceClassInterface
type resourceClasses struct {
	client rest.Interface
}

// newResourceClasses returns a ResourceClasses
func newResourceClasses(c *CoreV1Client) *resourceClasses {
	return &resourceClasses{
		client: c.RESTClient(),
	}
}

// Get takes name of the resourceClass, and returns the corresponding resourceClass object, and an error if there is any.
func (c *resourceClasses) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.ResourceClass, err error) {
	result = &v1.ResourceClass{}
	err = c.client.Get().
		Resource("resourceclasses").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ResourceClasses that match those selectors.
func (c *resourceClasses) List(ctx context.Context, opts metav1.ListOptions) (result *v1.ResourceClassList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.ResourceClassList{}
	err = c.client.Get().
		Resource("resourceclasses").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested resourceClasses.
func (c *resourceClasses) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("resourceclasses").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a resourceClass and creates it.  Returns the server's representation of the resourceClass, and an error, if there is any.
func (c *resourceClasses) Create(ctx context.Context, resourceClass *v1.ResourceClass, opts metav1.CreateOptions) (result *v1.ResourceClass, err error) {
	result = &v1.ResourceClass{}
	err = c.client.Post().
		Resource("resourceclasses").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(resourceClass).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a resourceClass and updates it. Returns the server's representation of the resourceClass, and an error, if there is any.
func (c *resourceClasses) Update(ctx context.Context, resourceClass *v1.ResourceClass, opts metav1.UpdateOptions) (result *v1.ResourceClass, err error) {
	result = &v1.ResourceClass{}
	err = c.client.Put().
		Resource("resourceclasses").
		Name(resourceClass.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(resourceClass).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the resourceClass and deletes it. Returns an error if one occurs.
func (c *resourceClasses) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Resource("resourceclasses").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *resourceClasses) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("resourceclasses").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched resourceClass.
func (c *resourceClasses) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.ResourceClass, err error) {
	result = &v1.ResourceClass{}
	err = c.client.Patch(pt).
		Resource("resourceclasses").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// Apply takes the given apply declarative configuration, applies it and returns the applied resourceClass.
func (c *resourceClasses) Apply(ctx context.Context, resourceClass *corev1.ResourceClassApplyConfiguration, opts metav1.ApplyOptions) (result *v1.ResourceClass, err error) {
	if resourceClass == nil {
		return nil, fmt.Errorf("resourceClass provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(resourceClass)
	if err != nil {
		return nil, err
	}
	name := resourceClass.Name
	if name == nil {
		return nil, fmt.Errorf("resourceClass.Name must be provided to Apply")
	}
	result = &v1.ResourceClass{}
	err = c.client.Patch(types.ApplyPatchType).
		Resource("resourceclasses").
		Name(*name).
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
