/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	// secretmanager "cloud.google.com/go/secretmanager/apiv1"
	// "google.golang.org/api/option"
	// secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"

	secretv1alpha1 "github.com/imrenagi/google-secret-k8s/secret-operator/api/v1alpha1"
)

// GoogleSecretEntryReconciler reconciles a GoogleSecretEntry object
type GoogleSecretEntryReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=secret.security.imrenagi.com,resources=googlesecretentries,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=secret.security.imrenagi.com,resources=googlesecretentries/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch

func (r *GoogleSecretEntryReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("googlesecretentry", req.NamespacedName)

	entry := &secretv1alpha1.GoogleSecretEntry{}
	err := r.Get(ctx, req.NamespacedName, entry)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("GoogleSecretEntry resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get GoogleSecretEntry")
		return ctrl.Result{}, err
	}

	if entry.Spec.SecretRef == nil {
		return ctrl.Result{}, fmt.Errorf("spec.secretRef is empty, but is required")
	}

	secretFound := &corev1.Secret{}
	err = r.Get(ctx, types.NamespacedName{Name: entry.Spec.SecretRef.Name, Namespace: entry.Spec.SecretRef.Namespace}, secretFound)
	if err != nil && errors.IsNotFound(err) {
		log.Error(err, "secret not found")
		return ctrl.Result{}, err
	} else if err != nil {
		log.Error(err, "failed to get secret")
		return ctrl.Result{}, err
	}

	if _, ok := secretFound.Data[entry.Spec.SecretRef.DataKey]; !ok {
		log.Error(err, fmt.Sprintf("no data for key %s", entry.Spec.SecretRef.DataKey))
		return ctrl.Result{}, fmt.Errorf("failed to get gcp service account")
	}

	val := secretFound.Data[entry.Spec.SecretRef.DataKey]

	type SA struct {
		ClientEmail string `json:"client_email"`
	}

	var sa SA
	if err := json.Unmarshal(val, &sa); err != nil {
		log.Error(err, "failed to init google secret manager client")
		return ctrl.Result{}, err
	}

	log.Info(fmt.Sprintf("get service account email: %s", sa.ClientEmail))

	entry.Status.ServiceAccountEmail = sa.ClientEmail
	if err := r.Status().Update(ctx, entry); err != nil {
		log.Error(err, "Failed to update GoogleServiceEntry status")
		return ctrl.Result{}, err
	}

	// googleSecretMngrClient, err := secretmanager.NewClient(ctx, option.WithCredentialsJSON(val))
	// if err != nil {
	// 	log.Error(err, "failed to init google secret manager client")
	// 	return ctrl.Result{}, err
	// }

	// for _, secret := range entry.Spec.Secrets {
	// 	res, err := googleSecretMngrClient.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
	// 		Name: secret.Path,
	// 	})
	// 	if err != nil {
	// 		log.Error(err, "failed to get secret from gcp secret manager")
	// 		return ctrl.Result{}, err
	// 	}
	// 	log.Info(fmt.Sprintf("get secret for %s: %s", secret.Path, res.Payload.Data))
	// }

	log.Info(fmt.Sprintf("secret %s found", entry.Spec.SecretRef.Name))

	return ctrl.Result{}, nil
}

func (r *GoogleSecretEntryReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&secretv1alpha1.GoogleSecretEntry{}).
		Owns(&corev1.Secret{}).
		Complete(r)
}
