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

package apps_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"k8s.io/kubernetes/pkg/kubectl/apps"
)

var _ = Describe("When GroupKindVisitor accepts a GroupKind", func() {

	var visitor *TestGroupKindVisitor

	BeforeEach(func() {
		visitor = &TestGroupKindVisitor{map[string]int{}}
	})

	It("should Visit DaemonSet iff the Kind is a DaemonSet", func() {
		kind := apps.GroupKindElement{
			Kind:  "DaemonSet",
			Group: "apps",
		}
		kind.Accept(visitor)
		Expect(visitor.visits).To(Equal(map[string]int{
			"DaemonSet": 1,
		}))

		kind = apps.GroupKindElement{
			Kind:  "DaemonSet",
			Group: "extensions",
		}
		kind.Accept(visitor)
		Expect(visitor.visits).To(Equal(map[string]int{
			"DaemonSet": 2,
		}))
	})

	It("should Visit Deployment iff the Kind is a Deployment", func() {
		kind := apps.GroupKindElement{
			Kind:  "Deployment",
			Group: "apps",
		}
		kind.Accept(visitor)
		Expect(visitor.visits).To(Equal(map[string]int{
			"Deployment": 1,
		}))

		kind = apps.GroupKindElement{
			Kind:  "Deployment",
			Group: "extensions",
		}
		kind.Accept(visitor)
		Expect(visitor.visits).To(Equal(map[string]int{
			"Deployment": 2,
		}))
	})

	It("should Visit Job iff the Kind is a Job", func() {
		kind := apps.GroupKindElement{
			Kind:  "Job",
			Group: "apps",
		}
		kind.Accept(visitor)
		Expect(visitor.visits).To(Equal(map[string]int{
			"Job": 1,
		}))

		kind = apps.GroupKindElement{
			Kind:  "Job",
			Group: "extensions",
		}
		kind.Accept(visitor)
		Expect(visitor.visits).To(Equal(map[string]int{
			"Job": 2,
		}))
	})

	It("should Visit Pod iff the Kind is a Pod", func() {
		kind := apps.GroupKindElement{
			Kind:  "Pod",
			Group: "",
		}
		kind.Accept(visitor)
		Expect(visitor.visits).To(Equal(map[string]int{
			"Pod": 1,
		}))

		kind = apps.GroupKindElement{
			Kind:  "Pod",
			Group: "core",
		}
		kind.Accept(visitor)
		Expect(visitor.visits).To(Equal(map[string]int{
			"Pod": 2,
		}))
	})

	It("should Visit ReplicationController iff the Kind is a ReplicationController", func() {
		kind := apps.GroupKindElement{
			Kind:  "ReplicationController",
			Group: "",
		}
		kind.Accept(visitor)
		Expect(visitor.visits).To(Equal(map[string]int{
			"ReplicationController": 1,
		}))

		kind = apps.GroupKindElement{
			Kind:  "ReplicationController",
			Group: "core",
		}
		kind.Accept(visitor)
		Expect(visitor.visits).To(Equal(map[string]int{
			"ReplicationController": 2,
		}))
	})

	It("should Visit ReplicaSet iff the Kind is a ReplicaSet", func() {
		kind := apps.GroupKindElement{
			Kind:  "ReplicaSet",
			Group: "apps",
		}
		kind.Accept(visitor)
		Expect(visitor.visits).To(Equal(map[string]int{
			"ReplicaSet": 1,
		}))

		kind = apps.GroupKindElement{
			Kind:  "ReplicaSet",
			Group: "extensions",
		}
		kind.Accept(visitor)
		Expect(visitor.visits).To(Equal(map[string]int{
			"ReplicaSet": 2,
		}))
	})

	It("should Visit StatefulSet iff the Kind is a StatefulSet", func() {
		kind := apps.GroupKindElement{
			Kind:  "StatefulSet",
			Group: "apps",
		}
		kind.Accept(visitor)
		Expect(visitor.visits).To(Equal(map[string]int{
			"StatefulSet": 1,
		}))
	})
})

var _ apps.GroupKindVisitor = &TestGroupKindVisitor{}

type TestGroupKindVisitor struct {
	visits map[string]int
}

func (t *TestGroupKindVisitor) Visit(kind apps.GroupKindElement) { t.visits[kind.Kind] += 1 }

func (t *TestGroupKindVisitor) VisitDaemonSet(kind apps.GroupKindElement)             { t.Visit(kind) }
func (t *TestGroupKindVisitor) VisitDeployment(kind apps.GroupKindElement)            { t.Visit(kind) }
func (t *TestGroupKindVisitor) VisitJob(kind apps.GroupKindElement)                   { t.Visit(kind) }
func (t *TestGroupKindVisitor) VisitPod(kind apps.GroupKindElement)                   { t.Visit(kind) }
func (t *TestGroupKindVisitor) VisitReplicaSet(kind apps.GroupKindElement)            { t.Visit(kind) }
func (t *TestGroupKindVisitor) VisitReplicationController(kind apps.GroupKindElement) { t.Visit(kind) }
func (t *TestGroupKindVisitor) VisitStatefulSet(kind apps.GroupKindElement)           { t.Visit(kind) }
