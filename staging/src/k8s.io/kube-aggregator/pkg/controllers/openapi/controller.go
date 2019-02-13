/*
Copyright 2016 The Kubernetes Authors.

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

package openapi

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	restful "github.com/emicklei/go-restful"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apiserver/pkg/server"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"
	"k8s.io/kube-aggregator/pkg/apis/apiregistration"
	"k8s.io/kube-aggregator/pkg/controllers/openapi/aggregator"
	"k8s.io/kube-aggregator/pkg/controllers/openapi/download"
	"k8s.io/kube-openapi/pkg/builder"
	"k8s.io/kube-openapi/pkg/common"
	"k8s.io/kube-openapi/pkg/handler"
)

const (
	successfulUpdateDelay   = time.Minute
	failedUpdateMaxExpDelay = time.Hour

	localDelegateChainNamePattern = "k8s_internal_local_delegation_chain_%010d"
)

type syncAction int

const (
	syncRequeue syncAction = iota
	syncRequeueRateLimited
	syncNothing
)

// BuildAndRegisterAggregator registered OpenAPI aggregator handler. This function is not thread safe as it only being called on startup.
func BuildAndRegisterAggregator(downloader *download.Downloader, delegationTarget server.DelegationTarget, webServices []*restful.WebService,
	config *common.Config, pathHandler common.PathHandler) (aggregator.SpecAggregator, error) {
	s := aggregator.NewSpecAggregator()

	i := 0
	// Build Aggregator's spec
	aggregatorOpenAPISpec, err := builder.BuildOpenAPISpec(webServices, config)
	if err != nil {
		return nil, err
	}

	// Reserving non-name spec for aggregator's Spec.
	s.addLocalSpec(aggregatorOpenAPISpec, nil, fmt.Sprintf(localDelegateChainNamePattern, i), "")
	i++
	for delegate := delegationTarget; delegate != nil; delegate = delegate.NextDelegate() {
		handler := delegate.UnprotectedHandler()
		if handler == nil {
			continue
		}
		delegateSpec, etag, _, err := downloader.Download(handler, "")
		if err != nil {
			return nil, err
		}
		if delegateSpec == nil {
			continue
		}
		s.addLocalSpec(delegateSpec, handler, fmt.Sprintf(localDelegateChainNamePattern, i), etag)
		i++
	}

	// Build initial spec to serve.
	specToServe, err := s.buildOpenAPISpec()
	if err != nil {
		return nil, err
	}

	// Install handler
	s.openAPIVersionedService, err = handler.RegisterOpenAPIVersionedService(nil, "/openapi/v2", pathHandler)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// AggregationController periodically check for changes in OpenAPI specs of APIServices and update/remove
// them if necessary.
type AggregationController struct {
	openAPIAggregationManager aggregator.SpecAggregator
	queue                     workqueue.RateLimitingInterface
	downloader                *download.Downloader

	lock     sync.Mutex
	handlers map[string]http.Handler

	// To allow injection for testing.
	syncHandler func(key string) (syncAction, error)
}

// NewAggregationController creates new OpenAPI aggregation controller.
func NewAggregationController(downloader *aggregator.Downloader, openAPIAggregationManager aggregator.SpecAggregator) *AggregationController {
	c := &AggregationController{
		openAPIAggregationManager: openAPIAggregationManager,
		queue: workqueue.NewNamedRateLimitingQueue(
			workqueue.NewItemExponentialFailureRateLimiter(successfulUpdateDelay, failedUpdateMaxExpDelay), "APIServiceOpenAPIAggregationControllerQueue1"),
		downloader: downloader,
		handlers:   map[string]http.Handler{},
	}

	c.syncHandler = c.sync

	return c
}

// Run starts OpenAPI AggregationController
func (c *AggregationController) Run(stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	klog.Infof("Starting OpenAPI AggregationController")
	defer klog.Infof("Shutting down OpenAPI AggregationController")

	go wait.Until(c.runWorker, time.Second, stopCh)

	<-stopCh
}

func (c *AggregationController) runWorker() {
	for c.processNextWorkItem() {
	}
}

// processNextWorkItem deals with one key off the queue.  It returns false when it's time to quit.
func (c *AggregationController) processNextWorkItem() bool {
	key, quit := c.queue.Get()
	defer c.queue.Done(key)
	if quit {
		return false
	}

	klog.Infof("OpenAPI AggregationController: Processing item %s", key)

	action, err := c.syncHandler(key.(string))
	if err == nil {
		c.queue.Forget(key)
	} else {
		utilruntime.HandleError(fmt.Errorf("loading OpenAPI spec for %q failed with: %v", key, err))
	}

	switch action {
	case syncRequeue:
		klog.Infof("OpenAPI AggregationController: action for item %s: Requeue.", key)
		c.queue.AddAfter(key, successfulUpdateDelay)
	case syncRequeueRateLimited:
		klog.Infof("OpenAPI AggregationController: action for item %s: Rate Limited Requeue.", key)
		c.queue.AddRateLimited(key)
	case syncNothing:
		klog.Infof("OpenAPI AggregationController: action for item %s: Nothing (removed from the queue).", key)
	}

	return true
}

func (c *AggregationController) sync(key string) (syncAction, error) {
	oldSpec, etag, exists := c.openAPIAggregationManager.Spec(key)
	if !exists {
		return syncNothing, nil
	}

	c.lock.Lock()
	defer c.lock.Unlock()
	handler := c.handlers[key]
	if handler == nil {
		return syncNothing, nil
	}

	returnSpec, newEtag, httpStatus, err := c.downloader.Download(handler, etag)
	switch {
	case err != nil:
		return syncRequeueRateLimited, err
	case httpStatus == http.StatusNotModified:
	case httpStatus == http.StatusNotFound || returnSpec == nil:
		return syncRequeueRateLimited, fmt.Errorf("OpenAPI spec does not exist")
	case httpStatus == http.StatusOK:
		if err := c.openAPIAggregationManager.UpdateSpec(key, returnSpec, newEtag); err != nil {
			return syncRequeueRateLimited, err
		}
	}
	return syncRequeue, nil
}

// AddAPIService adds a new API Service to OpenAPI Aggregation.
func (c *AggregationController) AddAPIService(handler http.Handler, apiService *apiregistration.APIService) {
	if apiService.Spec.Service == nil {
		return
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	// TODO: combine APIServices from the same aggregated apiserver by choosing the same key
	key := apiService.Name
	if err := c.openAPIAggregationManager.AddUpdateService(key, nil, apiService); err != nil {
		utilruntime.HandleError(fmt.Errorf("adding %q to AggregationController failed with: %v", apiService.Name, err))
	}
	c.handlers[key] = handler

	c.queue.AddAfter(apiService.Name, time.Second)
}

// UpdateAPIService updates API Service's info and handler.
func (c *AggregationController) UpdateAPIService(handler http.Handler, apiService *apiregistration.APIService) {
	if apiService.Spec.Service == nil {
		return
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	// TODO: combine APIServices from the same aggregated apiserver by choosing the same name
	key := apiService.Name
	if err := c.openAPIAggregationManager.AddUpdateSpec(key, apiService); err != nil {
		utilruntime.HandleError(fmt.Errorf("updating %q to AggregationController failed with: %v", apiService.Name, err))
	}
	c.handlers[key] = handler

	if c.queue.NumRequeues(key) > 0 {
		// The item has failed before. Remove it from failure queue and
		// update it in a second
		c.queue.Forget(key)
		c.queue.AddAfter(key, time.Second)
	}
	// Else: The item has been succeeded before and it will be updated soon (after successfulUpdateDelay)
	// we don't add it again as it will cause a duplication of items.
}

// RemoveAPIService removes API Service from OpenAPI Aggregation Controller.
func (c *AggregationController) RemoveAPIService(apiServiceName string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	// split APIService name which always has the version.group pattern
	ns := strings.SplitN(apiServiceName, ".", 2)
	version, group := ns[0], ns[1]

	key = apiServiceName
	if err := c.openAPIAggregationManager.RemoveService(key, schema.GroupVersion{group, version}); err != nil {
		utilruntime.HandleError(fmt.Errorf("removing %q from AggregationController failed with: %v", apiServiceName, err))
	}
	delete(c.handlers, apiServiceName)

	// This will only remove it if it was failing before. If it was successful, processNextWorkItem will figure it out
	// and will not add it again to the queue.
	c.queue.Forget(apiServiceName)
}
