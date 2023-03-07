module github.com/Yelp/paasta-tools-go

require (
	github.com/dlespiau/kube-test-harness v0.0.0-20200730130322-72c5b0037f4a
	github.com/go-logr/zapr v0.1.1 // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/hashicorp/golang-lru v0.5.3 // indirect
	github.com/mitchellh/mapstructure v1.2.2
	github.com/openzipkin/zipkin-go v0.2.2
	github.com/pmezard/go-difflib v1.0.0
	github.com/prometheus/client_golang v1.3.0 // indirect
	github.com/stretchr/testify v1.6.1
	github.com/subosito/gotenv v1.2.0
	go.uber.org/zap v1.13.0 // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/api v0.20.0
	k8s.io/apimachinery v0.20.0
	k8s.io/client-go v0.20.0
	k8s.io/klog v1.0.0
	sigs.k8s.io/controller-runtime v0.6.0
)

go 1.12
