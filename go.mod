module github.com/vmware-tanzu-labs/namespace-operator

go 1.16

require (
	github.com/banzaicloud/k8s-objectmatcher v1.6.1
	github.com/banzaicloud/operator-tools v0.26.0
	github.com/go-logr/logr v0.4.0
	github.com/imdario/mergo v0.3.12
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.14.0
	github.com/spf13/cobra v1.1.3
	k8s.io/api v0.21.3
	k8s.io/apiextensions-apiserver v0.21.3
	k8s.io/apimachinery v0.21.3
	k8s.io/client-go v0.21.3
	sigs.k8s.io/controller-runtime v0.9.5
	sigs.k8s.io/yaml v1.2.0
)
