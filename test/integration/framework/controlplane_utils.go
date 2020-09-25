/*
Copyright 2014 The Kubernetes Authors.

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

package framework

import (
	"context"
	"flag"
	"net"
	"net/http"
	"net/http/httptest"
	"path"
	"strconv"
	"time"

	"github.com/go-openapi/spec"
	"github.com/google/uuid"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	authauthenticator "k8s.io/apiserver/pkg/authentication/authenticator"
	"k8s.io/apiserver/pkg/authentication/authenticatorfactory"
	authenticatorunion "k8s.io/apiserver/pkg/authentication/request/union"
	"k8s.io/apiserver/pkg/authentication/user"
	"k8s.io/apiserver/pkg/authorization/authorizer"
	"k8s.io/apiserver/pkg/authorization/authorizerfactory"
	authorizerunion "k8s.io/apiserver/pkg/authorization/union"
	openapinamer "k8s.io/apiserver/pkg/endpoints/openapi"
	genericfeatures "k8s.io/apiserver/pkg/features"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/apiserver/pkg/server/options"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	"k8s.io/apiserver/pkg/storage/storagebackend"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	utilflowcontrol "k8s.io/apiserver/pkg/util/flowcontrol"
	"k8s.io/client-go/informers"
	clientset "k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/component-base/version"
	"k8s.io/klog/v2"
	openapicommon "k8s.io/kube-openapi/pkg/common"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	"k8s.io/kubernetes/pkg/controlplane"
	"k8s.io/kubernetes/pkg/generated/openapi"
	"k8s.io/kubernetes/pkg/kubeapiserver"
	kubeletclient "k8s.io/kubernetes/pkg/kubelet/client"
)

// alwaysAllow always allows an action
type alwaysAllow struct{}

func (alwaysAllow) Authorize(ctx context.Context, requestAttributes authorizer.Attributes) (authorizer.Decision, string, error) {
	return authorizer.DecisionAllow, "always allow", nil
}

// alwaysEmpty simulates "no authentication" for old tests
func alwaysEmpty(req *http.Request) (*authauthenticator.Response, bool, error) {
	return &authauthenticator.Response{
		User: &user.DefaultInfo{
			Name: "",
		},
	}, true, nil
}

// ControlPlaneReceiver can be used to provide the control plane to a custom incoming server function
type ControlPlaneReceiver interface {
	SetInstance(i *controlplane.Instance)
}

// ControlPlaneHolder implements ControlPlaneReceiver
type ControlPlaneHolder struct {
	Initialized chan struct{}
	Instance    *controlplane.Instance
}

// SetInstance assigns the current control plane.
func (h *ControlPlaneHolder) SetInstance(i *controlplane.Instance) {
	h.Instance = i
	close(h.Initialized)
}

// DefaultOpenAPIConfig returns an openapicommon.Config initialized to default values.
func DefaultOpenAPIConfig() *openapicommon.Config {
	openAPIConfig := genericapiserver.DefaultOpenAPIConfig(openapi.GetOpenAPIDefinitions, openapinamer.NewDefinitionNamer(legacyscheme.Scheme))
	openAPIConfig.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:   "Kubernetes",
			Version: "unversioned",
		},
	}
	openAPIConfig.DefaultResponse = &spec.Response{
		ResponseProps: spec.ResponseProps{
			Description: "Default Response.",
		},
	}
	openAPIConfig.GetDefinitions = openapi.GetOpenAPIDefinitions

	return openAPIConfig
}

// startControlPlaneOrDie starts a kubernetes control plane and an httpserver to handle api requests
func startControlPlaneOrDie(controlPlaneConfig *controlplane.Config, incomingServer *httptest.Server, controlPlaneReceiver ControlPlaneReceiver) (*controlplane.Instance, *httptest.Server, CloseFunc) {
	var instance *controlplane.Instance
	var s *httptest.Server

	// Ensure we log at least level 4
	v := flag.Lookup("v").Value
	level, _ := strconv.Atoi(v.String())
	if level < 4 {
		v.Set("4")
	}

	if incomingServer != nil {
		s = incomingServer
	} else {
		s = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			instance.GenericAPIServer.Handler.ServeHTTP(w, req)
		}))
	}

	stopCh := make(chan struct{})
	closeFn := func() {
		if instance != nil {
			instance.GenericAPIServer.RunPreShutdownHooks()
		}
		close(stopCh)
		s.Close()
	}

	if controlPlaneConfig == nil {
		controlPlaneConfig = NewControlPlaneConfig()
		controlPlaneConfig.GenericConfig.OpenAPIConfig = DefaultOpenAPIConfig()
	}

	// set the loopback client config
	if controlPlaneConfig.GenericConfig.LoopbackClientConfig == nil {
		controlPlaneConfig.GenericConfig.LoopbackClientConfig = &restclient.Config{QPS: 50, Burst: 100, ContentConfig: restclient.ContentConfig{NegotiatedSerializer: legacyscheme.Codecs}}
	}
	controlPlaneConfig.GenericConfig.LoopbackClientConfig.Host = s.URL

	privilegedLoopbackToken := uuid.New().String()
	// wrap any available authorizer
	tokens := make(map[string]*user.DefaultInfo)
	tokens[privilegedLoopbackToken] = &user.DefaultInfo{
		Name:   user.APIServerUser,
		UID:    uuid.New().String(),
		Groups: []string{user.SystemPrivilegedGroup},
	}

	tokenAuthenticator := authenticatorfactory.NewFromTokens(tokens)
	if controlPlaneConfig.GenericConfig.Authentication.Authenticator == nil {
		controlPlaneConfig.GenericConfig.Authentication.Authenticator = authenticatorunion.New(tokenAuthenticator, authauthenticator.RequestFunc(alwaysEmpty))
	} else {
		controlPlaneConfig.GenericConfig.Authentication.Authenticator = authenticatorunion.New(tokenAuthenticator, controlPlaneConfig.GenericConfig.Authentication.Authenticator)
	}

	if controlPlaneConfig.GenericConfig.Authorization.Authorizer != nil {
		tokenAuthorizer := authorizerfactory.NewPrivilegedGroups(user.SystemPrivilegedGroup)
		controlPlaneConfig.GenericConfig.Authorization.Authorizer = authorizerunion.New(tokenAuthorizer, controlPlaneConfig.GenericConfig.Authorization.Authorizer)
	} else {
		controlPlaneConfig.GenericConfig.Authorization.Authorizer = alwaysAllow{}
	}

	controlPlaneConfig.GenericConfig.LoopbackClientConfig.BearerToken = privilegedLoopbackToken

	clientset, err := clientset.NewForConfig(controlPlaneConfig.GenericConfig.LoopbackClientConfig)
	if err != nil {
		klog.Fatal(err)
	}

	controlPlaneConfig.ExtraConfig.VersionedInformers = informers.NewSharedInformerFactory(clientset, controlPlaneConfig.GenericConfig.LoopbackClientConfig.Timeout)

	if utilfeature.DefaultFeatureGate.Enabled(genericfeatures.APIPriorityAndFairness) {
		controlPlaneConfig.GenericConfig.FlowControl = utilflowcontrol.New(
			controlPlaneConfig.ExtraConfig.VersionedInformers,
			clientset.FlowcontrolV1alpha1(),
			controlPlaneConfig.GenericConfig.MaxRequestsInFlight+controlPlaneConfig.GenericConfig.MaxMutatingRequestsInFlight,
			controlPlaneConfig.GenericConfig.RequestTimeout/4,
		)
	}

	instance, err = controlPlaneConfig.Complete().New(genericapiserver.NewEmptyDelegate())
	if err != nil {
		// We log the error first so that even if closeFn crashes, the error is shown
		klog.Errorf("error in bringing up the control plane: %v", err)
		closeFn()
		klog.Fatalf("error in bringing up the control plane: %v", err)
	}
	if controlPlaneReceiver != nil {
		controlPlaneReceiver.SetInstance(instance)
	}

	// TODO have this start method actually use the normal start sequence for the API server
	// this method never actually calls the `Run` method for the API server
	// fire the post hooks ourselves
	instance.GenericAPIServer.PrepareRun()
	instance.GenericAPIServer.RunPostStartHooks(stopCh)

	cfg := *controlPlaneConfig.GenericConfig.LoopbackClientConfig
	cfg.ContentConfig.GroupVersion = &schema.GroupVersion{}
	privilegedClient, err := restclient.RESTClientFor(&cfg)
	if err != nil {
		closeFn()
		klog.Fatal(err)
	}
	var lastHealthContent []byte
	err = wait.PollImmediate(100*time.Millisecond, 30*time.Second, func() (bool, error) {
		result := privilegedClient.Get().AbsPath("/healthz").Do(context.TODO())
		status := 0
		result.StatusCode(&status)
		if status == 200 {
			return true, nil
		}
		lastHealthContent, _ = result.Raw()
		return false, nil
	})
	if err != nil {
		closeFn()
		klog.Errorf("last health content: %q", string(lastHealthContent))
		klog.Fatal(err)
	}

	return instance, s, closeFn
}

// NewIntegrationTestControlPlaneConfig returns the control plane config appropriate for most integration tests.
func NewIntegrationTestControlPlaneConfig() *controlplane.Config {
	return NewIntegrationTestControlPlaneConfigWithOptions(&ControlPlaneConfigOptions{})
}

// NewIntegrationTestControlPlaneConfigWithOptions returns the control plane config appropriate for most integration tests
// configured with the provided options.
func NewIntegrationTestControlPlaneConfigWithOptions(opts *ControlPlaneConfigOptions) *controlplane.Config {
	controlPlaneConfig := NewControlPlaneConfigWithOptions(opts)
	controlPlaneConfig.GenericConfig.PublicAddress = net.ParseIP("192.168.10.4")
	controlPlaneConfig.ExtraConfig.APIResourceConfigSource = controlplane.DefaultAPIResourceConfigSource()

	// TODO: get rid of these tests or port them to secure serving
	controlPlaneConfig.GenericConfig.SecureServing = &genericapiserver.SecureServingInfo{Listener: fakeLocalhost443Listener{}}

	return controlPlaneConfig
}

// ControlPlaneConfigOptions are the configurable options for a new integration test control plane config.
type ControlPlaneConfigOptions struct {
	EtcdOptions *options.EtcdOptions
}

// DefaultEtcdOptions are the default EtcdOptions for use with integration tests.
func DefaultEtcdOptions() *options.EtcdOptions {
	// This causes the integration tests to exercise the etcd
	// prefix code, so please don't change without ensuring
	// sufficient coverage in other ways.
	etcdOptions := options.NewEtcdOptions(storagebackend.NewDefaultConfig(uuid.New().String(), nil))
	etcdOptions.StorageConfig.Transport.ServerList = []string{GetEtcdURL()}
	return etcdOptions
}

// NewControlPlaneConfig returns a basic control plane config.
func NewControlPlaneConfig() *controlplane.Config {
	return NewControlPlaneConfigWithOptions(&ControlPlaneConfigOptions{})
}

// NewControlPlaneConfigWithOptions returns a basic control plane config configured with the provided options.
func NewControlPlaneConfigWithOptions(opts *ControlPlaneConfigOptions) *controlplane.Config {
	etcdOptions := DefaultEtcdOptions()
	if opts.EtcdOptions != nil {
		etcdOptions = opts.EtcdOptions
	}

	storageConfig := kubeapiserver.NewStorageFactoryConfig()
	storageConfig.APIResourceConfig = serverstorage.NewResourceConfig()
	completedStorageConfig, err := storageConfig.Complete(etcdOptions)
	if err != nil {
		panic(err)
	}
	storageFactory, err := completedStorageConfig.New()
	if err != nil {
		panic(err)
	}

	genericConfig := genericapiserver.NewConfig(legacyscheme.Codecs)
	kubeVersion := version.Get()
	genericConfig.Version = &kubeVersion
	genericConfig.Authorization.Authorizer = authorizerfactory.NewAlwaysAllowAuthorizer()

	// TODO: get rid of these tests or port them to secure serving
	genericConfig.SecureServing = &genericapiserver.SecureServingInfo{Listener: fakeLocalhost443Listener{}}

	err = etcdOptions.ApplyWithStorageFactoryTo(storageFactory, genericConfig)
	if err != nil {
		panic(err)
	}

	return &controlplane.Config{
		GenericConfig: genericConfig,
		ExtraConfig: controlplane.ExtraConfig{
			APIResourceConfigSource: controlplane.DefaultAPIResourceConfigSource(),
			StorageFactory:          storageFactory,
			KubeletClientConfig:     kubeletclient.KubeletClientConfig{Port: 10250},
			APIServerServicePort:    443,
			MasterCount:             1,
		},
	}
}

// CloseFunc can be called to cleanup the control plane
type CloseFunc func()

// RunAControlPlane starts a control plane with the provided config.
func RunAControlPlane(controlPlaneConfig *controlplane.Config) (*controlplane.Instance, *httptest.Server, CloseFunc) {
	if controlPlaneConfig == nil {
		controlPlaneConfig = NewControlPlaneConfig()
		controlPlaneConfig.GenericConfig.EnableProfiling = true
	}
	return startControlPlaneOrDie(controlPlaneConfig, nil, nil)
}

// RunAControlPlaneUsingServer starts up a control plane using the provided config on the specified server.
func RunAControlPlaneUsingServer(controlPlaneConfig *controlplane.Config, s *httptest.Server, controlPlaneReceiver ControlPlaneReceiver) (*controlplane.Instance, *httptest.Server, CloseFunc) {
	return startControlPlaneOrDie(controlPlaneConfig, s, controlPlaneReceiver)
}

// SharedEtcd creates a storage config for a shared etcd instance, with a unique prefix.
func SharedEtcd() *storagebackend.Config {
	cfg := storagebackend.NewDefaultConfig(path.Join(uuid.New().String(), "registry"), nil)
	cfg.Transport.ServerList = []string{GetEtcdURL()}
	return cfg
}

type fakeLocalhost443Listener struct{}

func (fakeLocalhost443Listener) Accept() (net.Conn, error) {
	return nil, nil
}

func (fakeLocalhost443Listener) Close() error {
	return nil
}

func (fakeLocalhost443Listener) Addr() net.Addr {
	return &net.TCPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 443,
	}
}
