// This is a generated file. Do not edit directly.

module k8s.io/cluster-bootstrap

go 1.19

require (
	github.com/stretchr/testify v1.8.0
	gopkg.in/square/go-jose.v2 v2.2.2
	k8s.io/api v0.0.0
	k8s.io/apimachinery v0.0.0
	k8s.io/klog/v2 v2.80.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/crypto v0.1.0 // indirect
	golang.org/x/net v0.1.1-0.20221027164007-c63010009c80 // indirect
	golang.org/x/text v0.4.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/utils v0.0.0-20220922133306-665eaaec4324 // indirect
	sigs.k8s.io/json v0.0.0-20220713155537-f223a00ba0e2 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
)

replace (
	k8s.io/api => ../api
	k8s.io/apimachinery => ../apimachinery
	k8s.io/cluster-bootstrap => ../cluster-bootstrap
)

replace google.golang.org/grpc => github.com/liggitt/grpc-go v1.51.0-dev.0.20221027215202-2901f263bef4

replace google.golang.org/api => github.com/liggitt/google-api-go-client v0.101.1-0.20221028054038-b4aad6d92467

replace go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc => github.com/liggitt/opentelemetry-go-contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.36.5-0.20221027221714-25c81eb49c35

replace github.com/google/cadvisor => github.com/liggitt/cadvisor v0.45.1-0.20221027232935-03ec2fc20e12
