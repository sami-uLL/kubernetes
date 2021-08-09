// This is a generated file. Do not edit directly.

module k8s.io/kubelet

go 1.15

require (
	github.com/gogo/protobuf v1.3.2
	golang.org/x/net v0.0.0-20201224014010-6772e930b67b
	google.golang.org/genproto v0.0.0-20201110150050-8816d57aaa9a // indirect
	google.golang.org/grpc v1.27.1
	k8s.io/api v0.0.0
	k8s.io/apimachinery v0.0.0
	k8s.io/component-base v0.0.0
)

replace (
	github.com/sirupsen/logrus => github.com/sirupsen/logrus v1.6.0
	golang.org/x/net => golang.org/x/net v0.0.0-20201110031124-69a78807bb2b
	golang.org/x/sys => golang.org/x/sys v0.0.0-20201112073958-5cba982894dd
	google.golang.org/protobuf => google.golang.org/protobuf v1.25.0
	k8s.io/api => ../api
	k8s.io/apimachinery => ../apimachinery
	k8s.io/client-go => ../client-go
	k8s.io/component-base => ../component-base
	k8s.io/kubelet => ../kubelet
)
