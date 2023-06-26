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

package windows

import (
	"github.com/onsi/ginkgo/v2"
	e2eskipper "k8s.io/kubernetes/test/e2e/framework/skipper"
)

// SIGDescribe annotates the test with the SIG label.
func SIGDescribe(text string, args ...interface{}) bool {
	funcs := []interface{}{
		ginkgo.BeforeEach(func() {
			// all tests in this package are Windows specific
			e2eskipper.SkipUnlessNodeOSDistroIs("windows")
		}),
	}
	args = append(funcs, args...)
	return ginkgo.Describe("[sig-windows] "+text, args...)
}
