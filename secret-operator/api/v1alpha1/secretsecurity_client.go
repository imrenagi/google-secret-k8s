package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
)

// SecretSecurityV1alpha1Interface ...
type SecretSecurityV1alpha1Interface interface {
	RESTClient() rest.Interface
	GoogleSecretEntryGetter
}

// SecretSecurityV1alpha1Client ...
type SecretSecurityV1alpha1Client struct {
	restClient rest.Interface
}

// RESTClient ...
func (s *SecretSecurityV1alpha1Client) RESTClient() rest.Interface {
	if s == nil {
		return nil
	}
	return s.restClient
}

// NewForConfig ...
func NewForConfig(cfg *rest.Config) (*SecretSecurityV1alpha1Client, error) {
	scheme := runtime.NewScheme()
	if err := SchemeBuilder.AddToScheme(scheme); err != nil {
		return nil, err
	}

	config := *cfg
	config.GroupVersion = &GroupVersion
	config.APIPath = "/apis"
	config.ContentType = runtime.ContentTypeJSON
	config.NegotiatedSerializer = serializer.NewCodecFactory(scheme)

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &SecretSecurityV1alpha1Client{
		restClient: client,
	}, nil
}
