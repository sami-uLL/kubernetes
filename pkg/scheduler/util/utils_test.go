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

package util

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	clientsetfake "k8s.io/client-go/kubernetes/fake"
	clienttesting "k8s.io/client-go/testing"
	extenderv1 "k8s.io/kube-scheduler/extender/v1"
)

func TestGetPodFullName(t *testing.T) {
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test",
			Name:      "pod",
		},
	}
	got := GetPodFullName(pod)
	expected := fmt.Sprintf("%s_%s", pod.Name, pod.Namespace)
	if got != expected {
		t.Errorf("Got wrong full name, got: %s, expected: %s", got, expected)
	}
}

func newPriorityPodWithStartTime(name string, priority int32, startTime time.Time) *v1.Pod {
	return &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: v1.PodSpec{
			Priority: &priority,
		},
		Status: v1.PodStatus{
			StartTime: &metav1.Time{Time: startTime},
		},
	}
}

func TestGetEarliestPodStartTime(t *testing.T) {
	currentTime := time.Now()
	pod1 := newPriorityPodWithStartTime("pod1", 1, currentTime.Add(time.Second))
	pod2 := newPriorityPodWithStartTime("pod2", 2, currentTime.Add(time.Second))
	pod3 := newPriorityPodWithStartTime("pod3", 2, currentTime)
	victims := &extenderv1.Victims{
		Pods: []*v1.Pod{pod1, pod2, pod3},
	}
	startTime := GetEarliestPodStartTime(victims)
	if !startTime.Equal(pod3.Status.StartTime) {
		t.Errorf("Got wrong earliest pod start time")
	}

	pod1 = newPriorityPodWithStartTime("pod1", 2, currentTime)
	pod2 = newPriorityPodWithStartTime("pod2", 2, currentTime.Add(time.Second))
	pod3 = newPriorityPodWithStartTime("pod3", 2, currentTime.Add(2*time.Second))
	victims = &extenderv1.Victims{
		Pods: []*v1.Pod{pod1, pod2, pod3},
	}
	startTime = GetEarliestPodStartTime(victims)
	if !startTime.Equal(pod1.Status.StartTime) {
		t.Errorf("Got wrong earliest pod start time, got %v, expected %v", startTime, pod1.Status.StartTime)
	}
}

func TestMoreImportantPod(t *testing.T) {
	currentTime := time.Now()
	pod1 := newPriorityPodWithStartTime("pod1", 1, currentTime)
	pod2 := newPriorityPodWithStartTime("pod2", 2, currentTime.Add(time.Second))
	pod3 := newPriorityPodWithStartTime("pod3", 2, currentTime)

	tests := map[string]struct {
		p1       *v1.Pod
		p2       *v1.Pod
		expected bool
	}{
		"Pod with higher priority": {
			p1:       pod1,
			p2:       pod2,
			expected: false,
		},
		"Pod with older created time": {
			p1:       pod2,
			p2:       pod3,
			expected: false,
		},
		"Pods with same start time": {
			p1:       pod3,
			p2:       pod1,
			expected: true,
		},
	}

	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			got := MoreImportantPod(v.p1, v.p2)
			if got != v.expected {
				t.Errorf("expected %t but got %t", v.expected, got)
			}
		})
	}
}

func TestClearNominatedNodeName(t *testing.T) {
	t.Parallel()
	statusErr := &apierrors.StatusError{
		ErrStatus: metav1.Status{Reason: metav1.StatusReasonUnknown},
	}

	tests := []struct {
		name                  string
		pods                  []*v1.Pod
		patchError            error
		expectedPatchError    utilerrors.Aggregate
		expectedPatchRequests int
		expectedPatchData     []string
	}{
		{
			name: "Should make patch requests to clear node name",
			pods: []*v1.Pod{
				{
					ObjectMeta: metav1.ObjectMeta{Name: "foo1"},
					Status:     v1.PodStatus{NominatedNodeName: "node1"},
				},
				{
					ObjectMeta: metav1.ObjectMeta{Name: "foo2"},
					Status:     v1.PodStatus{NominatedNodeName: "node2"},
				},
			},
			expectedPatchRequests: 2,
			expectedPatchData:     []string{`{"status":{"nominatedNodeName":null}}`, `{"status":{"nominatedNodeName":null}}`},
		},
		{
			name: "Should make patch request only for pods that have NominatedNodeName",
			pods: []*v1.Pod{
				{
					ObjectMeta: metav1.ObjectMeta{Name: "foo1"},
					Status:     v1.PodStatus{NominatedNodeName: ""},
				},
				{
					ObjectMeta: metav1.ObjectMeta{Name: "foo2"},
					Status:     v1.PodStatus{NominatedNodeName: "node1"},
				},
			},
			expectedPatchRequests: 1,
			expectedPatchData:     []string{`{"status":{"nominatedNodeName":null}}`},
		},
		{
			name: "Should make patch requests for all pods even if a patch for one pods fails",
			pods: []*v1.Pod{
				{
					// In this test, the pod named "err" returns an error in the patch.
					ObjectMeta: metav1.ObjectMeta{Name: "err"},
					Status:     v1.PodStatus{NominatedNodeName: "node1"},
				},
				{
					ObjectMeta: metav1.ObjectMeta{Name: "foo1"},
					Status:     v1.PodStatus{NominatedNodeName: "node1"},
				},
			},
			patchError:            statusErr,
			expectedPatchError:    utilerrors.NewAggregate([]error{statusErr}),
			expectedPatchRequests: 2,
			expectedPatchData:     []string{`{"status":{"nominatedNodeName":null}}`},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualPatchRequests := 0
			var actualPatchData []string
			cs := &clientsetfake.Clientset{}
			cs.AddReactor("patch", "pods", func(action clienttesting.Action) (bool, runtime.Object, error) {
				actualPatchRequests++
				patch := action.(clienttesting.PatchAction)
				// If the pod name is "err", return an error.
				if patch.GetName() == "err" {
					return true, nil, test.patchError
				}
				actualPatchData = append(actualPatchData, string(patch.GetPatch()))
				// For this test, we don't care about the result of the patched pod, just that we got the expected
				// patch request, so just returning &v1.Pod{} here is OK because scheduler doesn't use the response.
				return true, &v1.Pod{}, nil

			})

			if err := ClearNominatedNodeName(cs, test.pods...); err != nil && err.Is(test.expectedPatchError) {
				t.Fatalf("ClearNominatedNodeName's error mismatch: Actual was %v, but expected %v", err, test.expectedPatchError)
			}

			if actualPatchRequests != test.expectedPatchRequests {
				t.Fatalf("Actual patch requests (%d) dos not equal expected patch requests (%d)", actualPatchRequests, test.expectedPatchRequests)
			}

			if test.expectedPatchRequests > 0 && !reflect.DeepEqual(actualPatchData, test.expectedPatchData) {
				t.Fatalf("Patch data mismatch: Actual was %v, but expected %v", actualPatchData, test.expectedPatchData)
			}
		})
	}
}

func TestPatchPodStatus(t *testing.T) {
	tests := []struct {
		name           string
		pod            v1.Pod
		statusToUpdate v1.PodStatus
	}{
		{
			name: "Should update pod conditions successfully",
			pod: v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "ns",
					Name:      "pod1",
				},
				Spec: v1.PodSpec{
					ImagePullSecrets: []v1.LocalObjectReference{{Name: "foo"}},
				},
			},
			statusToUpdate: v1.PodStatus{
				Conditions: []v1.PodCondition{
					{
						Type:   v1.PodScheduled,
						Status: v1.ConditionFalse,
					},
				},
			},
		},
		{
			// ref: #101697, #94626 - ImagePullSecrets are allowed to have empty secret names
			// which would fail the 2-way merge patch generation on Pod patches
			// due to the mergeKey being the name field
			name: "Should update pod conditions successfully on a pod Spec with secrets with empty name",
			pod: v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "ns",
					Name:      "pod2",
				},
				Spec: v1.PodSpec{
					// this will serialize to imagePullSecrets:[{}]
					ImagePullSecrets: make([]v1.LocalObjectReference, 1),
				},
			},
			statusToUpdate: v1.PodStatus{
				Conditions: []v1.PodCondition{
					{
						Type:   v1.PodScheduled,
						Status: v1.ConditionFalse,
					},
				},
			},
		},
	}
	client := clientsetfake.NewSimpleClientset()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := client.CoreV1().Pods(tc.pod.Namespace).Create(context.TODO(), &tc.pod, metav1.CreateOptions{})
			if err != nil {
				t.Fatal(err)
			}

			err = PatchPodStatus(client, &tc.pod, &tc.statusToUpdate)
			if err != nil {
				t.Fatal(err)
			}

			retrievedPod, err := client.CoreV1().Pods(tc.pod.Namespace).Get(context.TODO(), tc.pod.Name, metav1.GetOptions{})
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tc.statusToUpdate, retrievedPod.Status); diff != "" {
				t.Errorf("unexpected pod status (-want,+got):\n%s", diff)
			}
		})
	}
}

func TestIsScalarResourceName(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		resourceName v1.ResourceName
		want         bool
	}{
		{
			name:         "Should be true with extended resources name",
			resourceName: "example.com/foo",
			want:         true,
		},
		{
			name:         "Should be false with not extended resources name",
			resourceName: "requests.foo",
			want:         false,
		},
		{
			name:         "Should be true with hugePage resource name",
			resourceName: "hugepages-2Mi",
			want:         true,
		},
		{
			name:         "Should be false with not HugePage resource name",
			resourceName: "cpu",
			want:         false,
		},
		{
			name:         "Should be true with native resource name",
			resourceName: "kubernetes.io/resource-foo",
			want:         true,
		},
		{
			name:         "Should be true with attachable volume resource name",
			resourceName: "attachable-volumes-foo",
			want:         true,
		},
	}
	for _, test := range tests {
		if IsScalarResourceName(test.resourceName) != test.want {
			t.Errorf("%s: expected: %t, got: %t", test.name, test.want, !test.want)
		}
	}
}
