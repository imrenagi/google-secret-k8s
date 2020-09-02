package cmd

import (
	"fmt"
	"path/filepath"
	"time"

	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/imrenagi/google-secret-k8s/agent-inject/injector"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// NewServerCmd returns a new `version` command to be used as a sub-command to root
func NewServerCmd() *cobra.Command {

	var (
		requireAnnotation bool
		agentImage        string
	)

	serverCmd := cobra.Command{
		Use:   "server",
		Short: fmt.Sprintf("run server"),
		Run: func(cmd *cobra.Command, args []string) {

			var config *rest.Config
			var err error

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

			// creates the clientset
			clientset, err := kubernetes.NewForConfig(config)
			if err != nil {
				panic(err.Error())
			}

			handler := injector.Handler{
				RequireAnnotation: requireAnnotation,
				ImageSidecar:      agentImage,
				Clientset:         clientset,
				Log:               log.With().Timestamp().Logger(),
			}

			mux := http.NewServeMux()

			mux.HandleFunc("/", home)
			mux.HandleFunc("/mutate", handler.Handle)

			s := &http.Server{
				Addr:           ":8080",
				Handler:        mux,
				ReadTimeout:    10 * time.Second,
				WriteTimeout:   10 * time.Second,
				MaxHeaderBytes: 1 << 20, // 1048576
			}

			go func() {
				log.Printf("listening on %s", s.Addr)
				err := s.ListenAndServeTLS("/cert/agent-injector.pem", "/cert/agent-injector.key")
				if err != nil {
					log.Fatal().Err(err).Msg("cant start server")
				}
			}()

			termChan := make(chan os.Signal)
			signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
			<-termChan

		},
	}

	serverCmd.Flags().StringVar(&agentImage, "agent-image", "imrenagi/gsecret-agent:latest", "Agent image used as init or sidecar container")
	serverCmd.Flags().BoolVar(&requireAnnotation, "require-annotation", true, "If it is true, annotation should be given so that sidecar can be injected")

	return &serverCmd
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("oke"))
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
