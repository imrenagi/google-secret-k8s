package agent

import (
	"context"
	"fmt"

	ioutil "github.com/imrenagi/google-secret-k8s/agent-sidecar/io"
	secretop "github.com/imrenagi/google-secret-k8s/secret-operator/api"
	"github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"google.golang.org/api/option"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

type Agent struct {
	Clientset               *kubernetes.Clientset
	SecretSecurityClientset *secretop.Clientset
	SecretVolumePath        string
}

func (a *Agent) SyncSecret() error {

	ctx := context.Background()

	entry, err := a.SecretSecurityClientset.SecretSecurityV1alpha1().
		GoogleSecretEntry("").Get(ctx, "googlesecretentry-sample")
	if err != nil {
		return err
	}

	fmt.Println(entry.Name)

	secret, err := a.Clientset.CoreV1().Secrets(entry.Spec.SecretRef.Namespace).
		Get(ctx, entry.Spec.SecretRef.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	b := secret.Data[entry.Spec.SecretRef.DataKey]

	googleSecretMngrClient, err := secretmanager.NewClient(ctx, option.WithCredentialsJSON(b))
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

		log.Info().Msg(fmt.Sprintf("get secret for %s: %s", secret.Path, res.Payload.Data))
	}

	return nil
}
