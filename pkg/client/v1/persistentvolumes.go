/*
Copyright 2014 The Kubernetes Authors All rights reserved.

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

package v1

import (
	"fmt"

	"k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/fields"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/watch"
)

type PersistentVolumesInterface interface {
	PersistentVolumes() PersistentVolumeInterface
}

// PersistentVolumeInterface has methods to work with PersistentVolume resources.
type PersistentVolumeInterface interface {
	List(label labels.Selector, field fields.Selector) (*v1.PersistentVolumeList, error)
	Get(name string) (*v1.PersistentVolume, error)
	Create(volume *v1.PersistentVolume) (*v1.PersistentVolume, error)
	Update(volume *v1.PersistentVolume) (*v1.PersistentVolume, error)
	UpdateStatus(persistentVolume *v1.PersistentVolume) (*v1.PersistentVolume, error)
	Delete(name string) error
	Watch(label labels.Selector, field fields.Selector, resourceVersion string) (watch.Interface, error)
}

// persistentVolumes implements PersistentVolumesInterface
type persistentVolumes struct {
	client *Client
}

func newPersistentVolumes(c *Client) *persistentVolumes {
	return &persistentVolumes{c}
}

func (c *persistentVolumes) List(label labels.Selector, field fields.Selector) (result *v1.PersistentVolumeList, err error) {
	result = &v1.PersistentVolumeList{}
	err = c.client.Get().
		Resource("persistentVolumes").
		LabelsSelectorParam(label).
		FieldsSelectorParam(field).
		Do().
		Into(result)

	return result, err
}

func (c *persistentVolumes) Get(name string) (result *v1.PersistentVolume, err error) {
	result = &v1.PersistentVolume{}
	err = c.client.Get().Resource("persistentVolumes").Name(name).Do().Into(result)
	return
}

func (c *persistentVolumes) Create(volume *v1.PersistentVolume) (result *v1.PersistentVolume, err error) {
	result = &v1.PersistentVolume{}
	err = c.client.Post().Resource("persistentVolumes").Body(volume).Do().Into(result)
	return
}

func (c *persistentVolumes) Update(volume *v1.PersistentVolume) (result *v1.PersistentVolume, err error) {
	result = &v1.PersistentVolume{}
	if len(volume.ResourceVersion) == 0 {
		err = fmt.Errorf("invalid update object, missing resource version: %v", volume)
		return
	}
	err = c.client.Put().Resource("persistentVolumes").Name(volume.Name).Body(volume).Do().Into(result)
	return
}

func (c *persistentVolumes) UpdateStatus(volume *v1.PersistentVolume) (result *v1.PersistentVolume, err error) {
	result = &v1.PersistentVolume{}
	err = c.client.Put().Resource("persistentVolumes").Name(volume.Name).SubResource("status").Body(volume).Do().Into(result)
	return
}

func (c *persistentVolumes) Delete(name string) error {
	return c.client.Delete().Resource("persistentVolumes").Name(name).Do().Error()
}

func (c *persistentVolumes) Watch(label labels.Selector, field fields.Selector, resourceVersion string) (watch.Interface, error) {
	return c.client.Get().
		Prefix("watch").
		Resource("persistentVolumes").
		Param("resourceVersion", resourceVersion).
		LabelsSelectorParam(label).
		FieldsSelectorParam(field).
		Watch()
}
