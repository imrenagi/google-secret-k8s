package cmd

import (
	"context"
	"os"
	"path/filepath"

	"github.com/imrenagi/google-secret-k8s/agent-sidecar/agent"
	secretop "github.com/imrenagi/google-secret-k8s/secret-operator/api"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// NewInitCmd returns init agent command
func NewInitCmd() *cobra.Command {

	var (
		googleSecretEntryName      string
		googleSecretEntryNamespace string
		secretVolumePath           string
	)

	initCmd := cobra.Command{
		Use:   "init",
		Short: "",
		Run: func(cmd *cobra.Command, args []string) {

			var config *rest.Config
			var err error

			kubeconfig := filepath.Join(homeDir(), ".kube", "config")

			// use the current context in kubeconfig
			config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
			if err != nil {
				panic(err.Error())
			}

			// creates the clientset
			clientset, err := kubernetes.NewForConfig(config)
			if err != nil {
				panic(err.Error())
			}

			secretClientset, err := secretop.NewClientSetForConfig(config)
			if err != nil {
				panic(err.Error())
			}

			agent := agent.Agent{
				Clientset:                  clientset,
				SecretSecurityClientset:    secretClientset,
				SecretVolumePath:           secretVolumePath,
				GoogleSecretEntryName:      googleSecretEntryName,
				GoogleSecretEntryNamespace: googleSecretEntryNamespace,
			}

			ctx := context.Background()

			err = agent.SyncSecret(ctx)
			if err != nil {
				panic(err.Error())
			}

		},
	}

	initCmd.Flags().StringVar(&googleSecretEntryName, "secret-entry", "", "google secret entry name")
	initCmd.Flags().StringVar(&googleSecretEntryNamespace, "namespace", "", "google secret entry namespace")
	initCmd.Flags().StringVar(&secretVolumePath, "secret-volume-path", "/google/secrets", "path for storing secret after being fetched from google secret manager")

	return &initCmd
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
