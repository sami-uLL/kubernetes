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


//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by kcp code-generator. DO NOT EDIT.

package v1

import (
	kcpclient "github.com/kcp-dev/apimachinery/pkg/client"
	"github.com/kcp-dev/logicalcluster/v2"
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"


	apiextensionsv1client "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1"
)

// CustomResourceDefinitionsClusterGetter has a method to return a CustomResourceDefinitionClusterInterface.
// A group's cluster client should implement this interface.
type CustomResourceDefinitionsClusterGetter interface {
	CustomResourceDefinitions() CustomResourceDefinitionClusterInterface
}

// CustomResourceDefinitionClusterInterface can operate on CustomResourceDefinitions across all clusters,
// or scope down to one cluster and return a apiextensionsv1client.CustomResourceDefinitionInterface.
type CustomResourceDefinitionClusterInterface interface {
	Cluster(logicalcluster.Name) apiextensionsv1client.CustomResourceDefinitionInterface
	List(ctx context.Context, opts metav1.ListOptions) (*apiextensionsv1.CustomResourceDefinitionList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
}

type customResourceDefinitionsClusterInterface struct {
	clientCache kcpclient.Cache[*apiextensionsv1client.ApiextensionsV1Client]
}

// Cluster scopes the client down to a particular cluster.
func (c *customResourceDefinitionsClusterInterface) Cluster(name logicalcluster.Name) apiextensionsv1client.CustomResourceDefinitionInterface {
	if name == logicalcluster.Wildcard {
		panic("A specific cluster must be provided when scoping, not the wildcard.")
	}

	return c.clientCache.ClusterOrDie(name).CustomResourceDefinitions()
}


// List returns the entire collection of all CustomResourceDefinitions across all clusters. 
func (c *customResourceDefinitionsClusterInterface) List(ctx context.Context, opts metav1.ListOptions) (*apiextensionsv1.CustomResourceDefinitionList, error) {
	return c.clientCache.ClusterOrDie(logicalcluster.Wildcard).CustomResourceDefinitions().List(ctx, opts)
}

// Watch begins to watch all CustomResourceDefinitions across all clusters.
func (c *customResourceDefinitionsClusterInterface) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.clientCache.ClusterOrDie(logicalcluster.Wildcard).CustomResourceDefinitions().Watch(ctx, opts)
}
