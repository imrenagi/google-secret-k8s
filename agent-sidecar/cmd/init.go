package cmd

import (
	"context"
	"os"
	"path/filepath"

	"github.com/imrenagi/google-secret-k8s/agent-sidecar/agent"
	secretop "github.com/imrenagi/google-secret-k8s/secret-operator/api"
	"github.com/rs/zerolog/log"
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

			log.Debug().Msg("Start init")

			if os.Getenv("ENV") == "development" {
				kubeconfig := filepath.Join(homeDir(), ".kube", "config")
				config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
				if err != nil {
					log.Fatal().Err(err).Msg("unable to create config from kubeconfig")
				}
			} else {
				config, err = rest.InClusterConfig()
				if err != nil {
					log.Fatal().Err(err).Msg("unable to create config from cluster service account")
				}
			}

			log.Debug().Msg("create kubernetes clientset")

			// creates the clientset
			clientset, err := kubernetes.NewForConfig(config)
			if err != nil {
				log.Fatal().Err(err).Msg("unable to create kubernetes clientset")
			}

			log.Debug().Msg("create secret operator clientset")

			secretClientset, err := secretop.NewClientSetForConfig(config)
			if err != nil {
				log.Fatal().Err(err).Msg("unable to create secret security clientsent")
			}

			agent := agent.Agent{
				Clientset:                  clientset,
				SecretSecurityClientset:    secretClientset,
				SecretVolumePath:           secretVolumePath,
				GoogleSecretEntryName:      googleSecretEntryName,
				GoogleSecretEntryNamespace: googleSecretEntryNamespace,
			}

			log.Debug().Msg("create kubernetes clientset")

			ctx := context.Background()
			err = agent.SyncSecret(ctx)
			if err != nil {
				log.Fatal().Err(err).Msg("unable to sync secret to init container")
			}
		},
	}

	initCmd.Flags().StringVar(&googleSecretEntryName, "secret-entry", "", "google secret entry name")
	initCmd.Flags().StringVar(&googleSecretEntryNamespace, "namespace", "", "google secret entry namespace")
	initCmd.Flags().StringVar(&secretVolumePath, "secret-volume-path", "/opt/google-secret/secrets", "path for storing secret after being fetched from google secret manager")

	return &initCmd
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
