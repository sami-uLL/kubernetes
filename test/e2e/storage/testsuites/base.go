/*
Copyright 2018 The Kubernetes Authors.

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

package testsuites

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/test/e2e/framework"
	"k8s.io/kubernetes/test/e2e/storage/drivers"
	"k8s.io/kubernetes/test/e2e/storage/testpatterns"
	"k8s.io/kubernetes/test/e2e/storage/types"
)

// GenericVolumeTestResource is a generic implementation of TestResource that wil be able to
// be used in most of types.TestSuites.
// See volume_io.go or volumes.go in test/e2e/storage/testsuites/ for how to use this resource.
// Also, see subpath.go in the same directory for how to extend and use it.
type GenericVolumeTestResource struct {
	driver    types.TestDriver
	volType   string
	volSource *v1.VolumeSource
	pvc       *v1.PersistentVolumeClaim
	pv        *v1.PersistentVolume
	sc        *storagev1.StorageClass

	DriverTestResources types.DriverTestResources
}

func getTestNameStr(suite types.TestSuite, pattern testpatterns.TestPattern) string {
	tsInfo := suite.GetTestSuiteInfo()
	return fmt.Sprintf("[Testpattern: %s]%s %s%s", pattern.Name, pattern.FeatureTag, tsInfo.Name, tsInfo.FeatureTag)
}

// RunTestSuite runs all testpatterns of all testSuites for a driver
func RunTestSuite(f *framework.Framework, config framework.VolumeTestConfig, driver types.TestDriver, tsInits []func() types.TestSuite) {
	for _, testSuiteInit := range tsInits {
		suite := testSuiteInit()
		tsInfo := suite.GetTestSuiteInfo()

		for _, pattern := range tsInfo.TestPatterns {
			suite.ExecTest(driver, pattern)
		}
	}
}

// SkipUnsupportedTest will skip tests if the combination of driver, testsuite, and testpattern
// is not suitable to be tested.
// Whether it needs to be skipped is checked by following steps:
// 1. Check if Whether volType is supported by driver from its interface
// 2. Check if fsType is supported by driver
// 3. Check with driver specific logic
// 4. Check with testSuite specific logic
func SkipUnsupportedTest(suite types.TestSuite, driver types.TestDriver, pattern testpatterns.TestPattern) {
	dInfo := driver.GetDriverInfo()

	// 1. Check if Whether volType is supported by driver from its interface
	var isSupported bool
	switch pattern.VolType {
	case testpatterns.InlineVolume:
		_, isSupported = driver.(types.InlineVolumeTestDriver)
	case testpatterns.PreprovisionedPV:
		_, isSupported = driver.(types.PreprovisionedPVTestDriver)
	case testpatterns.DynamicPV:
		_, isSupported = driver.(types.DynamicPVTestDriver)
	default:
		isSupported = false
	}

	if !isSupported {
		framework.Skipf("Driver %s doesn't support %v -- skipping", dInfo.Name, pattern.VolType)
	}

	// 2. Check if fsType is supported by driver
	if !dInfo.SupportedFsType.Has(pattern.FsType) {
		framework.Skipf("Driver %s doesn't support %v -- skipping", dInfo.Name, pattern.FsType)
	}

	// 3. Check with driver specific logic
	driver.SkipUnsupportedTest(pattern)

	// 4. Check with testSuite specific logic
	suite.SkipUnsupportedTest(pattern, driver)
}

var _ types.TestResource = &GenericVolumeTestResource{}

// SetupResource sets up GenericVolumeTestResource
func (r *GenericVolumeTestResource) SetupResource(driver types.TestDriver, pattern testpatterns.TestPattern) {
	r.driver = driver
	dInfo := driver.GetDriverInfo()
	f := dInfo.Framework
	cs := f.ClientSet
	fsType := pattern.FsType
	volType := pattern.VolType

	// Create volume for pre-provisioned volume tests
	r.DriverTestResources = drivers.CreateVolume(driver, volType)

	switch volType {
	case testpatterns.InlineVolume:
		framework.Logf("Creating resource for inline volume")
		if iDriver, ok := driver.(types.InlineVolumeTestDriver); ok {
			r.volSource = iDriver.GetVolumeSource(false, fsType, r.DriverTestResources)
			r.volType = dInfo.Name
		}
	case testpatterns.PreprovisionedPV:
		framework.Logf("Creating resource for pre-provisioned PV")
		if pDriver, ok := driver.(types.PreprovisionedPVTestDriver); ok {
			pvSource := pDriver.GetPersistentVolumeSource(false, fsType, r.DriverTestResources)
			if pvSource != nil {
				r.volSource, r.pv, r.pvc = createVolumeSourceWithPVCPV(f, dInfo.Name, pvSource, false)
			}
			r.volType = fmt.Sprintf("%s-preprovisionedPV", dInfo.Name)
		}
	case testpatterns.DynamicPV:
		framework.Logf("Creating resource for dynamic PV")
		if dDriver, ok := driver.(types.DynamicPVTestDriver); ok {
			claimSize := "2Gi"
			r.sc = dDriver.GetDynamicProvisionStorageClass(fsType)

			By("creating a StorageClass " + r.sc.Name)
			var err error
			r.sc, err = cs.StorageV1().StorageClasses().Create(r.sc)
			Expect(err).NotTo(HaveOccurred())

			if r.sc != nil {
				r.volSource, r.pv, r.pvc = createVolumeSourceWithPVCPVFromDynamicProvisionSC(
					f, dInfo.Name, claimSize, r.sc, false, nil)
			}
			r.volType = fmt.Sprintf("%s-dynamicPV", dInfo.Name)
		}
	default:
		framework.Failf("GenericVolumeTestResource doesn't support: %s", volType)
	}

	if r.volSource == nil {
		framework.Skipf("Driver %s doesn't support %v -- skipping", dInfo.Name, volType)
	}
}

// CleanupResource clean up GenericVolumeTestResource
func (r *GenericVolumeTestResource) CleanupResource(driver types.TestDriver, pattern testpatterns.TestPattern) {
	dInfo := driver.GetDriverInfo()
	f := dInfo.Framework
	volType := pattern.VolType

	if r.pvc != nil || r.pv != nil {
		By("Deleting pv and pvc")
		if errs := framework.PVPVCCleanup(f.ClientSet, f.Namespace.Name, r.pv, r.pvc); len(errs) != 0 {
			framework.Failf("Failed to delete PVC or PV: %v", utilerrors.NewAggregate(errs))
		}
	}

	if r.sc != nil {
		By("Deleting sc")
		deleteStorageClass(f.ClientSet, r.sc.Name)
	}

	// Cleanup volume for pre-provisioned volume tests
	drivers.DeleteVolume(driver, volType, r.DriverTestResources)
}

func createVolumeSourceWithPVCPV(
	f *framework.Framework,
	name string,
	pvSource *v1.PersistentVolumeSource,
	readOnly bool,
) (*v1.VolumeSource, *v1.PersistentVolume, *v1.PersistentVolumeClaim) {
	pvConfig := framework.PersistentVolumeConfig{
		NamePrefix:       fmt.Sprintf("%s-", name),
		StorageClassName: f.Namespace.Name,
		PVSource:         *pvSource,
	}
	pvcConfig := framework.PersistentVolumeClaimConfig{
		StorageClassName: &f.Namespace.Name,
	}

	framework.Logf("Creating PVC and PV")
	pv, pvc, err := framework.CreatePVCPV(f.ClientSet, pvConfig, pvcConfig, f.Namespace.Name, false)
	Expect(err).NotTo(HaveOccurred(), "PVC, PV creation failed")

	err = framework.WaitOnPVandPVC(f.ClientSet, f.Namespace.Name, pv, pvc)
	Expect(err).NotTo(HaveOccurred(), "PVC, PV failed to bind")

	volSource := &v1.VolumeSource{
		PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
			ClaimName: pvc.Name,
			ReadOnly:  readOnly,
		},
	}
	return volSource, pv, pvc
}

func createVolumeSourceWithPVCPVFromDynamicProvisionSC(
	f *framework.Framework,
	name string,
	claimSize string,
	sc *storagev1.StorageClass,
	readOnly bool,
	volMode *v1.PersistentVolumeMode,
) (*v1.VolumeSource, *v1.PersistentVolume, *v1.PersistentVolumeClaim) {
	cs := f.ClientSet
	ns := f.Namespace.Name

	By("creating a claim")
	pvc := getClaim(claimSize, ns)
	pvc.Spec.StorageClassName = &sc.Name
	if volMode != nil {
		pvc.Spec.VolumeMode = volMode
	}

	var err error
	pvc, err = cs.CoreV1().PersistentVolumeClaims(ns).Create(pvc)
	Expect(err).NotTo(HaveOccurred())

	err = framework.WaitForPersistentVolumeClaimPhase(v1.ClaimBound, cs, pvc.Namespace, pvc.Name, framework.Poll, framework.ClaimProvisionTimeout)
	Expect(err).NotTo(HaveOccurred())

	pvc, err = cs.CoreV1().PersistentVolumeClaims(pvc.Namespace).Get(pvc.Name, metav1.GetOptions{})
	Expect(err).NotTo(HaveOccurred())

	pv, err := cs.CoreV1().PersistentVolumes().Get(pvc.Spec.VolumeName, metav1.GetOptions{})
	Expect(err).NotTo(HaveOccurred())

	volSource := &v1.VolumeSource{
		PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
			ClaimName: pvc.Name,
			ReadOnly:  readOnly,
		},
	}
	return volSource, pv, pvc
}

func getClaim(claimSize string, ns string) *v1.PersistentVolumeClaim {
	claim := v1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "pvc-",
			Namespace:    ns,
		},
		Spec: v1.PersistentVolumeClaimSpec{
			AccessModes: []v1.PersistentVolumeAccessMode{
				v1.ReadWriteOnce,
			},
			Resources: v1.ResourceRequirements{
				Requests: v1.ResourceList{
					v1.ResourceName(v1.ResourceStorage): resource.MustParse(claimSize),
				},
			},
		},
	}

	return &claim
}

// deleteStorageClass deletes the passed in StorageClass and catches errors other than "Not Found"
func deleteStorageClass(cs clientset.Interface, className string) {
	err := cs.StorageV1().StorageClasses().Delete(className, nil)
	if err != nil && !apierrs.IsNotFound(err) {
		Expect(err).NotTo(HaveOccurred())
	}
}
