package v1alpha1

import (
	"github.com/imrenagi/google-secret-k8s/secret-operator/api/v1alpha1"
	"k8s.io/client-go/rest"
)

// Interface ...
type Interface interface {
	SecretSecurityV1alpha1() v1alpha1.SecretSecurityV1alpha1Interface
}

// Clientset ...
type Clientset struct {
	secretSecurityV1alpha1 *v1alpha1.SecretSecurityV1alpha1Client
}

// SecretSecurityV1alpha1 retrieves the SecretSecurityV1alpha1Client
func (c *Clientset) SecretSecurityV1alpha1() v1alpha1.SecretSecurityV1alpha1Interface {
	return c.secretSecurityV1alpha1
}

// NewClientSetForConfig ...
func NewClientSetForConfig(c *rest.Config) (*Clientset, error) {
	var cs Clientset
	var err error
	cs.secretSecurityV1alpha1, err = v1alpha1.NewForConfig(c)
	if err != nil {
		return nil, err
	}
	return &cs, nil
}
