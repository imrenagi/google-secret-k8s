module github.com/imrenagi/google-secret-k8s/agent-sidecar

go 1.15

require (
	cloud.google.com/go v0.65.0
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/googleapis/gnostic v0.5.1 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/imrenagi/google-secret-k8s/secret-operator v0.0.0-20200903152213-75074c077c01
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/rs/zerolog v1.19.0
	github.com/spf13/cobra v1.0.0	
	google.golang.org/api v0.30.0
	google.golang.org/genproto v0.0.0-20200825200019-8632dd797987	
	k8s.io/api v0.19.0
	k8s.io/apimachinery v0.19.0
	k8s.io/client-go v0.19.0
	sigs.k8s.io/controller-runtime v0.6.2
)
