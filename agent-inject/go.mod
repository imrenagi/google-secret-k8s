module github.com/imrenagi/google-secret-k8s/agent-inject

go 1.15

require (
	github.com/hashicorp/go-hclog v0.14.1
	github.com/hashicorp/vault v0.10.3
	github.com/hashicorp/vault-k8s v0.5.0
	github.com/hashicorp/vault/sdk v0.1.14-0.20191205220236-47cffd09f972
	github.com/mattbaird/jsonpatch v0.0.0-20200820163806-098863c1fc24
	github.com/pkg/errors v0.9.1
	github.com/rs/zerolog v1.19.0
	github.com/spf13/cobra v1.0.0
	github.com/stretchr/testify v1.4.0
	gomodules.xyz/jsonpatch/v2 v2.1.0
	k8s.io/api v0.0.0-20190409021203-6e4e0e4f393b
	k8s.io/apimachinery v0.0.0-20190404173353-6a84e37a896d
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	k8s.io/klog v1.0.0 // indirect
	k8s.io/utils v0.0.0-20191030222137-2b95a09bc58d // indirect
	sigs.k8s.io/yaml v1.1.0 // indirect
)
