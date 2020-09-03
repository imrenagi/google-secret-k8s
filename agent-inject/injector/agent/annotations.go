package agent

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
)

const (
	DefaultVaultImage = "imrenagi/gsecret-agent:latest"
)

const (

	// AnnotationAgentStatus is the key of the annotation that is added to
	// a pod after an injection is done.
	// There's only one valid status we care about: "injected".
	AnnotationAgentStatus = "google.secret.security.imrenagi.com/agent-inject-status"

	// AnnotationAgentInject is the key of the annotation that controls whether
	// injection is explicitly enabled or disabled for a pod. This should
	// be set to a true or false value, as parseable by strconv.ParseBool
	AnnotationAgentInject = "google.secret.security.imrenagi.com/agent-inject"

	// AnnotationAgentGoogleSecretEntryCRD stores the name of GoogleSecretEntry CRD that will be used to fetch
	// all secret from Google Secret Manager
	AnnotationAgentGoogleSecretEntryCRD = "google.secret.security.imrenagi.com/agent-google-secret-crd"

	// AnnotationAgentRequestNamespace is the Kubernetes namespace where the request
	// originated from.
	AnnotationAgentRequestNamespace = "google.secret.security.imrenagi.com/agent-request-namespace"

	// AnnotationAgentInitFirst makes the initialization container the first container
	// to run when a pod starts. Default is last.
	AnnotationAgentInitFirst = "google.secret.security.imrenagi.com/agent-init-first"

	// AnnotationAgentImage is the name of the Vault docker image to use.
	AnnotationAgentImage = "google.secret.security.imrenagi.com/agent-image"

	// AnnotationAgentPrePopulate controls whether an init container is included
	// to pre-populate the shared memory volume with secrets prior to the application
	// starting.
	AnnotationAgentPrePopulate = "google.secret.security.imrenagi.com/agent-pre-populate"

	// AnnotationAgentPrePopulateOnly controls whether an init container is the only
	// injected container.  If true, no sidecar container will be injected at runtime
	// of the application.
	AnnotationAgentPrePopulateOnly = "google.secret.security.imrenagi.com/agent-pre-populate-only"

	// AnnotationAgentConfigMap is the name of the configuration map where Vault Agent
	// configuration file and templates can be found.
	AnnotationAgentConfigMap = "google.secret.security.imrenagi.com/agent-configmap"

	// AnnotationAgentLimitsCPU sets the CPU limit on the Vault Agent containers.
	AnnotationAgentLimitsCPU = "google.secret.security.imrenagi.com/agent-limits-cpu"

	// AnnotationAgentLimitsMem sets the memory limit on the Vault Agent containers.
	AnnotationAgentLimitsMem = "google.secret.security.imrenagi.com/agent-limits-mem"

	// AnnotationAgentRequestsCPU sets the requested CPU amount on the Vault Agent containers.
	AnnotationAgentRequestsCPU = "google.secret.security.imrenagi.com/agent-requests-cpu"

	// AnnotationAgentRequestsMem sets the requested memory amount on the Vault Agent containers.
	AnnotationAgentRequestsMem = "google.secret.security.imrenagi.com/agent-requests-mem"

	// AnnotationAgentLogLevel sets the Vault Agent log level.
	AnnotationAgentLogLevel = "google.secret.security.imrenagi.com/log-level"

	// AnnotationGoogleSecretManagerClientMaxRetries is the number of retry attempts when 5xx errors are encountered.
	AnnotationGoogleSecretManagerClientMaxRetries = "google.secret.security.imrenagi.com/client-max-retries"

	// AnnotationGoogleSecretManagerClientTimeout sets the request timeout when communicating with Vault.
	AnnotationGoogleSecretManagerClientTimeout = "google.secret.security.imrenagi.com/client-timeout"

	// AnnotationAgentSecretVolumePath specifies where the secrets are to be
	// Mounted after fetching.
	AnnotationAgentSecretVolumePath = "google.secret.security.imrenagi.com/secret-volume-path"
)

type AgentConfig struct {
	Image     string
	Namespace string
}

// Init configures the expected annotations required to create a new instance
// of Agent.  This should be run before running new to ensure all annotations are
// present.
func Init(pod *corev1.Pod, cfg AgentConfig) error {

	if pod == nil {
		return errors.New("pod is empty")
	}

	if cfg.Namespace == "" {
		return errors.New("kubernetes namespace required")
	}

	if pod.ObjectMeta.Annotations == nil {
		pod.ObjectMeta.Annotations = make(map[string]string)
	}

	if _, ok := pod.ObjectMeta.Annotations[AnnotationAgentImage]; !ok {
		if cfg.Image == "" {
			cfg.Image = DefaultVaultImage
		}
		pod.ObjectMeta.Annotations[AnnotationAgentImage] = cfg.Image
	}

	if _, ok := pod.ObjectMeta.Annotations[AnnotationAgentRequestNamespace]; !ok {
		pod.ObjectMeta.Annotations[AnnotationAgentRequestNamespace] = cfg.Namespace
	}

	if _, ok := pod.ObjectMeta.Annotations[AnnotationAgentLimitsCPU]; !ok {
		pod.ObjectMeta.Annotations[AnnotationAgentLimitsCPU] = DefaultResourceLimitCPU
	}

	if _, ok := pod.ObjectMeta.Annotations[AnnotationAgentLimitsMem]; !ok {
		pod.ObjectMeta.Annotations[AnnotationAgentLimitsMem] = DefaultResourceLimitMem
	}

	if _, ok := pod.ObjectMeta.Annotations[AnnotationAgentRequestsCPU]; !ok {
		pod.ObjectMeta.Annotations[AnnotationAgentRequestsCPU] = DefaultResourceRequestCPU
	}

	if _, ok := pod.ObjectMeta.Annotations[AnnotationAgentRequestsMem]; !ok {
		pod.ObjectMeta.Annotations[AnnotationAgentRequestsMem] = DefaultResourceRequestMem
	}

	if _, ok := pod.ObjectMeta.Annotations[AnnotationAgentSecretVolumePath]; !ok {
		pod.ObjectMeta.Annotations[AnnotationAgentSecretVolumePath] = secretVolumePath
	}

	if _, ok := pod.ObjectMeta.Annotations[AnnotationAgentLogLevel]; !ok {
		pod.ObjectMeta.Annotations[AnnotationAgentLogLevel] = DefaultAgentLogLevel
	}

	return nil
}

func (a *Agent) secretEntry() (string, error) {
	name, ok := a.Annotations[AnnotationAgentGoogleSecretEntryCRD]
	if !ok {
		return "", fmt.Errorf("google secret entry annotation doesn't exist")
	}
	return name, nil
}

func (a *Agent) inject() (bool, error) {
	raw, ok := a.Annotations[AnnotationAgentInject]
	if !ok {
		return true, nil
	}
	return strconv.ParseBool(raw)
}

func (a *Agent) initFirst() (bool, error) {
	raw, ok := a.Annotations[AnnotationAgentInitFirst]
	if !ok {
		return false, nil
	}
	return strconv.ParseBool(raw)
}

func (a *Agent) prePopulate() (bool, error) {
	raw, ok := a.Annotations[AnnotationAgentPrePopulate]
	if !ok {
		return true, nil
	}
	return strconv.ParseBool(raw)
}

func (a *Agent) prePopulateOnly() (bool, error) {
	raw, ok := a.Annotations[AnnotationAgentPrePopulateOnly]
	if !ok {
		// TODO this is set to true until we are ready with sidecar container
		return true, nil
	}
	return strconv.ParseBool(raw)
}
