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

package upgrades

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/test/e2e/framework"

	. "github.com/onsi/ginkgo"
)

// SecretUpgradeTest test that a secret is available before and after
// a cluster upgrade.
type SecretUpgradeTest struct {
	secret *v1.Secret
}

func (t *SecretUpgradeTest) Setup(f *framework.Framework) {
	secretName := "upgrade-secret"

	// Grab a unique namespace so we don't collide.
	ns, err := f.CreateNamespace("secret-upgrade", nil)
	framework.ExpectNoError(err)

	t.secret = &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns.Name,
			Name:      secretName,
		},
		Data: map[string][]byte{
			"data": []byte("keep it secret"),
		},
	}

	By("Creating a secret")
	if t.secret, err = f.ClientSet.Core().Secrets(ns.Name).Create(t.secret); err != nil {
		framework.Failf("unable to create test secret %s: %v", t.secret.Name, err)
	}

	By("Making sure the secret is consumable")
	t.testPod(f)
}

func (t *SecretUpgradeTest) Test(f *framework.Framework, done <-chan struct{}, upgrade UpgradeType) {
	<-done
	By("Consuming the secret after upgrade")
	t.testPod(f)
}

// Teardown cleans up any remaining resources.
func (t *SecretUpgradeTest) Teardown(f *framework.Framework) {
	// rely on the namespace deletion to clean up everything
}

func (t *SecretUpgradeTest) testPod(f *framework.Framework) {
	volumeName := "secret-volume"
	volumeMountPath := "/etc/secret-volume"

	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pod-secrets-" + string(uuid.NewUUID()),
			Namespace: t.secret.ObjectMeta.Namespace,
		},
		Spec: v1.PodSpec{
			Volumes: []v1.Volume{
				{
					Name: volumeName,
					VolumeSource: v1.VolumeSource{
						Secret: &v1.SecretVolumeSource{
							SecretName: t.secret.ObjectMeta.Name,
						},
					},
				},
			},
			Containers: []v1.Container{
				{
					Name:  "secret-volume-test",
					Image: "gcr.io/google_containers/mounttest:0.7",
					Args: []string{
						fmt.Sprintf("--file_content=%s/data", volumeMountPath),
						fmt.Sprintf("--file_mode=%s/data", volumeMountPath),
					},
					VolumeMounts: []v1.VolumeMount{
						{
							Name:      volumeName,
							MountPath: volumeMountPath,
						},
					},
				},
			},
			RestartPolicy: v1.RestartPolicyNever,
		},
	}

	expectedOutput := []string{
		"content of file \"/etc/secret-volume/data\": keep it secret",
		"mode of file \"/etc/secret-volume/data\": -rw-r--r--",
	}

	f.TestContainerOutput("consume secrets", pod, 0, expectedOutput)
}
