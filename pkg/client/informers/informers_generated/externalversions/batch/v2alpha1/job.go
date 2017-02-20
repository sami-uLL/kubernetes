/*
Copyright 2017 The Kubernetes Authors.

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

// This file was automatically generated by informer-gen

package v2alpha1

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	batch_v2alpha1 "k8s.io/apis/pkg/apis/batch/v2alpha1"
	cache "k8s.io/client-go/tools/cache"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/clientset"
	internalinterfaces "k8s.io/kubernetes/pkg/client/informers/informers_generated/externalversions/internalinterfaces"
	v2alpha1 "k8s.io/kubernetes/pkg/client/listers/batch/v2alpha1"
	time "time"
)

// JobInformer provides access to a shared informer and lister for
// Jobs.
type JobInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v2alpha1.JobLister
}

type jobInformer struct {
	factory internalinterfaces.SharedInformerFactory
}

func newJobInformer(client clientset.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	sharedIndexInformer := cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				return client.BatchV2alpha1().Jobs(v1.NamespaceAll).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				return client.BatchV2alpha1().Jobs(v1.NamespaceAll).Watch(options)
			},
		},
		&batch_v2alpha1.Job{},
		resyncPeriod,
		cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
	)

	return sharedIndexInformer
}

func (f *jobInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&batch_v2alpha1.Job{}, newJobInformer)
}

func (f *jobInformer) Lister() v2alpha1.JobLister {
	return v2alpha1.NewJobLister(f.Informer().GetIndexer())
}
