// +build !providerless

/*
Copyright 2020 The Kubernetes Authors.

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

package network

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	compute "google.golang.org/api/compute/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/wait"
	clientset "k8s.io/client-go/kubernetes"
	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/kubernetes/test/e2e/framework"
	e2edeployment "k8s.io/kubernetes/test/e2e/framework/deployment"
	e2ekubesystem "k8s.io/kubernetes/test/e2e/framework/kubesystem"
	e2enetwork "k8s.io/kubernetes/test/e2e/framework/network"
	e2enode "k8s.io/kubernetes/test/e2e/framework/node"
	"k8s.io/kubernetes/test/e2e/framework/providers/gce"
	e2erc "k8s.io/kubernetes/test/e2e/framework/rc"
	e2eservice "k8s.io/kubernetes/test/e2e/framework/service"
	e2eskipper "k8s.io/kubernetes/test/e2e/framework/skipper"
	gcecloud "k8s.io/legacy-cloud-providers/gce"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = SIGDescribe("Services with Cloud LoadBalancers", func() {

	f := framework.NewDefaultFramework("services")

	var cs clientset.Interface
	serviceLBNames := []string{}

	ginkgo.BeforeEach(func() {
		cs = f.ClientSet
	})

	ginkgo.AfterEach(func() {
		if ginkgo.CurrentGinkgoTestDescription().Failed {
			DescribeSvc(f.Namespace.Name)
		}
		for _, lb := range serviceLBNames {
			framework.Logf("cleaning load balancer resource for %s", lb)
			e2eservice.CleanupServiceResources(cs, lb, framework.TestContext.CloudConfig.Region, framework.TestContext.CloudConfig.Zone)
		}
		//reset serviceLBNames
		serviceLBNames = []string{}
	})

	// TODO: Get rid of [DisabledForLargeClusters] tag when issue #56138 is fixed
	ginkgo.It("should be able to change the type and ports of a service [Slow] [DisabledForLargeClusters]", func() {
		// requires cloud load-balancer support
		e2eskipper.SkipUnlessProviderIs("gce", "gke", "aws")

		loadBalancerSupportsUDP := !framework.ProviderIs("aws")

		loadBalancerLagTimeout := e2eservice.LoadBalancerLagTimeoutDefault
		if framework.ProviderIs("aws") {
			loadBalancerLagTimeout = e2eservice.LoadBalancerLagTimeoutAWS
		}
		loadBalancerCreateTimeout := e2eservice.GetServiceLoadBalancerCreationTimeout(cs)

		// This test is more monolithic than we'd like because LB turnup can be
		// very slow, so we lumped all the tests into one LB lifecycle.

		serviceName := "mutability-test"
		ns1 := f.Namespace.Name // LB1 in ns1 on TCP
		framework.Logf("namespace for TCP test: %s", ns1)

		ginkgo.By("creating a second namespace")
		namespacePtr, err := f.CreateNamespace("services", nil)
		framework.ExpectNoError(err, "failed to create namespace")
		ns2 := namespacePtr.Name // LB2 in ns2 on UDP
		framework.Logf("namespace for UDP test: %s", ns2)

		nodeIP, err := e2enode.PickIP(cs) // for later
		framework.ExpectNoError(err)

		// Test TCP and UDP Services.  Services with the same name in different
		// namespaces should get different node ports and load balancers.

		ginkgo.By("creating a TCP service " + serviceName + " with type=ClusterIP in namespace " + ns1)
		tcpJig := e2eservice.NewTestJig(cs, ns1, serviceName)
		tcpService, err := tcpJig.CreateTCPService(nil)
		framework.ExpectNoError(err)

		ginkgo.By("creating a UDP service " + serviceName + " with type=ClusterIP in namespace " + ns2)
		udpJig := e2eservice.NewTestJig(cs, ns2, serviceName)
		udpService, err := udpJig.CreateUDPService(nil)
		framework.ExpectNoError(err)

		ginkgo.By("verifying that TCP and UDP use the same port")
		if tcpService.Spec.Ports[0].Port != udpService.Spec.Ports[0].Port {
			framework.Failf("expected to use the same port for TCP and UDP")
		}
		svcPort := int(tcpService.Spec.Ports[0].Port)
		framework.Logf("service port (TCP and UDP): %d", svcPort)

		ginkgo.By("creating a pod to be part of the TCP service " + serviceName)
		_, err = tcpJig.Run(nil)
		framework.ExpectNoError(err)

		ginkgo.By("creating a pod to be part of the UDP service " + serviceName)
		_, err = udpJig.Run(nil)
		framework.ExpectNoError(err)

		// Change the services to NodePort.

		ginkgo.By("changing the TCP service to type=NodePort")
		tcpService, err = tcpJig.UpdateService(func(s *v1.Service) {
			s.Spec.Type = v1.ServiceTypeNodePort
		})
		framework.ExpectNoError(err)
		tcpNodePort := int(tcpService.Spec.Ports[0].NodePort)
		framework.Logf("TCP node port: %d", tcpNodePort)

		ginkgo.By("changing the UDP service to type=NodePort")
		udpService, err = udpJig.UpdateService(func(s *v1.Service) {
			s.Spec.Type = v1.ServiceTypeNodePort
		})
		framework.ExpectNoError(err)
		udpNodePort := int(udpService.Spec.Ports[0].NodePort)
		framework.Logf("UDP node port: %d", udpNodePort)

		ginkgo.By("hitting the TCP service's NodePort")
		e2eservice.TestReachableHTTP(nodeIP, tcpNodePort, e2eservice.KubeProxyLagTimeout)

		ginkgo.By("hitting the UDP service's NodePort")
		testReachableUDP(nodeIP, udpNodePort, e2eservice.KubeProxyLagTimeout)

		// Change the services to LoadBalancer.

		// Here we test that LoadBalancers can receive static IP addresses.  This isn't
		// necessary, but is an additional feature this monolithic test checks.
		requestedIP := ""
		staticIPName := ""
		if framework.ProviderIs("gce", "gke") {
			ginkgo.By("creating a static load balancer IP")
			staticIPName = fmt.Sprintf("e2e-external-lb-test-%s", framework.RunID)
			gceCloud, err := gce.GetGCECloud()
			framework.ExpectNoError(err, "failed to get GCE cloud provider")

			err = gceCloud.ReserveRegionAddress(&compute.Address{Name: staticIPName}, gceCloud.Region())
			defer func() {
				if staticIPName != "" {
					// Release GCE static IP - this is not kube-managed and will not be automatically released.
					if err := gceCloud.DeleteRegionAddress(staticIPName, gceCloud.Region()); err != nil {
						framework.Logf("failed to release static IP %s: %v", staticIPName, err)
					}
				}
			}()
			framework.ExpectNoError(err, "failed to create region address: %s", staticIPName)
			reservedAddr, err := gceCloud.GetRegionAddress(staticIPName, gceCloud.Region())
			framework.ExpectNoError(err, "failed to get region address: %s", staticIPName)

			requestedIP = reservedAddr.Address
			framework.Logf("Allocated static load balancer IP: %s", requestedIP)
		}

		ginkgo.By("changing the TCP service to type=LoadBalancer")
		tcpService, err = tcpJig.UpdateService(func(s *v1.Service) {
			s.Spec.LoadBalancerIP = requestedIP // will be "" if not applicable
			s.Spec.Type = v1.ServiceTypeLoadBalancer
		})
		framework.ExpectNoError(err)

		if loadBalancerSupportsUDP {
			ginkgo.By("changing the UDP service to type=LoadBalancer")
			udpService, err = udpJig.UpdateService(func(s *v1.Service) {
				s.Spec.Type = v1.ServiceTypeLoadBalancer
			})
			framework.ExpectNoError(err)
		}
		serviceLBNames = append(serviceLBNames, cloudprovider.DefaultLoadBalancerName(tcpService))
		if loadBalancerSupportsUDP {
			serviceLBNames = append(serviceLBNames, cloudprovider.DefaultLoadBalancerName(udpService))
		}

		ginkgo.By("waiting for the TCP service to have a load balancer")
		// Wait for the load balancer to be created asynchronously
		tcpService, err = tcpJig.WaitForLoadBalancer(loadBalancerCreateTimeout)
		framework.ExpectNoError(err)
		if int(tcpService.Spec.Ports[0].NodePort) != tcpNodePort {
			framework.Failf("TCP Spec.Ports[0].NodePort changed (%d -> %d) when not expected", tcpNodePort, tcpService.Spec.Ports[0].NodePort)
		}
		if requestedIP != "" && e2eservice.GetIngressPoint(&tcpService.Status.LoadBalancer.Ingress[0]) != requestedIP {
			framework.Failf("unexpected TCP Status.LoadBalancer.Ingress (expected %s, got %s)", requestedIP, e2eservice.GetIngressPoint(&tcpService.Status.LoadBalancer.Ingress[0]))
		}
		tcpIngressIP := e2eservice.GetIngressPoint(&tcpService.Status.LoadBalancer.Ingress[0])
		framework.Logf("TCP load balancer: %s", tcpIngressIP)

		if framework.ProviderIs("gce", "gke") {
			// Do this as early as possible, which overrides the `defer` above.
			// This is mostly out of fear of leaking the IP in a timeout case
			// (as of this writing we're not 100% sure where the leaks are
			// coming from, so this is first-aid rather than surgery).
			ginkgo.By("demoting the static IP to ephemeral")
			if staticIPName != "" {
				gceCloud, err := gce.GetGCECloud()
				framework.ExpectNoError(err, "failed to get GCE cloud provider")
				// Deleting it after it is attached "demotes" it to an
				// ephemeral IP, which can be auto-released.
				if err := gceCloud.DeleteRegionAddress(staticIPName, gceCloud.Region()); err != nil {
					framework.Failf("failed to release static IP %s: %v", staticIPName, err)
				}
				staticIPName = ""
			}
		}

		var udpIngressIP string
		if loadBalancerSupportsUDP {
			ginkgo.By("waiting for the UDP service to have a load balancer")
			// 2nd one should be faster since they ran in parallel.
			udpService, err = udpJig.WaitForLoadBalancer(loadBalancerCreateTimeout)
			framework.ExpectNoError(err)
			if int(udpService.Spec.Ports[0].NodePort) != udpNodePort {
				framework.Failf("UDP Spec.Ports[0].NodePort changed (%d -> %d) when not expected", udpNodePort, udpService.Spec.Ports[0].NodePort)
			}
			udpIngressIP = e2eservice.GetIngressPoint(&udpService.Status.LoadBalancer.Ingress[0])
			framework.Logf("UDP load balancer: %s", udpIngressIP)

			ginkgo.By("verifying that TCP and UDP use different load balancers")
			if tcpIngressIP == udpIngressIP {
				framework.Failf("Load balancers are not different: %s", e2eservice.GetIngressPoint(&tcpService.Status.LoadBalancer.Ingress[0]))
			}
		}

		ginkgo.By("hitting the TCP service's NodePort")
		e2eservice.TestReachableHTTP(nodeIP, tcpNodePort, e2eservice.KubeProxyLagTimeout)

		ginkgo.By("hitting the UDP service's NodePort")
		testReachableUDP(nodeIP, udpNodePort, e2eservice.KubeProxyLagTimeout)

		ginkgo.By("hitting the TCP service's LoadBalancer")
		e2eservice.TestReachableHTTP(tcpIngressIP, svcPort, loadBalancerLagTimeout)

		if loadBalancerSupportsUDP {
			ginkgo.By("hitting the UDP service's LoadBalancer")
			testReachableUDP(udpIngressIP, svcPort, loadBalancerLagTimeout)
		}

		// Change the services' node ports.

		ginkgo.By("changing the TCP service's NodePort")
		tcpService, err = tcpJig.ChangeServiceNodePort(tcpNodePort)
		framework.ExpectNoError(err)
		tcpNodePortOld := tcpNodePort
		tcpNodePort = int(tcpService.Spec.Ports[0].NodePort)
		if tcpNodePort == tcpNodePortOld {
			framework.Failf("TCP Spec.Ports[0].NodePort (%d) did not change", tcpNodePort)
		}
		if e2eservice.GetIngressPoint(&tcpService.Status.LoadBalancer.Ingress[0]) != tcpIngressIP {
			framework.Failf("TCP Status.LoadBalancer.Ingress changed (%s -> %s) when not expected", tcpIngressIP, e2eservice.GetIngressPoint(&tcpService.Status.LoadBalancer.Ingress[0]))
		}
		framework.Logf("TCP node port: %d", tcpNodePort)

		ginkgo.By("changing the UDP service's NodePort")
		udpService, err = udpJig.ChangeServiceNodePort(udpNodePort)
		framework.ExpectNoError(err)
		udpNodePortOld := udpNodePort
		udpNodePort = int(udpService.Spec.Ports[0].NodePort)
		if udpNodePort == udpNodePortOld {
			framework.Failf("UDP Spec.Ports[0].NodePort (%d) did not change", udpNodePort)
		}
		if loadBalancerSupportsUDP && e2eservice.GetIngressPoint(&udpService.Status.LoadBalancer.Ingress[0]) != udpIngressIP {
			framework.Failf("UDP Status.LoadBalancer.Ingress changed (%s -> %s) when not expected", udpIngressIP, e2eservice.GetIngressPoint(&udpService.Status.LoadBalancer.Ingress[0]))
		}
		framework.Logf("UDP node port: %d", udpNodePort)

		ginkgo.By("hitting the TCP service's new NodePort")
		e2eservice.TestReachableHTTP(nodeIP, tcpNodePort, e2eservice.KubeProxyLagTimeout)

		ginkgo.By("hitting the UDP service's new NodePort")
		testReachableUDP(nodeIP, udpNodePort, e2eservice.KubeProxyLagTimeout)

		ginkgo.By("checking the old TCP NodePort is closed")
		testNotReachableHTTP(nodeIP, tcpNodePortOld, e2eservice.KubeProxyLagTimeout)

		ginkgo.By("checking the old UDP NodePort is closed")
		testNotReachableUDP(nodeIP, udpNodePortOld, e2eservice.KubeProxyLagTimeout)

		ginkgo.By("hitting the TCP service's LoadBalancer")
		e2eservice.TestReachableHTTP(tcpIngressIP, svcPort, loadBalancerLagTimeout)

		if loadBalancerSupportsUDP {
			ginkgo.By("hitting the UDP service's LoadBalancer")
			testReachableUDP(udpIngressIP, svcPort, loadBalancerLagTimeout)
		}

		// Change the services' main ports.

		ginkgo.By("changing the TCP service's port")
		tcpService, err = tcpJig.UpdateService(func(s *v1.Service) {
			s.Spec.Ports[0].Port++
		})
		framework.ExpectNoError(err)
		svcPortOld := svcPort
		svcPort = int(tcpService.Spec.Ports[0].Port)
		if svcPort == svcPortOld {
			framework.Failf("TCP Spec.Ports[0].Port (%d) did not change", svcPort)
		}
		if int(tcpService.Spec.Ports[0].NodePort) != tcpNodePort {
			framework.Failf("TCP Spec.Ports[0].NodePort (%d) changed", tcpService.Spec.Ports[0].NodePort)
		}
		if e2eservice.GetIngressPoint(&tcpService.Status.LoadBalancer.Ingress[0]) != tcpIngressIP {
			framework.Failf("TCP Status.LoadBalancer.Ingress changed (%s -> %s) when not expected", tcpIngressIP, e2eservice.GetIngressPoint(&tcpService.Status.LoadBalancer.Ingress[0]))
		}

		ginkgo.By("changing the UDP service's port")
		udpService, err = udpJig.UpdateService(func(s *v1.Service) {
			s.Spec.Ports[0].Port++
		})
		framework.ExpectNoError(err)
		if int(udpService.Spec.Ports[0].Port) != svcPort {
			framework.Failf("UDP Spec.Ports[0].Port (%d) did not change", udpService.Spec.Ports[0].Port)
		}
		if int(udpService.Spec.Ports[0].NodePort) != udpNodePort {
			framework.Failf("UDP Spec.Ports[0].NodePort (%d) changed", udpService.Spec.Ports[0].NodePort)
		}
		if loadBalancerSupportsUDP && e2eservice.GetIngressPoint(&udpService.Status.LoadBalancer.Ingress[0]) != udpIngressIP {
			framework.Failf("UDP Status.LoadBalancer.Ingress changed (%s -> %s) when not expected", udpIngressIP, e2eservice.GetIngressPoint(&udpService.Status.LoadBalancer.Ingress[0]))
		}

		framework.Logf("service port (TCP and UDP): %d", svcPort)

		ginkgo.By("hitting the TCP service's NodePort")
		e2eservice.TestReachableHTTP(nodeIP, tcpNodePort, e2eservice.KubeProxyLagTimeout)

		ginkgo.By("hitting the UDP service's NodePort")
		testReachableUDP(nodeIP, udpNodePort, e2eservice.KubeProxyLagTimeout)

		ginkgo.By("hitting the TCP service's LoadBalancer")
		e2eservice.TestReachableHTTP(tcpIngressIP, svcPort, loadBalancerCreateTimeout)

		if loadBalancerSupportsUDP {
			ginkgo.By("hitting the UDP service's LoadBalancer")
			testReachableUDP(udpIngressIP, svcPort, loadBalancerCreateTimeout)
		}

		ginkgo.By("Scaling the pods to 0")
		err = tcpJig.Scale(0)
		framework.ExpectNoError(err)
		err = udpJig.Scale(0)
		framework.ExpectNoError(err)

		ginkgo.By("looking for ICMP REJECT on the TCP service's NodePort")
		testRejectedHTTP(nodeIP, tcpNodePort, e2eservice.KubeProxyLagTimeout)

		ginkgo.By("looking for ICMP REJECT on the UDP service's NodePort")
		testRejectedUDP(nodeIP, udpNodePort, e2eservice.KubeProxyLagTimeout)

		ginkgo.By("looking for ICMP REJECT on the TCP service's LoadBalancer")
		testRejectedHTTP(tcpIngressIP, svcPort, loadBalancerCreateTimeout)

		if loadBalancerSupportsUDP {
			ginkgo.By("looking for ICMP REJECT on the UDP service's LoadBalancer")
			testRejectedUDP(udpIngressIP, svcPort, loadBalancerCreateTimeout)
		}

		ginkgo.By("Scaling the pods to 1")
		err = tcpJig.Scale(1)
		framework.ExpectNoError(err)
		err = udpJig.Scale(1)
		framework.ExpectNoError(err)

		ginkgo.By("hitting the TCP service's NodePort")
		e2eservice.TestReachableHTTP(nodeIP, tcpNodePort, e2eservice.KubeProxyLagTimeout)

		ginkgo.By("hitting the UDP service's NodePort")
		testReachableUDP(nodeIP, udpNodePort, e2eservice.KubeProxyLagTimeout)

		ginkgo.By("hitting the TCP service's LoadBalancer")
		e2eservice.TestReachableHTTP(tcpIngressIP, svcPort, loadBalancerCreateTimeout)

		if loadBalancerSupportsUDP {
			ginkgo.By("hitting the UDP service's LoadBalancer")
			testReachableUDP(udpIngressIP, svcPort, loadBalancerCreateTimeout)
		}

		// Change the services back to ClusterIP.

		ginkgo.By("changing TCP service back to type=ClusterIP")
		_, err = tcpJig.UpdateService(func(s *v1.Service) {
			s.Spec.Type = v1.ServiceTypeClusterIP
			s.Spec.Ports[0].NodePort = 0
		})
		framework.ExpectNoError(err)
		// Wait for the load balancer to be destroyed asynchronously
		_, err = tcpJig.WaitForLoadBalancerDestroy(tcpIngressIP, svcPort, loadBalancerCreateTimeout)
		framework.ExpectNoError(err)

		ginkgo.By("changing UDP service back to type=ClusterIP")
		_, err = udpJig.UpdateService(func(s *v1.Service) {
			s.Spec.Type = v1.ServiceTypeClusterIP
			s.Spec.Ports[0].NodePort = 0
		})
		framework.ExpectNoError(err)
		if loadBalancerSupportsUDP {
			// Wait for the load balancer to be destroyed asynchronously
			_, err = udpJig.WaitForLoadBalancerDestroy(udpIngressIP, svcPort, loadBalancerCreateTimeout)
			framework.ExpectNoError(err)
		}

		ginkgo.By("checking the TCP NodePort is closed")
		testNotReachableHTTP(nodeIP, tcpNodePort, e2eservice.KubeProxyLagTimeout)

		ginkgo.By("checking the UDP NodePort is closed")
		testNotReachableUDP(nodeIP, udpNodePort, e2eservice.KubeProxyLagTimeout)

		ginkgo.By("checking the TCP LoadBalancer is closed")
		testNotReachableHTTP(tcpIngressIP, svcPort, loadBalancerLagTimeout)

		if loadBalancerSupportsUDP {
			ginkgo.By("checking the UDP LoadBalancer is closed")
			testNotReachableUDP(udpIngressIP, svcPort, loadBalancerLagTimeout)
		}
	})

	ginkgo.It("should be able to create an internal type load balancer [Slow]", func() {
		e2eskipper.SkipUnlessProviderIs("azure", "gke", "gce")

		createTimeout := e2eservice.GetServiceLoadBalancerCreationTimeout(cs)
		pollInterval := framework.Poll * 10

		namespace := f.Namespace.Name
		serviceName := "lb-internal"
		jig := e2eservice.NewTestJig(cs, namespace, serviceName)

		ginkgo.By("creating pod to be part of service " + serviceName)
		_, err := jig.Run(nil)
		framework.ExpectNoError(err)

		enableILB, disableILB := enableAndDisableInternalLB()

		isInternalEndpoint := func(lbIngress *v1.LoadBalancerIngress) bool {
			ingressEndpoint := e2eservice.GetIngressPoint(lbIngress)
			// Needs update for providers using hostname as endpoint.
			return strings.HasPrefix(ingressEndpoint, "10.")
		}

		ginkgo.By("creating a service with type LoadBalancer and cloud specific Internal-LB annotation enabled")
		svc, err := jig.CreateTCPService(func(svc *v1.Service) {
			svc.Spec.Type = v1.ServiceTypeLoadBalancer
			enableILB(svc)
		})
		framework.ExpectNoError(err)

		defer func() {
			ginkgo.By("Clean up loadbalancer service")
			e2eservice.WaitForServiceDeletedWithFinalizer(cs, svc.Namespace, svc.Name)
		}()

		svc, err = jig.WaitForLoadBalancer(createTimeout)
		framework.ExpectNoError(err)
		lbIngress := &svc.Status.LoadBalancer.Ingress[0]
		svcPort := int(svc.Spec.Ports[0].Port)
		// should have an internal IP.
		framework.ExpectEqual(isInternalEndpoint(lbIngress), true)

		// ILBs are not accessible from the test orchestrator, so it's necessary to use
		//  a pod to test the service.
		ginkgo.By("hitting the internal load balancer from pod")
		framework.Logf("creating pod with host network")
		hostExec := launchHostExecPod(f.ClientSet, f.Namespace.Name, "ilb-host-exec")

		framework.Logf("Waiting up to %v for service %q's internal LB to respond to requests", createTimeout, serviceName)
		tcpIngressIP := e2eservice.GetIngressPoint(lbIngress)
		if pollErr := wait.PollImmediate(pollInterval, createTimeout, func() (bool, error) {
			cmd := fmt.Sprintf(`curl -m 5 'http://%v:%v/echo?msg=hello'`, tcpIngressIP, svcPort)
			stdout, err := framework.RunHostCmd(hostExec.Namespace, hostExec.Name, cmd)
			if err != nil {
				framework.Logf("error curling; stdout: %v. err: %v", stdout, err)
				return false, nil
			}

			if !strings.Contains(stdout, "hello") {
				framework.Logf("Expected output to contain 'hello', got %q; retrying...", stdout)
				return false, nil
			}

			framework.Logf("Successful curl; stdout: %v", stdout)
			return true, nil
		}); pollErr != nil {
			framework.Failf("ginkgo.Failed to hit ILB IP, err: %v", pollErr)
		}

		ginkgo.By("switching to external type LoadBalancer")
		svc, err = jig.UpdateService(func(svc *v1.Service) {
			disableILB(svc)
		})
		framework.ExpectNoError(err)
		framework.Logf("Waiting up to %v for service %q to have an external LoadBalancer", createTimeout, serviceName)
		if pollErr := wait.PollImmediate(pollInterval, createTimeout, func() (bool, error) {
			svc, err := cs.CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
			if err != nil {
				return false, err
			}
			lbIngress = &svc.Status.LoadBalancer.Ingress[0]
			return !isInternalEndpoint(lbIngress), nil
		}); pollErr != nil {
			framework.Failf("Loadbalancer IP not changed to external.")
		}
		// should have an external IP.
		gomega.Expect(isInternalEndpoint(lbIngress)).To(gomega.BeFalse())

		ginkgo.By("hitting the external load balancer")
		framework.Logf("Waiting up to %v for service %q's external LB to respond to requests", createTimeout, serviceName)
		tcpIngressIP = e2eservice.GetIngressPoint(lbIngress)
		e2eservice.TestReachableHTTP(tcpIngressIP, svcPort, e2eservice.LoadBalancerLagTimeoutDefault)

		// GCE cannot test a specific IP because the test may not own it. This cloud specific condition
		// will be removed when GCP supports similar functionality.
		if framework.ProviderIs("azure") {
			ginkgo.By("switching back to interal type LoadBalancer, with static IP specified.")
			internalStaticIP := "10.240.11.11"
			svc, err = jig.UpdateService(func(svc *v1.Service) {
				svc.Spec.LoadBalancerIP = internalStaticIP
				enableILB(svc)
			})
			framework.ExpectNoError(err)
			framework.Logf("Waiting up to %v for service %q to have an internal LoadBalancer", createTimeout, serviceName)
			if pollErr := wait.PollImmediate(pollInterval, createTimeout, func() (bool, error) {
				svc, err := cs.CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
				if err != nil {
					return false, err
				}
				lbIngress = &svc.Status.LoadBalancer.Ingress[0]
				return isInternalEndpoint(lbIngress), nil
			}); pollErr != nil {
				framework.Failf("Loadbalancer IP not changed to internal.")
			}
			// should have the given static internal IP.
			framework.ExpectEqual(e2eservice.GetIngressPoint(lbIngress), internalStaticIP)
		}
	})

	// This test creates a load balancer, make sure its health check interval
	// equals to gceHcCheckIntervalSeconds. Then the interval is manipulated
	// to be something else, see if the interval will be reconciled.
	ginkgo.It("should reconcile LB health check interval [Slow][Serial]", func() {
		const gceHcCheckIntervalSeconds = int64(8)
		// This test is for clusters on GCE.
		// (It restarts kube-controller-manager, which we don't support on GKE)
		e2eskipper.SkipUnlessProviderIs("gce")
		e2eskipper.SkipUnlessSSHKeyPresent()

		clusterID, err := gce.GetClusterID(cs)
		if err != nil {
			framework.Failf("framework.GetClusterID(cs) = _, %v; want nil", err)
		}
		gceCloud, err := gce.GetGCECloud()
		if err != nil {
			framework.Failf("framework.GetGCECloud() = _, %v; want nil", err)
		}

		namespace := f.Namespace.Name
		serviceName := "lb-hc-int"
		jig := e2eservice.NewTestJig(cs, namespace, serviceName)

		ginkgo.By("create load balancer service")
		// Create loadbalancer service with source range from node[0] and podAccept
		svc, err := jig.CreateTCPService(func(svc *v1.Service) {
			svc.Spec.Type = v1.ServiceTypeLoadBalancer
		})
		framework.ExpectNoError(err)

		defer func() {
			ginkgo.By("Clean up loadbalancer service")
			e2eservice.WaitForServiceDeletedWithFinalizer(cs, svc.Namespace, svc.Name)
		}()

		svc, err = jig.WaitForLoadBalancer(e2eservice.GetServiceLoadBalancerCreationTimeout(cs))
		framework.ExpectNoError(err)

		hcName := gcecloud.MakeNodesHealthCheckName(clusterID)
		hc, err := gceCloud.GetHTTPHealthCheck(hcName)
		if err != nil {
			framework.Failf("gceCloud.GetHttpHealthCheck(%q) = _, %v; want nil", hcName, err)
		}
		framework.ExpectEqual(hc.CheckIntervalSec, gceHcCheckIntervalSeconds)

		ginkgo.By("modify the health check interval")
		hc.CheckIntervalSec = gceHcCheckIntervalSeconds - 1
		if err = gceCloud.UpdateHTTPHealthCheck(hc); err != nil {
			framework.Failf("gcecloud.UpdateHttpHealthCheck(%#v) = %v; want nil", hc, err)
		}

		ginkgo.By("restart kube-controller-manager")
		if err := e2ekubesystem.RestartControllerManager(); err != nil {
			framework.Failf("e2ekubesystem.RestartControllerManager() = %v; want nil", err)
		}
		if err := e2ekubesystem.WaitForControllerManagerUp(); err != nil {
			framework.Failf("e2ekubesystem.WaitForControllerManagerUp() = %v; want nil", err)
		}

		ginkgo.By("health check should be reconciled")
		pollInterval := framework.Poll * 10
		loadBalancerPropagationTimeout := e2eservice.GetServiceLoadBalancerPropagationTimeout(cs)
		if pollErr := wait.PollImmediate(pollInterval, loadBalancerPropagationTimeout, func() (bool, error) {
			hc, err := gceCloud.GetHTTPHealthCheck(hcName)
			if err != nil {
				framework.Logf("ginkgo.Failed to get HttpHealthCheck(%q): %v", hcName, err)
				return false, err
			}
			framework.Logf("hc.CheckIntervalSec = %v", hc.CheckIntervalSec)
			return hc.CheckIntervalSec == gceHcCheckIntervalSeconds, nil
		}); pollErr != nil {
			framework.Failf("Health check %q does not reconcile its check interval to %d.", hcName, gceHcCheckIntervalSeconds)
		}
	})

	var _ = SIGDescribe("ESIPP [Slow]", func() {
		f := framework.NewDefaultFramework("esipp")
		var loadBalancerCreateTimeout time.Duration

		var cs clientset.Interface
		serviceLBNames := []string{}

		ginkgo.BeforeEach(func() {
			// requires cloud load-balancer support - this feature currently supported only on GCE/GKE
			e2eskipper.SkipUnlessProviderIs("gce", "gke")

			cs = f.ClientSet
			loadBalancerCreateTimeout = e2eservice.GetServiceLoadBalancerCreationTimeout(cs)
		})

		ginkgo.AfterEach(func() {
			if ginkgo.CurrentGinkgoTestDescription().Failed {
				DescribeSvc(f.Namespace.Name)
			}
			for _, lb := range serviceLBNames {
				framework.Logf("cleaning load balancer resource for %s", lb)
				e2eservice.CleanupServiceResources(cs, lb, framework.TestContext.CloudConfig.Region, framework.TestContext.CloudConfig.Zone)
			}
			//reset serviceLBNames
			serviceLBNames = []string{}
		})

		ginkgo.It("should work for type=LoadBalancer", func() {
			namespace := f.Namespace.Name
			serviceName := "external-local-lb"
			jig := e2eservice.NewTestJig(cs, namespace, serviceName)

			svc, err := jig.CreateOnlyLocalLoadBalancerService(loadBalancerCreateTimeout, true, nil)
			framework.ExpectNoError(err)
			serviceLBNames = append(serviceLBNames, cloudprovider.DefaultLoadBalancerName(svc))
			healthCheckNodePort := int(svc.Spec.HealthCheckNodePort)
			if healthCheckNodePort == 0 {
				framework.Failf("Service HealthCheck NodePort was not allocated")
			}
			defer func() {
				err = jig.ChangeServiceType(v1.ServiceTypeClusterIP, loadBalancerCreateTimeout)
				framework.ExpectNoError(err)

				// Make sure we didn't leak the health check node port.
				threshold := 2
				nodes, err := jig.GetEndpointNodes()
				framework.ExpectNoError(err)
				for _, ips := range nodes {
					err := TestHTTPHealthCheckNodePort(ips[0], healthCheckNodePort, "/healthz", e2eservice.KubeProxyEndpointLagTimeout, false, threshold)
					framework.ExpectNoError(err)
				}
				err = cs.CoreV1().Services(svc.Namespace).Delete(context.TODO(), svc.Name, metav1.DeleteOptions{})
				framework.ExpectNoError(err)
			}()

			svcTCPPort := int(svc.Spec.Ports[0].Port)
			ingressIP := e2eservice.GetIngressPoint(&svc.Status.LoadBalancer.Ingress[0])

			ginkgo.By("reading clientIP using the TCP service's service port via its external VIP")
			content := GetHTTPContent(ingressIP, svcTCPPort, e2eservice.KubeProxyLagTimeout, "/clientip")
			clientIP := content.String()
			framework.Logf("ClientIP detected by target pod using VIP:SvcPort is %s", clientIP)

			ginkgo.By("checking if Source IP is preserved")
			if strings.HasPrefix(clientIP, "10.") {
				framework.Failf("Source IP was NOT preserved")
			}
		})

		ginkgo.It("should work for type=NodePort", func() {
			namespace := f.Namespace.Name
			serviceName := "external-local-nodeport"
			jig := e2eservice.NewTestJig(cs, namespace, serviceName)

			svc, err := jig.CreateOnlyLocalNodePortService(true)
			framework.ExpectNoError(err)
			defer func() {
				err := cs.CoreV1().Services(svc.Namespace).Delete(context.TODO(), svc.Name, metav1.DeleteOptions{})
				framework.ExpectNoError(err)
			}()

			tcpNodePort := int(svc.Spec.Ports[0].NodePort)
			endpointsNodeMap, err := jig.GetEndpointNodes()
			framework.ExpectNoError(err)
			path := "/clientip"

			for nodeName, nodeIPs := range endpointsNodeMap {
				nodeIP := nodeIPs[0]
				ginkgo.By(fmt.Sprintf("reading clientIP using the TCP service's NodePort, on node %v: %v%v%v", nodeName, nodeIP, tcpNodePort, path))
				content := GetHTTPContent(nodeIP, tcpNodePort, e2eservice.KubeProxyLagTimeout, path)
				clientIP := content.String()
				framework.Logf("ClientIP detected by target pod using NodePort is %s", clientIP)
				if strings.HasPrefix(clientIP, "10.") {
					framework.Failf("Source IP was NOT preserved")
				}
			}
		})

		ginkgo.It("should only target nodes with endpoints", func() {
			namespace := f.Namespace.Name
			serviceName := "external-local-nodes"
			jig := e2eservice.NewTestJig(cs, namespace, serviceName)
			nodes, err := e2enode.GetBoundedReadySchedulableNodes(cs, e2eservice.MaxNodesForEndpointsTests)
			framework.ExpectNoError(err)

			svc, err := jig.CreateOnlyLocalLoadBalancerService(loadBalancerCreateTimeout, false,
				func(svc *v1.Service) {
					// Change service port to avoid collision with opened hostPorts
					// in other tests that run in parallel.
					if len(svc.Spec.Ports) != 0 {
						svc.Spec.Ports[0].TargetPort = intstr.FromInt(int(svc.Spec.Ports[0].Port))
						svc.Spec.Ports[0].Port = 8081
					}

				})
			framework.ExpectNoError(err)
			serviceLBNames = append(serviceLBNames, cloudprovider.DefaultLoadBalancerName(svc))
			defer func() {
				err = jig.ChangeServiceType(v1.ServiceTypeClusterIP, loadBalancerCreateTimeout)
				framework.ExpectNoError(err)
				err := cs.CoreV1().Services(svc.Namespace).Delete(context.TODO(), svc.Name, metav1.DeleteOptions{})
				framework.ExpectNoError(err)
			}()

			healthCheckNodePort := int(svc.Spec.HealthCheckNodePort)
			if healthCheckNodePort == 0 {
				framework.Failf("Service HealthCheck NodePort was not allocated")
			}

			ips := e2enode.CollectAddresses(nodes, v1.NodeExternalIP)

			ingressIP := e2eservice.GetIngressPoint(&svc.Status.LoadBalancer.Ingress[0])
			svcTCPPort := int(svc.Spec.Ports[0].Port)

			threshold := 2
			path := "/healthz"
			for i := 0; i < len(nodes.Items); i++ {
				endpointNodeName := nodes.Items[i].Name

				ginkgo.By("creating a pod to be part of the service " + serviceName + " on node " + endpointNodeName)
				_, err = jig.Run(func(rc *v1.ReplicationController) {
					rc.Name = serviceName
					if endpointNodeName != "" {
						rc.Spec.Template.Spec.NodeName = endpointNodeName
					}
				})
				framework.ExpectNoError(err)

				ginkgo.By(fmt.Sprintf("waiting for service endpoint on node %v", endpointNodeName))
				err = jig.WaitForEndpointOnNode(endpointNodeName)
				framework.ExpectNoError(err)

				// HealthCheck should pass only on the node where num(endpoints) > 0
				// All other nodes should fail the healthcheck on the service healthCheckNodePort
				for n, publicIP := range ips {
					// Make sure the loadbalancer picked up the health check change.
					// Confirm traffic can reach backend through LB before checking healthcheck nodeport.
					e2eservice.TestReachableHTTP(ingressIP, svcTCPPort, e2eservice.KubeProxyLagTimeout)
					expectedSuccess := nodes.Items[n].Name == endpointNodeName
					port := strconv.Itoa(healthCheckNodePort)
					ipPort := net.JoinHostPort(publicIP, port)
					framework.Logf("Health checking %s, http://%s%s, expectedSuccess %v", nodes.Items[n].Name, ipPort, path, expectedSuccess)
					err := TestHTTPHealthCheckNodePort(publicIP, healthCheckNodePort, path, e2eservice.KubeProxyEndpointLagTimeout, expectedSuccess, threshold)
					framework.ExpectNoError(err)
				}
				framework.ExpectNoError(e2erc.DeleteRCAndWaitForGC(f.ClientSet, namespace, serviceName))
			}
		})

		ginkgo.It("should work from pods", func() {
			var err error
			namespace := f.Namespace.Name
			serviceName := "external-local-pods"
			jig := e2eservice.NewTestJig(cs, namespace, serviceName)

			svc, err := jig.CreateOnlyLocalLoadBalancerService(loadBalancerCreateTimeout, true, nil)
			framework.ExpectNoError(err)
			serviceLBNames = append(serviceLBNames, cloudprovider.DefaultLoadBalancerName(svc))
			defer func() {
				err = jig.ChangeServiceType(v1.ServiceTypeClusterIP, loadBalancerCreateTimeout)
				framework.ExpectNoError(err)
				err := cs.CoreV1().Services(svc.Namespace).Delete(context.TODO(), svc.Name, metav1.DeleteOptions{})
				framework.ExpectNoError(err)
			}()

			ingressIP := e2eservice.GetIngressPoint(&svc.Status.LoadBalancer.Ingress[0])
			port := strconv.Itoa(int(svc.Spec.Ports[0].Port))
			ipPort := net.JoinHostPort(ingressIP, port)
			path := fmt.Sprintf("%s/clientip", ipPort)

			ginkgo.By("Creating pause pod deployment to make sure, pausePods are in desired state")
			deployment := createPausePodDeployment(cs, "pause-pod-deployment", namespace, 1)
			framework.ExpectNoError(e2edeployment.WaitForDeploymentComplete(cs, deployment), "Failed to complete pause pod deployment")

			defer func() {
				framework.Logf("Deleting deployment")
				err = cs.AppsV1().Deployments(namespace).Delete(context.TODO(), deployment.Name, metav1.DeleteOptions{})
				framework.ExpectNoError(err, "Failed to delete deployment %s", deployment.Name)
			}()

			deployment, err = cs.AppsV1().Deployments(namespace).Get(context.TODO(), deployment.Name, metav1.GetOptions{})
			framework.ExpectNoError(err, "Error in retrieving pause pod deployment")
			labelSelector, err := metav1.LabelSelectorAsSelector(deployment.Spec.Selector)
			framework.ExpectNoError(err, "Error in setting LabelSelector as selector from deployment")

			pausePods, err := cs.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelSelector.String()})
			framework.ExpectNoError(err, "Error in listing pods associated with pause pod deployments")

			pausePod := pausePods.Items[0]
			framework.Logf("Waiting up to %v curl %v", e2eservice.KubeProxyLagTimeout, path)
			cmd := fmt.Sprintf(`curl -q -s --connect-timeout 30 %v`, path)

			var srcIP string
			loadBalancerPropagationTimeout := e2eservice.GetServiceLoadBalancerPropagationTimeout(cs)
			ginkgo.By(fmt.Sprintf("Hitting external lb %v from pod %v on node %v", ingressIP, pausePod.Name, pausePod.Spec.NodeName))
			if pollErr := wait.PollImmediate(framework.Poll, loadBalancerPropagationTimeout, func() (bool, error) {
				stdout, err := framework.RunHostCmd(pausePod.Namespace, pausePod.Name, cmd)
				if err != nil {
					framework.Logf("got err: %v, retry until timeout", err)
					return false, nil
				}
				srcIP = strings.TrimSpace(strings.Split(stdout, ":")[0])
				return srcIP == pausePod.Status.PodIP, nil
			}); pollErr != nil {
				framework.Failf("Source IP not preserved from %v, expected '%v' got '%v'", pausePod.Name, pausePod.Status.PodIP, srcIP)
			}
		})

		// TODO: Get rid of [DisabledForLargeClusters] tag when issue #90047 is fixed.
		ginkgo.It("should handle updates to ExternalTrafficPolicy field [DisabledForLargeClusters]", func() {
			namespace := f.Namespace.Name
			serviceName := "external-local-update"
			jig := e2eservice.NewTestJig(cs, namespace, serviceName)

			nodes, err := e2enode.GetBoundedReadySchedulableNodes(cs, e2eservice.MaxNodesForEndpointsTests)
			framework.ExpectNoError(err)
			if len(nodes.Items) < 2 {
				framework.Failf("Need at least 2 nodes to verify source ip from a node without endpoint")
			}

			svc, err := jig.CreateOnlyLocalLoadBalancerService(loadBalancerCreateTimeout, true, nil)
			framework.ExpectNoError(err)
			serviceLBNames = append(serviceLBNames, cloudprovider.DefaultLoadBalancerName(svc))
			defer func() {
				err = jig.ChangeServiceType(v1.ServiceTypeClusterIP, loadBalancerCreateTimeout)
				framework.ExpectNoError(err)
				err := cs.CoreV1().Services(svc.Namespace).Delete(context.TODO(), svc.Name, metav1.DeleteOptions{})
				framework.ExpectNoError(err)
			}()

			// save the health check node port because it disappears when ESIPP is turned off.
			healthCheckNodePort := int(svc.Spec.HealthCheckNodePort)

			ginkgo.By("turning ESIPP off")
			svc, err = jig.UpdateService(func(svc *v1.Service) {
				svc.Spec.ExternalTrafficPolicy = v1.ServiceExternalTrafficPolicyTypeCluster
			})
			framework.ExpectNoError(err)
			if svc.Spec.HealthCheckNodePort > 0 {
				framework.Failf("Service HealthCheck NodePort still present")
			}

			endpointNodeMap, err := jig.GetEndpointNodes()
			framework.ExpectNoError(err)
			noEndpointNodeMap := map[string][]string{}
			for _, n := range nodes.Items {
				if _, ok := endpointNodeMap[n.Name]; ok {
					continue
				}
				noEndpointNodeMap[n.Name] = e2enode.GetAddresses(&n, v1.NodeExternalIP)
			}

			svcTCPPort := int(svc.Spec.Ports[0].Port)
			svcNodePort := int(svc.Spec.Ports[0].NodePort)
			ingressIP := e2eservice.GetIngressPoint(&svc.Status.LoadBalancer.Ingress[0])
			path := "/clientip"

			ginkgo.By(fmt.Sprintf("endpoints present on nodes %v, absent on nodes %v", endpointNodeMap, noEndpointNodeMap))
			for nodeName, nodeIPs := range noEndpointNodeMap {
				ginkgo.By(fmt.Sprintf("Checking %v (%v:%v%v) proxies to endpoints on another node", nodeName, nodeIPs[0], svcNodePort, path))
				GetHTTPContent(nodeIPs[0], svcNodePort, e2eservice.KubeProxyLagTimeout, path)
			}

			for nodeName, nodeIPs := range endpointNodeMap {
				ginkgo.By(fmt.Sprintf("checking kube-proxy health check fails on node with endpoint (%s), public IP %s", nodeName, nodeIPs[0]))
				var body bytes.Buffer
				pollfn := func() (bool, error) {
					result := e2enetwork.PokeHTTP(nodeIPs[0], healthCheckNodePort, "/healthz", nil)
					if result.Code == 0 {
						return true, nil
					}
					body.Reset()
					body.Write(result.Body)
					return false, nil
				}
				if pollErr := wait.PollImmediate(framework.Poll, e2eservice.TestTimeout, pollfn); pollErr != nil {
					framework.Failf("Kube-proxy still exposing health check on node %v:%v, after ESIPP was turned off. body %s",
						nodeName, healthCheckNodePort, body.String())
				}
			}

			// Poll till kube-proxy re-adds the MASQUERADE rule on the node.
			ginkgo.By(fmt.Sprintf("checking source ip is NOT preserved through loadbalancer %v", ingressIP))
			var clientIP string
			pollErr := wait.PollImmediate(framework.Poll, e2eservice.KubeProxyLagTimeout, func() (bool, error) {
				content := GetHTTPContent(ingressIP, svcTCPPort, e2eservice.KubeProxyLagTimeout, "/clientip")
				clientIP = content.String()
				if strings.HasPrefix(clientIP, "10.") {
					return true, nil
				}
				return false, nil
			})
			if pollErr != nil {
				framework.Failf("Source IP WAS preserved even after ESIPP turned off. Got %v, expected a ten-dot cluster ip.", clientIP)
			}

			// TODO: We need to attempt to create another service with the previously
			// allocated healthcheck nodePort. If the health check nodePort has been
			// freed, the new service creation will succeed, upon which we cleanup.
			// If the health check nodePort has NOT been freed, the new service
			// creation will fail.

			ginkgo.By("setting ExternalTraffic field back to OnlyLocal")
			svc, err = jig.UpdateService(func(svc *v1.Service) {
				svc.Spec.ExternalTrafficPolicy = v1.ServiceExternalTrafficPolicyTypeLocal
				// Request the same healthCheckNodePort as before, to test the user-requested allocation path
				svc.Spec.HealthCheckNodePort = int32(healthCheckNodePort)
			})
			framework.ExpectNoError(err)
			pollErr = wait.PollImmediate(framework.Poll, e2eservice.KubeProxyLagTimeout, func() (bool, error) {
				content := GetHTTPContent(ingressIP, svcTCPPort, e2eservice.KubeProxyLagTimeout, path)
				clientIP = content.String()
				ginkgo.By(fmt.Sprintf("Endpoint %v:%v%v returned client ip %v", ingressIP, svcTCPPort, path, clientIP))
				if !strings.HasPrefix(clientIP, "10.") {
					return true, nil
				}
				return false, nil
			})
			if pollErr != nil {
				framework.Failf("Source IP (%v) is not the client IP even after ESIPP turned on, expected a public IP.", clientIP)
			}
		})
	})
})
