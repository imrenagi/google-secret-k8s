module github.com/imrenagi/google-secret-k8s/secret-operator

go 1.13

require (
	cloud.google.com/go v0.64.0
	github.com/go-logr/logr v0.1.0
	github.com/onsi/ginkgo v1.11.0
	github.com/onsi/gomega v1.8.1
	github.com/rs/zerolog v1.19.0
	google.golang.org/api v0.30.0
	google.golang.org/genproto v0.0.0-20200825200019-8632dd797987
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/api v0.18.2
	k8s.io/apimachinery v0.18.2
	k8s.io/client-go v0.18.2
	sigs.k8s.io/controller-runtime v0.6.0
)
