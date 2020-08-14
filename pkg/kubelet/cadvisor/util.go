/*
Copyright 2015 The Kubernetes Authors.

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

package cadvisor

import (
	goruntime "runtime"

	cadvisormetrics "github.com/google/cadvisor/container"
	cadvisorapi "github.com/google/cadvisor/info/v1"
	cadvisorapi2 "github.com/google/cadvisor/info/v2"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/klog/v2"
	v1helper "k8s.io/kubernetes/pkg/apis/core/v1/helper"
	kubetypes "k8s.io/kubernetes/pkg/kubelet/types"
)

const (
	// CrioSocket is the path to the CRI-O socket.
	// Please keep this in sync with the one in:
	// github.com/google/cadvisor/container/crio/client.go
	CrioSocket = "/var/run/crio/crio.sock"
)

var (
	// IncludedMetrics refers to default cadvisor metrics collected.
	IncludedMetrics = cadvisormetrics.MetricSet{
		cadvisormetrics.CpuUsageMetrics:     struct{}{},
		cadvisormetrics.MemoryUsageMetrics:  struct{}{},
		cadvisormetrics.CpuLoadMetrics:      struct{}{},
		cadvisormetrics.DiskIOMetrics:       struct{}{},
		cadvisormetrics.NetworkUsageMetrics: struct{}{},
		cadvisormetrics.AppMetrics:          struct{}{},
		cadvisormetrics.ProcessMetrics:      struct{}{},
	}
)

// CapacityFromMachineInfo returns the capacity of the resources from the machine info.
func CapacityFromMachineInfo(info *cadvisorapi.MachineInfo) v1.ResourceList {
	c := v1.ResourceList{
		v1.ResourceCPU: *resource.NewMilliQuantity(
			int64(info.NumCores*1000),
			resource.DecimalSI),
		v1.ResourceMemory: *resource.NewQuantity(
			int64(info.MemoryCapacity),
			resource.BinarySI),
	}

	// if huge pages are enabled, we report them as a schedulable resource on the node
	for _, hugepagesInfo := range info.HugePages {
		pageSizeBytes := int64(hugepagesInfo.PageSize * 1024)
		hugePagesBytes := pageSizeBytes * int64(hugepagesInfo.NumPages)
		pageSizeQuantity := resource.NewQuantity(pageSizeBytes, resource.BinarySI)
		c[v1helper.HugePageResourceName(*pageSizeQuantity)] = *resource.NewQuantity(hugePagesBytes, resource.BinarySI)
	}

	return c
}

// EphemeralStorageCapacityFromFsInfo returns the capacity of the ephemeral storage from the FsInfo.
func EphemeralStorageCapacityFromFsInfo(info cadvisorapi2.FsInfo) v1.ResourceList {
	c := v1.ResourceList{
		v1.ResourceEphemeralStorage: *resource.NewQuantity(
			int64(info.Capacity),
			resource.BinarySI),
	}
	return c
}

// UsingLegacyCadvisorStats returns true if container stats are provided by cadvisor instead of through the CRI.
// CRI integrations should get container metrics via CRI. Docker
// uses the built-in cadvisor to gather such metrics on Linux for
// historical reasons.
// TODO: cri-o relies on cadvisor as a temporary workaround. The code should
// be removed. Related issue:
// https://github.com/kubernetes/kubernetes/issues/51798
func UsingLegacyCadvisorStats(runtime, runtimeEndpoint string) bool {
	return (runtime == kubetypes.DockerContainerRuntime && goruntime.GOOS == "linux") ||
		runtimeEndpoint == CrioSocket || runtimeEndpoint == "unix://"+CrioSocket
}

// DefaultCadvisorMetricSet returns default metrics set for cadivsor interface of kubelet work background.
func DefaultCadvisorMetricSet(usingLegacyStats, usingLocalStorageCapacityIsolation, disableAcceleratorUsageMetrics bool) cadvisormetrics.MetricSet {
	metricSet := cadvisormetrics.MetricSet{}
	for metricKind := range IncludedMetrics {
		metricSet.Add(metricKind)
	}

	if usingLegacyStats || usingLocalStorageCapacityIsolation {
		metricSet.Add(cadvisormetrics.DiskUsageMetrics)
	}

	if !disableAcceleratorUsageMetrics {
		metricSet.Add(cadvisormetrics.AcceleratorUsageMetrics)
	}

	return metricSet
}

// CustomCadvisorMetricSet returns custom metrics set for cadvisor interface of kubelet work background.
// If feature gate CustomCadvisorMetrics equals to true,
// cadivsor interface can work with metrics specified.
func CustomCadvisorMetricSet(usingLegacyStats, usingLocalStorageCapacityIsolation, disableAcceleratorUsageMetrics bool, cadvisorMetricList []string) cadvisormetrics.MetricSet {
	metricSet := cadvisormetrics.MetricSet{}
	for _, metric := range cadvisorMetricList {
		if !cadvisormetrics.AllMetrics.Has(cadvisormetrics.MetricKind(metric)) {
			klog.Warningf("illegal cadvisor metric kind provided: %q", metric)
			continue
		}
		metricSet.Add(cadvisormetrics.MetricKind(metric))
	}

	// mapping usingLegacyStats and local storage feature gate, if false, remove custom disk metric from set.
	if !usingLegacyStats && !usingLocalStorageCapacityIsolation {
		delete(metricSet, cadvisormetrics.DiskUsageMetrics)
	}

	// if disable accelerator usage metric feature gate enabled, just remove it from set.
	if disableAcceleratorUsageMetrics {
		delete(metricSet, cadvisormetrics.AcceleratorUsageMetrics)
	}

	return metricSet
}
