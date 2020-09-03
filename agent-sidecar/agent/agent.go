package agent

import (
	"context"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	ioutil "github.com/imrenagi/google-secret-k8s/agent-sidecar/io"
	secretop "github.com/imrenagi/google-secret-k8s/secret-operator/api"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/option"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Agent struct {
	Clientset               *kubernetes.Clientset
	SecretSecurityClientset *secretop.Clientset

	// SecretVolumePath for storing secret retrieved from google secret manager
	SecretVolumePath string

	// GoogleSecretEntryName name of CRD used for storing secret definition
	GoogleSecretEntryName string

	// GoogleSecretEntryNamespace namespace where CRD is located
	GoogleSecretEntryNamespace string
}

// SyncSecret ...
func (a *Agent) SyncSecret(ctx context.Context) error {

	log.Debug().Msg("Start secret sync")

	entry, err := a.SecretSecurityClientset.SecretSecurityV1alpha1().
		GoogleSecretEntry(a.GoogleSecretEntryNamespace).
		Get(ctx, a.GoogleSecretEntryName)
	if err != nil {
		return err
	}

	log.Debug().Msg("retrieving secret")

	secret, err := a.Clientset.CoreV1().Secrets(entry.Spec.SecretRef.Namespace).
		Get(ctx, entry.Spec.SecretRef.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	log.Debug().Msg("creating gcp client api")

	jsoncreds := secret.Data[entry.Spec.SecretRef.DataKey]
	googleSecretMngrClient, err := secretmanager.NewClient(ctx, option.WithCredentialsJSON(jsoncreds))
	if err != nil {
		return err
	}

	for _, secret := range entry.Spec.Secrets {
		res, err := googleSecretMngrClient.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
			Name: secret.Path,
		})
		if err != nil {
			return err
		}

		_, err = ioutil.WriteToFile(a.SecretVolumePath, secret.Name, string(res.Payload.Data), false)
		if err != nil {
			return err
		}

		log.Info().Msgf("successfully get secret for %s", secret.Path)
	}

	return nil
}
