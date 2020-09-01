package v1alpha1

import (
	"context"

	"k8s.io/client-go/rest"
)

type GoogleSecretEntryGetter interface {
	GoogleSecretEntry(namespace string) GoogleSecretEntryInterface
}

// GoogleSecretEntryInterface has methods to work with GoogleSecretEntry resources.
type GoogleSecretEntryInterface interface {
	Get(ctx context.Context, name string) (*GoogleSecretEntry, error)
}

// GoogleSecretEntry ...
func (s SecretSecurityV1alpha1Client) GoogleSecretEntry(namespace string) GoogleSecretEntryInterface {
	var ns string
	if ns == "" {
		ns = "default"
	}
	return &googleSecretEntryClient{
		client:    s.RESTClient(),
		namespace: ns,
	}
}

type googleSecretEntryClient struct {
	client    rest.Interface
	namespace string
}

func (g *googleSecretEntryClient) Get(ctx context.Context, name string) (*GoogleSecretEntry, error) {
	var entry GoogleSecretEntry
	err := g.client.Get().
		Namespace(g.namespace).
		Resource("googlesecretentries").
		Name(name).
		Do(ctx).Into(&entry)
	return &entry, err
}
